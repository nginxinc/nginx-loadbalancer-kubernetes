/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/observation"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/probation"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/synchronization"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

func main() {
	err := run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()
	var err error

	k8sClient, err := buildKubernetesClient()
	if err != nil {
		return fmt.Errorf(`error building a Kubernetes client: %w`, err)
	}

	settings, err := configuration.NewSettings(ctx, k8sClient)
	if err != nil {
		return fmt.Errorf(`error occurred creating settings: %w`, err)
	}

	err = settings.Initialize()
	if err != nil {
		return fmt.Errorf(`error occurred initializing settings: %w`, err)
	}

	go settings.Run()

	synchronizerWorkqueue, err := buildWorkQueue(settings.Synchronizer.WorkQueueSettings)
	if err != nil {
		return fmt.Errorf(`error occurred building a workqueue: %w`, err)
	}

	synchronizer, err := synchronization.NewSynchronizer(settings, synchronizerWorkqueue)
	if err != nil {
		return fmt.Errorf(`error initializing synchronizer: %w`, err)
	}

	handlerWorkqueue, err := buildWorkQueue(settings.Synchronizer.WorkQueueSettings)
	if err != nil {
		return fmt.Errorf(`error occurred building a workqueue: %w`, err)
	}

	handler := observation.NewHandler(settings, synchronizer, handlerWorkqueue)

	watcher, err := observation.NewWatcher(settings, handler)
	if err != nil {
		return fmt.Errorf(`error occurred creating a watcher: %w`, err)
	}

	err = watcher.Initialize()
	if err != nil {
		return fmt.Errorf(`error occurred initializing the watcher: %w`, err)
	}

	go handler.Run(ctx.Done())
	go synchronizer.Run(ctx.Done())

	probeServer := probation.NewHealthServer()
	probeServer.Start()

	err = watcher.Watch()
	if err != nil {
		return fmt.Errorf(`error occurred watching for events: %w`, err)
	}

	<-ctx.Done()
	return nil
}

// buildKubernetesClient builds a Kubernetes clientset, supporting both in-cluster and out-of-cluster (kubeconfig) configurations.
func buildKubernetesClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first
	config, err = rest.InClusterConfig()
	if err != nil {
		if err == rest.ErrNotInCluster {
			// Not running in a cluster, fall back to kubeconfig
			kubeconfigPath := os.Getenv("KUBECONFIG")
			if kubeconfigPath == "" {
				kubeconfigPath = clientcmd.RecommendedHomeFile // ~/.kube/config
			}

			config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
			if err != nil {
				return nil, fmt.Errorf("could not get Kubernetes config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error occurred getting the in-cluster config: %w", err)
		}
	}

	// Create the clientset
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("error occurred creating a Kubernetes client: %w", err)
	}

	return client, nil
}

func buildWorkQueue(settings configuration.WorkQueueSettings) (workqueue.RateLimitingInterface, error) {
	logrus.Debug("Watcher::buildSynchronizerWorkQueue")

	rateLimiter := workqueue.NewItemExponentialFailureRateLimiter(settings.RateLimiterBase, settings.RateLimiterMax)
	return workqueue.NewNamedRateLimitingQueue(rateLimiter, settings.Name), nil
}

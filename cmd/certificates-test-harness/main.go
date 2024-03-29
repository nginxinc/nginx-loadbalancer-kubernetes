package main

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	err := run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	logrus.Info("certificates-test-harness::run")

	ctx := context.Background()
	var err error

	k8sClient, err := buildKubernetesClient()
	if err != nil {
		return fmt.Errorf(`error building a Kubernetes client: %w`, err)
	}

	certificates := certification.NewCertificates(ctx, k8sClient)

	err = certificates.Initialize()
	if err != nil {
		return fmt.Errorf(`error occurred initializing certificates: %w`, err)
	}

	go certificates.Run()

	<-ctx.Done()
	return nil
}

func buildKubernetesClient() (*kubernetes.Clientset, error) {
	logrus.Debug("Watcher::buildKubernetesClient")

	var kubeconfig *string
	var k8sConfig *rest.Config

	k8sConfig, err := rest.InClusterConfig()
	if errors.Is(err, rest.ErrNotInCluster) {
		if home := homedir.HomeDir(); home != "" {
			path := filepath.Join(home, ".kube", "config")
			kubeconfig = &path

			k8sConfig, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
			if err != nil {
				return nil, fmt.Errorf(`error occurred building the kubeconfig: %w`, err)
			}
		} else {
			return nil, fmt.Errorf(`not running in a Cluster: %w`, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf(`error occurred getting the Cluster config: %w`, err)
	}

	client, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		return nil, fmt.Errorf(`error occurred creating a client: %w`, err)
	}
	return client, nil
}

/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 *
 * Establishes a Watcher for the Kubernetes Secrets that contain the various certificates and keys used to generate a tls.Config object;
 * exposes the certificates and keys.
 */

package certification

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	// SecretsNamespace is the value used to filter the Secrets Resource in the Informer.
	SecretsNamespace = "nlk"
)

type Certificates struct {
	// CACertificate is the PEM formatted Certificate Authority certificate.
	CACertificate string

	// ClientCertificate is the PEM formatted Client certificate.
	ClientCertificate string

	// ClientKey is the PEM formatted Client key.
	ClientKey string

	// Context is the context used to control the application.
	Context context.Context

	// informer is the SharedInformer used to watch for changes to the Secrets .
	informer cache.SharedInformer

	// K8sClient is the Kubernetes client used to communicate with the Kubernetes API.
	k8sClient *kubernetes.Clientset

	// eventHandlerRegistration is the object used to track the event handlers with the SharedInformer.
	eventHandlerRegistration cache.ResourceEventHandlerRegistration
}

func NewCertificates(ctx context.Context, k8sClient *kubernetes.Clientset) (*Certificates, error) {
	return &Certificates{
		k8sClient: k8sClient,
		Context:   ctx,
	}, nil
}

// GetCACertificate returns the Certificate Authority certificate.
func (c *Certificates) GetCACertificate() []byte {
	return []byte(c.CACertificate)
}

// GetClientCertificate returns the Client certificate and key.
func (c *Certificates) GetClientCertificate() ([]byte, []byte) {
	clientKey := []byte(c.ClientKey)
	clientCertificate := []byte(c.ClientCertificate)

	return clientKey, clientCertificate
}

// Initialize initializes the Certificates object. Sets up a SharedInformer for the Secrets Resource.
func (c *Certificates) Initialize() error {
	logrus.Info("Certificates::Initialize")

	var err error

	informer, err := c.buildInformer()
	if err != nil {
		return fmt.Errorf(`error occurred building an informer: %w`, err)
	}

	c.informer = informer

	err = c.initializeEventHandlers()
	if err != nil {
		return fmt.Errorf(`error occurred initializing event handlers: %w`, err)
	}

	return nil
}

func (c *Certificates) Run() {
	logrus.Info("Certificates::Run")

	c.informer.Run(c.Context.Done())

	<-c.Context.Done()
}

func (c *Certificates) buildInformer() (cache.SharedInformer, error) {
	logrus.Debug("Certificates::buildInformer")

	options := informers.WithNamespace(SecretsNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(c.k8sClient, 0, options)
	informer := factory.Core().V1().Secrets().Informer()

	return informer, nil
}

func (c *Certificates) initializeEventHandlers() error {
	logrus.Debug("Certificates::initializeEventHandlers")

	var err error

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    c.handleAddEvent,
		DeleteFunc: c.handleDeleteEvent,
		UpdateFunc: c.handleUpdateEvent,
	}

	c.eventHandlerRegistration, err = c.informer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred registering event handlers: %w`, err)
	}

	return nil
}

func (c *Certificates) handleAddEvent(obj interface{}) {
	logrus.Debug("Certificates::handleAddEvent")

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleAddEvent: unable to cast object to Secret")
		return
	}

	logrus.Warnf("Certificates::handleAddEvent: secret name: %v", secret.Name)
}

func (c *Certificates) handleDeleteEvent(obj interface{}) {
	logrus.Debug("Certificates::handleDeleteEvent")

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleDeleteEvent: unable to cast object to Secret")
		return
	}

	logrus.Warnf("Certificates::handleDeleteEvent: secret name: %v", secret.Name)
}

func (c *Certificates) handleUpdateEvent(obj interface{}, obj2 interface{}) {
	logrus.Debug("Certificates::handleUpdateEvent")
	prevSecret, ok := obj.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleUpdateEvent: unable to cast object to Secret")
		return
	}
	newSecret, ok := obj.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleUpdateEvent: unable to cast object to Secret")
		return
	}

	logrus.Warnf("Certificates::handleUpdateEvent: previous secret name: %v, new secret name: %v", prevSecret.Name, newSecret.Name)
}

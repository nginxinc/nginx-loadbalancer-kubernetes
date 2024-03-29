/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 *
 * Establishes a Watcher for the Kubernetes Secrets that contain the various certificates
 * and keys used to generate a tls.Config object;
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

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/core"
)

const (
	// SecretsNamespace is the value used to filter the Secrets Resource in the Informer.
	SecretsNamespace = "nlk"

	// CertificateKey is the key for the certificate in the Secret.
	CertificateKey = "tls.crt"

	// CertificateKeyKey is the key for the certificate key in the Secret.
	CertificateKeyKey = "tls.key"
)

type Certificates struct {
	Certificates map[string]map[string]core.SecretBytes

	// Context is the context used to control the application.
	Context context.Context

	// CaCertificateSecretKey is the name of the Secret that contains the Certificate Authority certificate.
	CaCertificateSecretKey string

	// ClientCertificateSecretKey is the name of the Secret that contains the Client certificate.
	ClientCertificateSecretKey string

	// informer is the SharedInformer used to watch for changes to the Secrets .
	informer cache.SharedInformer

	// K8sClient is the Kubernetes client used to communicate with the Kubernetes API.
	k8sClient kubernetes.Interface

	// eventHandlerRegistration is the object used to track the event handlers with the SharedInformer.
	eventHandlerRegistration cache.ResourceEventHandlerRegistration
}

// NewCertificates factory method that returns a new Certificates object.
func NewCertificates(ctx context.Context, k8sClient kubernetes.Interface) *Certificates {
	return &Certificates{
		k8sClient:    k8sClient,
		Context:      ctx,
		Certificates: nil,
	}
}

// GetCACertificate returns the Certificate Authority certificate.
func (c *Certificates) GetCACertificate() core.SecretBytes {
	bytes := c.Certificates[c.CaCertificateSecretKey][CertificateKey]

	return bytes
}

// GetClientCertificate returns the Client certificate and key.
func (c *Certificates) GetClientCertificate() (core.SecretBytes, core.SecretBytes) {
	keyBytes := c.Certificates[c.ClientCertificateSecretKey][CertificateKeyKey]
	certificateBytes := c.Certificates[c.ClientCertificateSecretKey][CertificateKey]

	return keyBytes, certificateBytes
}

// Initialize initializes the Certificates object. Sets up a SharedInformer for the Secrets Resource.
func (c *Certificates) Initialize() error {
	logrus.Info("Certificates::Initialize")

	var err error

	c.Certificates = make(map[string]map[string]core.SecretBytes)

	informer := c.buildInformer()

	c.informer = informer

	err = c.initializeEventHandlers()
	if err != nil {
		return fmt.Errorf(`error occurred initializing event handlers: %w`, err)
	}

	return nil
}

// Run starts the SharedInformer.
func (c *Certificates) Run() error {
	logrus.Info("Certificates::Run")

	if c.informer == nil {
		return fmt.Errorf(`initialize must be called before Run`)
	}

	c.informer.Run(c.Context.Done())

	<-c.Context.Done()

	return nil
}

func (c *Certificates) buildInformer() cache.SharedInformer {
	logrus.Debug("Certificates::buildInformer")

	options := informers.WithNamespace(SecretsNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(c.k8sClient, 0, options)
	informer := factory.Core().V1().Secrets().Informer()

	return informer
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

	c.Certificates[secret.Name] = map[string]core.SecretBytes{}

	// Input from the secret comes in the form
	//   tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVCVEN...
	//   tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0l...
	// Where the keys are `tls.crt` and `tls.key` and the values are []byte
	for k, v := range secret.Data {
		c.Certificates[secret.Name][k] = core.SecretBytes(v)
	}

	logrus.Debugf("Certificates::handleAddEvent: certificates (%d)", len(c.Certificates))
}

func (c *Certificates) handleDeleteEvent(obj interface{}) {
	logrus.Debug("Certificates::handleDeleteEvent")

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleDeleteEvent: unable to cast object to Secret")
		return
	}

	if c.Certificates[secret.Name] != nil {
		delete(c.Certificates, secret.Name)
	}

	logrus.Debugf("Certificates::handleDeleteEvent: certificates (%d)", len(c.Certificates))
}

func (c *Certificates) handleUpdateEvent(_ interface{}, newValue interface{}) {
	logrus.Debug("Certificates::handleUpdateEvent")

	secret, ok := newValue.(*corev1.Secret)
	if !ok {
		logrus.Errorf("Certificates::handleUpdateEvent: unable to cast object to Secret")
		return
	}

	for k, v := range secret.Data {
		c.Certificates[secret.Name][k] = v
	}

	logrus.Debugf("Certificates::handleUpdateEvent: certificates (%d)", len(c.Certificates))
}

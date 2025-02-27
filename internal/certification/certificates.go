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
	"log/slog"
	"sync"

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
	mu           sync.Mutex // guards Certificates
	certificates map[string]map[string]core.SecretBytes

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
func NewCertificates(
	k8sClient kubernetes.Interface, certificates map[string]map[string]core.SecretBytes,
) *Certificates {
	return &Certificates{
		k8sClient:    k8sClient,
		certificates: certificates,
	}
}

// GetCACertificate returns the Certificate Authority certificate.
func (c *Certificates) GetCACertificate() core.SecretBytes {
	c.mu.Lock()
	defer c.mu.Unlock()

	bytes := c.certificates[c.CaCertificateSecretKey][CertificateKey]

	return bytes
}

// GetClientCertificate returns the Client certificate and key.
func (c *Certificates) GetClientCertificate() (core.SecretBytes, core.SecretBytes) {
	c.mu.Lock()
	defer c.mu.Unlock()

	keyBytes := c.certificates[c.ClientCertificateSecretKey][CertificateKeyKey]
	certificateBytes := c.certificates[c.ClientCertificateSecretKey][CertificateKey]

	return keyBytes, certificateBytes
}

// Initialize initializes the Certificates object. Sets up a SharedInformer for the Secrets Resource.
func (c *Certificates) Initialize() error {
	slog.Info("Certificates::Initialize")

	var err error

	c.mu.Lock()
	c.certificates = make(map[string]map[string]core.SecretBytes)
	c.mu.Unlock()

	informer := c.buildInformer()

	c.informer = informer

	err = c.initializeEventHandlers()
	if err != nil {
		return fmt.Errorf(`error occurred initializing event handlers: %w`, err)
	}

	return nil
}

// Run starts the SharedInformer.
func (c *Certificates) Run(ctx context.Context) error {
	slog.Info("Certificates::Run")

	if c.informer == nil {
		return fmt.Errorf(`initialize must be called before Run`)
	}

	c.informer.Run(ctx.Done())

	<-ctx.Done()

	return nil
}

func (c *Certificates) buildInformer() cache.SharedInformer {
	slog.Debug("Certificates::buildInformer")

	options := informers.WithNamespace(SecretsNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(c.k8sClient, 0, options)
	informer := factory.Core().V1().Secrets().Informer()

	return informer
}

func (c *Certificates) initializeEventHandlers() error {
	slog.Debug("Certificates::initializeEventHandlers")

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
	slog.Debug("Certificates::handleAddEvent")

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		slog.Error("Certificates::handleAddEvent: unable to cast object to Secret")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.certificates[secret.Name] = map[string]core.SecretBytes{}

	// Input from the secret comes in the form
	//   tls.crt: LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUVCVEN...
	//   tls.key: LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUV2Z0l...
	// Where the keys are `tls.crt` and `tls.key` and the values are []byte
	for k, v := range secret.Data {
		c.certificates[secret.Name][k] = core.SecretBytes(v)
	}

	slog.Debug("Certificates::handleAddEvent", slog.Int("certCount", len(c.certificates)))
}

func (c *Certificates) handleDeleteEvent(obj interface{}) {
	slog.Debug("Certificates::handleDeleteEvent")

	secret, ok := obj.(*corev1.Secret)
	if !ok {
		slog.Error("Certificates::handleDeleteEvent: unable to cast object to Secret")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.certificates[secret.Name] != nil {
		delete(c.certificates, secret.Name)
	}

	slog.Debug("Certificates::handleDeleteEvent", slog.Int("certCount", len(c.certificates)))
}

func (c *Certificates) handleUpdateEvent(_ interface{}, newValue interface{}) {
	slog.Debug("Certificates::handleUpdateEvent")

	secret, ok := newValue.(*corev1.Secret)
	if !ok {
		slog.Error("Certificates::handleUpdateEvent: unable to cast object to Secret")
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range secret.Data {
		c.certificates[secret.Name][k] = v
	}

	slog.Debug("Certificates::handleUpdateEvent", slog.Int("certCount", len(c.certificates)))
}

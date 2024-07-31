/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package configuration

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const (
	// ConfigMapsNamespace is the value used to filter the ConfigMaps Resource in the Informer.
	ConfigMapsNamespace = "nlk"

	// ConfigMapName is the name of the ConfigMap that contains the configuration for the application.
	ConfigMapName = "nlk-config"

	// ResyncPeriod is the value used to set the resync period for the Informer.
	ResyncPeriod = 0

	// NlkPrefix is used to determine if a Port definition should be handled and used to update a Border Server.
	// The Port name () must start with this prefix, e.g.:
	//   nlk-<my-upstream-name>
	NlkPrefix = ConfigMapsNamespace + "-"

	// PortAnnotationPrefix defines the prefix used when looking up a Port in the Service Annotations.
	// The value of the annotation determines which BorderServer implementation will be used.
	// See the documentation in the `application/application_constants.go` file for details.
	PortAnnotationPrefix = "nginxinc.io"

	// ServiceAnnotationMatchKey is the key name of the annotation in the application's config map
	// that identifies the ingress service whose events will be monitored.
	ServiceAnnotationMatchKey = "service-annotation-match"

	// DefaultServiceAnnotation is the default name of the ingress service whose events will be
	// monitored.
	DefaultServiceAnnotation = "nginxaas"
)

// WorkQueueSettings contains the configuration values needed by the Work Queues.
// There are two work queues in the application:
// 1. nlk-handler queue, used to move messages between the Watcher and the Handler.
// 2. nlk-synchronizer queue, used to move message between the Handler and the Synchronizer.
// The queues are NamedDelayingQueue objects that use an ItemExponentialFailureRateLimiter
// as the underlying rate limiter.
type WorkQueueSettings struct {
	// Name is the name of the queue.
	Name string

	// RateLimiterBase is the value used to calculate the exponential backoff rate limiter.
	// The formula is: RateLimiterBase * 2 ^ (num_retries - 1)
	RateLimiterBase time.Duration

	// RateLimiterMax limits the amount of time retries are allowed to be attempted.
	RateLimiterMax time.Duration
}

// HandlerSettings contains the configuration values needed by the Handler.
type HandlerSettings struct {
	// RetryCount is the number of times the Handler will attempt to process a message before giving up.
	RetryCount int

	// Threads is the number of threads that will be used to process messages.
	Threads int

	// WorkQueueSettings is the configuration for the Handler's queue.
	WorkQueueSettings WorkQueueSettings
}

// WatcherSettings contains the configuration values needed by the Watcher.
type WatcherSettings struct {
	// ServiceAnnotation is the annotation of the ingress service whose events the watcher should monitor.
	ServiceAnnotation string

	// ResyncPeriod is the value used to set the resync period for the underlying SharedInformer.
	ResyncPeriod time.Duration
}

// SynchronizerSettings contains the configuration values needed by the Synchronizer.
type SynchronizerSettings struct {
	// MaxMillisecondsJitter is the maximum number of milliseconds that will be applied when adding an event to the queue.
	MaxMillisecondsJitter int

	// MinMillisecondsJitter is the minimum number of milliseconds that will be applied when adding an event to the queue.
	MinMillisecondsJitter int

	// RetryCount is the number of times the Synchronizer will attempt to process a message before giving up.
	RetryCount int

	// Threads is the number of threads that will be used to process messages.
	Threads int

	// WorkQueueSettings is the configuration for the Synchronizer's queue.
	WorkQueueSettings WorkQueueSettings
}

// Settings contains the configuration values needed by the application.
type Settings struct {
	// Context is the context used to control the application.
	Context context.Context

	// NginxPlusHosts is a list of Nginx Plus hosts that will be used to update the Border Servers.
	NginxPlusHosts []string

	// TlsMode is the value used to determine which of the five TLS modes will be used to communicate
	// with the Border Servers (see: ../../docs/tls/README.md).
	TLSMode TLSMode

	// APIKey is the api key used to authenticate with the dataplane API.
	APIKey string

	// Certificates is the object used to retrieve the certificates and keys used to communicate with the Border Servers.
	Certificates *certification.Certificates

	// K8sClient is the Kubernetes client used to communicate with the Kubernetes API.
	K8sClient kubernetes.Interface

	// informer is the SharedInformer used to watch for changes to the ConfigMap .
	informer cache.SharedInformer

	// eventHandlerRegistration is the object used to track the event handlers with the SharedInformer.
	eventHandlerRegistration cache.ResourceEventHandlerRegistration

	// Handler contains the configuration values needed by the Handler.
	Handler HandlerSettings

	// Synchronizer contains the configuration values needed by the Synchronizer.
	Synchronizer SynchronizerSettings

	// Watcher contains the configuration values needed by the Watcher.
	Watcher WatcherSettings
}

// NewSettings creates a new Settings object with default values.
func NewSettings(ctx context.Context, k8sClient kubernetes.Interface) (*Settings, error) {
	// get base64 encoded version of raw api key set by user
	apiKey := base64.StdEncoding.EncodeToString([]byte(os.Getenv("NGINXAAS_DATAPLANE_API_KEY")))

	settings := &Settings{
		Context:      ctx,
		K8sClient:    k8sClient,
		TLSMode:      NoTLS,
		APIKey:       apiKey,
		Certificates: nil,
		Handler: HandlerSettings{
			RetryCount: 5,
			Threads:    1,
			WorkQueueSettings: WorkQueueSettings{
				RateLimiterBase: time.Second * 2,
				RateLimiterMax:  time.Second * 60,
				Name:            "nlk-handler",
			},
		},
		Synchronizer: SynchronizerSettings{
			MaxMillisecondsJitter: 750,
			MinMillisecondsJitter: 250,
			RetryCount:            5,
			Threads:               1,
			WorkQueueSettings: WorkQueueSettings{
				RateLimiterBase: time.Second * 2,
				RateLimiterMax:  time.Second * 60,
				Name:            "nlk-synchronizer",
			},
		},
		Watcher: WatcherSettings{
			ResyncPeriod:      0,
			ServiceAnnotation: DefaultServiceAnnotation,
		},
	}

	return settings, nil
}

// Initialize initializes the Settings object. Sets up a SharedInformer to watch for changes to the ConfigMap.
// This method must be called before the Run method.
func (s *Settings) Initialize() error {
	logrus.Info("Settings::Initialize")

	var err error

	certificates := certification.NewCertificates(s.Context, s.K8sClient)

	err = certificates.Initialize()
	if err != nil {
		return fmt.Errorf(`error occurred initializing certificates: %w`, err)
	}

	s.Certificates = certificates

	go certificates.Run() //nolint:errcheck

	logrus.Debug(">>>>>>>>>> Settings::Initialize: retrieving nlk-config ConfigMap")
	configMap, err := s.K8sClient.CoreV1().ConfigMaps(ConfigMapsNamespace).Get(
		s.Context, "nlk-config", metav1.GetOptions{},
	)
	if err != nil {
		return err
	}

	s.handleUpdateEvent(nil, configMap)
	logrus.Debug(">>>>>>>>>> Settings::Initialize: retrieved nlk-config ConfigMap")

	informer := s.buildInformer()

	s.informer = informer

	err = s.initializeEventListeners()
	if err != nil {
		return fmt.Errorf(`error occurred initializing event listeners: %w`, err)
	}

	return nil
}

// Run starts the SharedInformer and waits for the Context to be canceled.
func (s *Settings) Run() {
	logrus.Debug("Settings::Run")

	defer utilruntime.HandleCrash()

	go s.informer.Run(s.Context.Done())

	<-s.Context.Done()
}

func (s *Settings) buildInformer() cache.SharedInformer {
	options := informers.WithNamespace(ConfigMapsNamespace)
	factory := informers.NewSharedInformerFactoryWithOptions(s.K8sClient, ResyncPeriod, options)
	informer := factory.Core().V1().ConfigMaps().Informer()

	return informer
}

func (s *Settings) initializeEventListeners() error {
	logrus.Debug("Settings::initializeEventListeners")

	var err error

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc:    s.handleAddEvent,
		UpdateFunc: s.handleUpdateEvent,
		DeleteFunc: s.handleDeleteEvent,
	}

	s.eventHandlerRegistration, err = s.informer.AddEventHandler(handlers)
	if err != nil {
		return fmt.Errorf(`error occurred registering event handlers: %w`, err)
	}

	return nil
}

func (s *Settings) handleAddEvent(obj interface{}) {
	logrus.Debug("Settings::handleAddEvent")

	if _, yes := isOurConfig(obj); yes {
		s.handleUpdateEvent(nil, obj)
	}
}

func (s *Settings) handleDeleteEvent(obj interface{}) {
	logrus.Debug("Settings::handleDeleteEvent")

	if _, yes := isOurConfig(obj); yes {
		s.updateHosts([]string{})
	}
}

func (s *Settings) handleUpdateEvent(_ interface{}, newValue interface{}) {
	logrus.Debug("Settings::handleUpdateEvent")

	configMap, yes := isOurConfig(newValue)
	if !yes {
		return
	}

	hosts, found := configMap.Data["nginx-hosts"]
	if found {
		newHosts := s.parseHosts(hosts)
		s.updateHosts(newHosts)
	} else {
		logrus.Warnf("Settings::handleUpdateEvent: nginx-hosts key not found in ConfigMap")
	}

	tlsMode, err := validateTLSMode(configMap)
	if err != nil {
		// NOTE: the TLSMode defaults to NoTLS on startup, or the last known good value if previously set.
		logrus.Errorf(
			"Error with configured TLS Mode. TLS Mode has NOT been changed. The current mode is: '%v'. Error: %v. ",
			s.TLSMode, err,
		)
	} else {
		s.TLSMode = tlsMode
	}

	caCertificateSecretKey, found := configMap.Data["ca-certificate"]
	if found {
		s.Certificates.CaCertificateSecretKey = caCertificateSecretKey
		logrus.Debugf("Settings::handleUpdateEvent: ca-certificate: %s", s.Certificates.CaCertificateSecretKey)
	} else {
		s.Certificates.CaCertificateSecretKey = ""
		logrus.Warnf("Settings::handleUpdateEvent: ca-certificate key not found in ConfigMap")
	}

	clientCertificateSecretKey, found := configMap.Data["client-certificate"]
	if found {
		s.Certificates.ClientCertificateSecretKey = clientCertificateSecretKey
		logrus.Debugf("Settings::handleUpdateEvent: client-certificate: %s", s.Certificates.ClientCertificateSecretKey)
	} else {
		s.Certificates.ClientCertificateSecretKey = ""
		logrus.Warnf("Settings::handleUpdateEvent: client-certificate key not found in ConfigMap")
	}

	if serviceAnnotation, found := configMap.Data[ServiceAnnotationMatchKey]; found {
		s.Watcher.ServiceAnnotation = serviceAnnotation
	} else {
		s.Watcher.ServiceAnnotation = DefaultServiceAnnotation
	}
	logrus.Debugf("Settings::handleUpdateEvent: %s: %s", ServiceAnnotationMatchKey, s.Watcher.ServiceAnnotation)

	setLogLevel(configMap.Data["log-level"])

	logrus.Debugf("Settings::handleUpdateEvent: \n\tHosts: %v,\n\tSettings: %v ", s.NginxPlusHosts, configMap)
}

func validateTLSMode(configMap *corev1.ConfigMap) (TLSMode, error) {
	tlsConfigMode, tlsConfigModeFound := configMap.Data["tls-mode"]
	if !tlsConfigModeFound {
		return NoTLS, fmt.Errorf(`tls-mode key not found in ConfigMap`)
	}

	if tlsMode, tlsModeFound := TLSModeMap[tlsConfigMode]; tlsModeFound {
		return tlsMode, nil
	}

	return NoTLS, fmt.Errorf(`invalid tls-mode value: %s`, tlsConfigMode)
}

func (s *Settings) parseHosts(hosts string) []string {
	return strings.Split(hosts, ",")
}

func (s *Settings) updateHosts(hosts []string) {
	s.NginxPlusHosts = hosts
}

func isOurConfig(obj interface{}) (*corev1.ConfigMap, bool) {
	configMap, ok := obj.(*corev1.ConfigMap)
	return configMap, ok && configMap.Name == ConfigMapName && configMap.Namespace == ConfigMapsNamespace
}

func setLogLevel(logLevel string) {
	logrus.Debugf("Settings::setLogLevel: %s", logLevel)
	switch logLevel {
	case "panic":
		logrus.SetLevel(logrus.PanicLevel)

	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)

	case "error":
		logrus.SetLevel(logrus.ErrorLevel)

	case "warn":
		logrus.SetLevel(logrus.WarnLevel)

	case "info":
		logrus.SetLevel(logrus.InfoLevel)

	case "debug":
		logrus.SetLevel(logrus.DebugLevel)

	case "trace":
		logrus.SetLevel(logrus.TraceLevel)

	default:
		logrus.SetLevel(logrus.WarnLevel)
	}
}

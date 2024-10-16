/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package configuration

import (
	"encoding/base64"
	"fmt"
	"log/slog"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"

	"github.com/spf13/viper"
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
	// LogLevel is the user-specified log level. Defaults to warn.
	LogLevel string

	// NginxPlusHosts is a list of Nginx Plus hosts that will be used to update the Border Servers.
	NginxPlusHosts []string

	// TlsMode is the value used to determine which of the five TLS modes will be used to communicate
	// with the Border Servers (see: ../../docs/tls/README.md).
	TLSMode TLSMode

	// APIKey is the api key used to authenticate with the dataplane API.
	APIKey string

	// Certificates is the object used to retrieve the certificates and keys used to communicate with the Border Servers.
	Certificates *certification.Certificates

	// Handler contains the configuration values needed by the Handler.
	Handler HandlerSettings

	// Synchronizer contains the configuration values needed by the Synchronizer.
	Synchronizer SynchronizerSettings

	// Watcher contains the configuration values needed by the Watcher.
	Watcher WatcherSettings
}

// Read parses all the config and returns the values
func Read(configName, configPath string) (s Settings, err error) {
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)
	if err = v.ReadInConfig(); err != nil {
		return s, err
	}

	if err = v.BindEnv("NGINXAAS_DATAPLANE_API_KEY"); err != nil {
		return s, err
	}

	tlsMode := NoTLS
	if t, err := validateTLSMode(v.GetString("tls-mode")); err != nil {
		slog.Error("could not validate tls mode", "error", err)
	} else {
		tlsMode = t
	}

	serviceAnnotation := DefaultServiceAnnotation
	if sa := v.GetString(ServiceAnnotationMatchKey); sa != "" {
		serviceAnnotation = sa
	}

	return Settings{
		LogLevel:       v.GetString("log-level"),
		NginxPlusHosts: v.GetStringSlice("nginx-hosts"),
		TLSMode:        tlsMode,
		APIKey:         base64.StdEncoding.EncodeToString([]byte(v.GetString("NGINXAAS_DATAPLANE_API_KEY"))),
		Certificates: &certification.Certificates{
			CaCertificateSecretKey:     v.GetString("ca-certificate"),
			ClientCertificateSecretKey: v.GetString("client-certificate"),
		},
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
			ServiceAnnotation: serviceAnnotation,
		},
	}, nil
}

func validateTLSMode(tlsConfigMode string) (TLSMode, error) {
	if tlsMode, tlsModeFound := TLSModeMap[tlsConfigMode]; tlsModeFound {
		return tlsMode, nil
	}

	return NoTLS, fmt.Errorf(`invalid tls-mode value: %s`, tlsConfigMode)
}

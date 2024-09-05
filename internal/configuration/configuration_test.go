package configuration_test

import (
	"testing"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/certification"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"

	"github.com/stretchr/testify/require"
)

func TestConfiguration(t *testing.T) {
	t.Parallel()
	expectedSettings := configuration.Settings{
		LogLevel:       "warn",
		NginxPlusHosts: []string{"https://10.0.0.1:9000/api"},
		TLSMode:        configuration.NoTLS,
		Certificates: &certification.Certificates{
			CaCertificateSecretKey:     "fakeCAKey",
			ClientCertificateSecretKey: "fakeCertKey",
		},
		Handler: configuration.HandlerSettings{
			RetryCount: 5,
			Threads:    1,
			WorkQueueSettings: configuration.WorkQueueSettings{
				RateLimiterBase: time.Second * 2,
				RateLimiterMax:  time.Second * 60,
				Name:            "nlk-handler",
			},
		},
		Synchronizer: configuration.SynchronizerSettings{
			MaxMillisecondsJitter: 750,
			MinMillisecondsJitter: 250,
			RetryCount:            5,
			Threads:               1,
			WorkQueueSettings: configuration.WorkQueueSettings{
				RateLimiterBase: time.Second * 2,
				RateLimiterMax:  time.Second * 60,
				Name:            "nlk-synchronizer",
			},
		},
		Watcher: configuration.WatcherSettings{
			ResyncPeriod:      0,
			ServiceAnnotation: "fakeServiceMatch",
		},
	}

	settings, err := configuration.Read("test", "./testdata")
	require.NoError(t, err)
	require.Equal(t, expectedSettings, settings)
}

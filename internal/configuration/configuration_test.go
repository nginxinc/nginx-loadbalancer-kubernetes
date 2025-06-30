package configuration_test

import (
	"testing"
	"time"

	"github.com/nginxinc/kubernetes-nginx-ingress/internal/configuration"

	"github.com/stretchr/testify/require"
)

func TestConfiguration_Read(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		testFile         string
		expectedSettings configuration.Settings
	}{
		"one nginx plus host": {
			testFile: "one-nginx-host",
			expectedSettings: configuration.Settings{
				LogLevel:       "warn",
				NginxPlusHosts: []string{"https://10.0.0.1:9000/api"},
				SkipVerifyTLS:  false,
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
			},
		},
		"multiple nginx plus hosts": {
			testFile: "multiple-nginx-hosts",
			expectedSettings: configuration.Settings{
				LogLevel:       "warn",
				NginxPlusHosts: []string{"https://10.0.0.1:9000/api", "https://10.0.0.2:9000/api"},
				SkipVerifyTLS:  true,
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
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			settings, err := configuration.Read(tc.testFile, "./testdata")
			require.NoError(t, err)
			require.Equal(t, tc.expectedSettings, settings)
		})
	}
}

func TestConfiguration_TLS(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		tlsMode               string
		expectedSkipVerifyTLS bool
		expectedErr           bool
	}{
		"no input": {
			tlsMode:               "",
			expectedSkipVerifyTLS: false,
		},
		"no tls": {
			tlsMode:               "no-tls",
			expectedSkipVerifyTLS: true,
		},
		"skip verify tls": {
			tlsMode:               "skip-verify-tls",
			expectedSkipVerifyTLS: true,
		},
		"ca tls": {
			tlsMode:               "ca-tls",
			expectedSkipVerifyTLS: false,
		},
		"unexpected input": {
			tlsMode:     "unexpected-tls-mode",
			expectedErr: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			skipVerifyTLS, err := configuration.ValidateTLSMode(tc.tlsMode)
			if tc.expectedErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.expectedSkipVerifyTLS, skipVerifyTLS)
		})
	}
}

package config

import (
	"os"
	"testing"
)

func TestNewSettings(t *testing.T) {
	const EnvVarName = "NGINX_PLUS_HOST"
	const ExpectedValue = "https://locahost:443"

	defer os.Unsetenv(EnvVarName)
	os.Setenv(EnvVarName, ExpectedValue)

	settings, err := NewSettings()
	if err != nil {
		t.Fatalf(`Did not expect an error. %v`, err)
	}

	if settings.NginxPlusHost != ExpectedValue {
		t.Fatalf(`Expected %v to be %v`, settings.NginxPlusHost, ExpectedValue)
	}
}

func TestNewSettings_NginxPlusHostNotSet(t *testing.T) {
	_, err := NewSettings()
	if err == nil {
		t.Fatalf(`Expected an error that the NGINX_PLUS_HOST variable was not set`)
	}
}

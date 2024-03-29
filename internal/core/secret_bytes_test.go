/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestSecretBytesToString(t *testing.T) {
	t.Parallel()
	sensitive := SecretBytes([]byte("If you can see this we have a problem"))

	expected := "foo [REDACTED] bar"
	result := fmt.Sprintf("foo %v bar", sensitive)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestSecretBytesToJSON(t *testing.T) {
	t.Parallel()
	sensitive, _ := json.Marshal(SecretBytes([]byte("If you can see this we have a problem")))
	expected := `foo "[REDACTED]" bar`
	result := fmt.Sprintf("foo %v bar", string(sensitive))
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

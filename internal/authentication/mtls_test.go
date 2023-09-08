/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package authentication

import "testing"

func TestMtls(t *testing.T) {
	mtls, err := NewMtls()

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if mtls == nil {
		t.Fatalf(`mtls should not be nil`)
	}
}

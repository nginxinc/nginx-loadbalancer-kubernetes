/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package certification

import "testing"

func TestNewCertificate(t *testing.T) {
	cert, err := NewCertificate()

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if cert == nil {
		t.Fatalf(`cert should not be nil`)
	}
}

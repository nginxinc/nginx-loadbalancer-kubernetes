/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package certification

import (
	"context"
	"testing"
)

func TestNewCertificate(t *testing.T) {
	ctx := context.Background()

	cert, err := NewCertificates(ctx, nil)

	if err != nil {
		t.Fatalf(`Unexpected error: %v`, err)
	}

	if cert == nil {
		t.Fatalf(`cert should not be nil`)
	}
}

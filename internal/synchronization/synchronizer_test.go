// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"testing"
)

func TestNewNginxPlusSynchronizer(t *testing.T) {
	synchronizer, err := NewSynchronizer()
	if err != nil {
		t.Fatalf(`should have been no error, %v`, err)
	}

	if synchronizer == nil {
		t.Fatal("should have an Synchronizer instance")
	}
}

/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package core

import "testing"

func TestNewUpstreamServer(t *testing.T) {
	t.Parallel()
	host := "localhost"
	us := NewUpstreamServer(host)
	if us.Host != host {
		t.Errorf("NewUpstreamServer(%s) = %s; want %s", host, us.Host, host)
	}
}

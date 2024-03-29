/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

import "testing"

func TestCheck_LiveCheck(t *testing.T) {
	t.Parallel()
	check := LiveCheck{}
	if !check.Check() {
		t.Errorf("LiveCheck should return true")
	}
}

func TestCheck_ReadyCheck(t *testing.T) {
	t.Parallel()
	check := ReadyCheck{}
	if !check.Check() {
		t.Errorf("ReadyCheck should return true")
	}
}

func TestCheck_StartupCheck(t *testing.T) {
	t.Parallel()
	check := StartupCheck{}
	if !check.Check() {
		t.Errorf("StartupCheck should return true")
	}
}

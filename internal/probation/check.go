// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package probation

type Check interface {
	Check() bool
}

type LiveCheck struct {
}

type ReadyCheck struct {
}

type StartupCheck struct {
}

func (l *LiveCheck) Check() bool {
	return true
}

func (r *ReadyCheck) Check() bool {
	return true
}

func (s *StartupCheck) Check() bool {
	return true
}

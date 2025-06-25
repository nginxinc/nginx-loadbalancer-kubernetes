/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package probation

// Check defines a single method that can be implemented for various health checks.
type Check interface {
	Check() bool
}

// LiveCheck is a check that can be used for the k8s "livez" endpoint.
type LiveCheck struct {
}

// ReadyCheck is a check that can be used for the k8s "readyz" endpoint.
type ReadyCheck struct {
}

// StartupCheck is a check that can be used for the k8s "startupz" endpoint.
type StartupCheck struct {
}

// Check implements the Check interface for the LiveCheck type.
func (l *LiveCheck) Check() bool {
	return true
}

// Check implements the Check interface for the ReadyCheck type.
func (r *ReadyCheck) Check() bool {
	return true
}

// Check implements the Check interface for the StartupCheck type.
func (s *StartupCheck) Check() bool {
	return true
}

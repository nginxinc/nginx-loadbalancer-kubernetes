// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package mocks

type MockCheck struct {
	result bool
}

func NewMockCheck(result bool) *MockCheck {
	return &MockCheck{result: result}
}

func (m *MockCheck) Check() bool {
	return m.result
}

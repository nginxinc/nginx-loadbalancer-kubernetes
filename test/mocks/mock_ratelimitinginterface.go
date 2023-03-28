/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import "time"

type MockRateLimiter struct {
	items []interface{}
}

func (m *MockRateLimiter) Add(_ interface{}) {
}

func (m *MockRateLimiter) Len() int {
	return len(m.items)
}

func (m *MockRateLimiter) Get() (item interface{}, shutdown bool) {
	if len(m.items) > 0 {
		item = m.items[0]
		m.items = m.items[1:]
		return item, false
	}
	return nil, false
}

func (m *MockRateLimiter) Done(_ interface{}) {
}

func (m *MockRateLimiter) ShutDown() {
}

func (m *MockRateLimiter) ShutDownWithDrain() {
}

func (m *MockRateLimiter) ShuttingDown() bool {
	return true
}

func (m *MockRateLimiter) AddAfter(item interface{}, _ time.Duration) {
	m.items = append(m.items, item)
}

func (m *MockRateLimiter) AddRateLimited(item interface{}) {
	m.items = append(m.items, item)
}

func (m *MockRateLimiter) Forget(_ interface{}) {

}

func (m *MockRateLimiter) NumRequeues(_ interface{}) int {
	return 0
}

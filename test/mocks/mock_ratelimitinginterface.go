/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package mocks

import "time"

type MockRateLimiter[T any] struct {
	items []T
}

func (m *MockRateLimiter[T]) Add(_ T) {
}

func (m *MockRateLimiter[T]) Len() int {
	return len(m.items)
}

func (m *MockRateLimiter[T]) Get() (item T, shutdown bool) {
	if len(m.items) > 0 {
		item = m.items[0]
		m.items = m.items[1:]
		return item, false
	}
	return item, false
}

func (m *MockRateLimiter[T]) Done(_ T) {
}

func (m *MockRateLimiter[T]) ShutDown() {
}

func (m *MockRateLimiter[T]) ShutDownWithDrain() {
}

func (m *MockRateLimiter[T]) ShuttingDown() bool {
	return true
}

func (m *MockRateLimiter[T]) AddAfter(item T, _ time.Duration) {
	m.items = append(m.items, item)
}

func (m *MockRateLimiter[T]) AddRateLimited(item T) {
	m.items = append(m.items, item)
}

func (m *MockRateLimiter[T]) Forget(_ T) {
}

func (m *MockRateLimiter[T]) NumRequeues(_ T) int {
	return 0
}

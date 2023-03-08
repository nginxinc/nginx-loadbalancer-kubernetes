// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package mocks

import "net/http"

type MockResponseWriter struct {
	body []byte
}

func NewMockResponseWriter() *MockResponseWriter {
	return &MockResponseWriter{}
}

func (m *MockResponseWriter) Header() http.Header {
	return nil
}

func (m *MockResponseWriter) Write(body []byte) (int, error) {
	m.body = append(m.body, body...)
	return len(m.body), nil
}

func (m *MockResponseWriter) WriteHeader(int) {

}

func (m *MockResponseWriter) Body() []byte {
	return m.body
}

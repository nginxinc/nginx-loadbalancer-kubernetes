// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package synchronization

import (
	"math/rand"
	"time"
)

var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var number = []byte("0123456789")
var alphaNumeric = append(charset, number...)

// RandomString where n is the length of random string we want to generate
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))]
	}
	return string(b)
}

func RandomMilliseconds(min, max int) time.Duration {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := randomizer.Intn(max-min) + min

	return time.Millisecond * time.Duration(random)
}

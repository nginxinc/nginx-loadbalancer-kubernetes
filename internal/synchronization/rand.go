/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package synchronization

import (
	// Try using crpyto if needed.
	"math/rand"
	"time"
)

// charset contains all characters that can be used in random string generation
var charset = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// number contains all numbers that can be used in random string generation
var number = []byte("0123456789")

// alphaNumeric contains all characters and numbers that can be used in random string generation
var alphaNumeric = append(charset, number...)

// RandomString where n is the length of random string we want to generate
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		// randomly select 1 character from given charset
		b[i] = alphaNumeric[rand.Intn(len(alphaNumeric))] //nolint:gosec
	}
	return string(b)
}

// RandomMilliseconds returns a random duration between min and max milliseconds
func RandomMilliseconds(min, max int) time.Duration {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec
	random := randomizer.Intn(max-min) + min

	return time.Millisecond * time.Duration(random)
}

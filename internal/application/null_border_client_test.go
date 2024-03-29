/*
 * Copyright (c) 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package application

import "testing"

func TestNullBorderClient_Delete(t *testing.T) {
	t.Parallel()
	client := NullBorderClient{}
	err := client.Delete(nil)
	if err != nil {
		t.Errorf(`expected no error deleting border client, got: %v`, err)
	}
}

func TestNullBorderClient_Update(t *testing.T) {
	t.Parallel()
	client := NullBorderClient{}
	err := client.Update(nil)
	if err != nil {
		t.Errorf(`expected no error updating border client, got: %v`, err)
	}
}

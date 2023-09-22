/*
 * Copyright 2023 F5 Inc. All rights reserved.
 * Use of this source code is governed by the Apache License that can be found in the LICENSE file.
 */

package certification

type Certificates struct {
	caCertificate     string
	clientCertificate string
	clientKey         string
}

func NewCertificates() (*Certificates, error) {
	return &Certificates{}, nil
}

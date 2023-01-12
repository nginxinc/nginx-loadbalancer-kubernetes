// Copyright 2023 f5 Inc. All rights reserved.
// Use of this source code is governed by the Apache
// license that can be found in the LICENSE file.

package nginxplus

import (
	"fmt"
	"github.com/nginxinc/kubernetes-nginx-ingress/internal/translation"
	"testing"
)

// NOTE: this is just a placeholoder to ensure each implementation conforms to the translation.Translator interface.
func TestResults(t *testing.T) {
	ct := CreatedTranslator{}
	dt := UpdatedTranslator{}
	ut := DeletedTranslator{}

	printResult(ct)
	printResult(dt)
	printResult(ut)
}

func printResult(t translation.Translator) {
	i, err := t.Translate()
	if err != nil {
		fmt.Printf(`there was an error: %v`, err)
	}
	fmt.Println("  ")
	fmt.Printf(`success! %v`, i)
}

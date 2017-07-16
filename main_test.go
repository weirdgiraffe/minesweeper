//
// main_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"strings"
	"testing"
)

func TestNewField(t *testing.T) {
	tc := []struct {
		inText      string
		expectError bool
	}{
		{"4 4\n..*.\n*...\n....\n.**.\n", false},
		{"-1 4\n..*.\n*...\n....\n.**.\n", true},
		{"4 101\n..*.\n*...\n....\n.**.\n", true},
		{"4 4\n..*.\n*...\n.A..\n.**.\n", true},
		{"4 4\n..*.\n*............\n....\n.**.\n", true},
	}
	for i := range tc {
		field, err := NewField(strings.NewReader(tc[i].inText))
		if err != nil {
			if !tc[i].expectError {
				t.Fatal("'%s' Unexpected error: %v", tc[i].inText, err)
			}
			continue
		}
		if field.String() != tc[i].inText {
			t.Fatal(
				"Field not match:\nExpected\n'%s'Have\n'%s'\n",
				tc[i].inText, field.String(),
			)
		}
	}
}

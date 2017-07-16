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
		outText     string
		expectError bool
	}{
		{
			"0 0\n",
			"",
			false,
		},
		{
			"4 4\n..*.\n*...\n....\n.**.\n",
			"12*1\n*211\n2321\n1**1\n",
			false,
		},
		{
			"4 4\n*...\n....\n.*..\n....\n",
			"*100\n2210\n1*10\n1110\n",
			false,
		},
		{
			"3 5\n**...\n.....\n.*...\n",
			"**100\n33200\n1*100\n",
			false,
		},
		{
			"-1 4\n..*.\n*...\n....\n.**.\n",
			"",
			true,
		},
		{
			"4 101\n..*.\n*...\n....\n.**.\n",
			"",
			true,
		},
		{
			"4 4\n..*.\n*...\n.A..\n.**.\n",
			"",
			true,
		},
		{
			"4 4\n..*.\n*............\n....\n.**.\n",
			"",
			true,
		},
	}
	for i := range tc {
		field, err := NewField(strings.NewReader(tc[i].inText))
		if err != nil {
			if !tc[i].expectError {
				t.Fatalf("\n%s Unexpected error: %v", tc[i].inText, err)
			}
			continue
		}
		field.Enumerate()
		if field.String() != tc[i].outText {
			t.Fatalf(
				"Field not match:\nInput\n%s\nExpected\n%s\nHave\n%s\n",
				tc[i].inText,
				tc[i].outText,
				field.String(),
			)
		}
	}
}

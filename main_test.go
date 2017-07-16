//
// main_test.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"bytes"
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
			"Field #1:\n12*1\n*211\n2321\n1**1\n\n",
			false,
		},
		{
			"4 4\n*...\n....\n.*..\n....\n",
			"Field #1:\n*100\n2210\n1*10\n1110\n\n",
			false,
		},
		{
			"3 5\n**...\n.....\n.*...\n",
			"Field #1:\n**100\n33200\n1*100\n\n",
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
		buf := new(bytes.Buffer)
		r := strings.NewReader(tc[i].inText)
		err := DoMinesweeper(r, buf)
		if err != nil {
			if !tc[i].expectError {
				t.Fatalf("TestCase #%d\n%s Unexpected error: %v", i, tc[i].inText, err)
			}
			continue
		}
		if buf.String() != tc[i].outText {
			t.Fatalf(
				"TestCase #%d Field not match:\nInput\n%s\nExpected\n%s\nHave\n%s\n",
				i, tc[i].inText, tc[i].outText, buf.String(),
			)
		}
	}
}

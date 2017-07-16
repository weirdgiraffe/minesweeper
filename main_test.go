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

func TestReadDimensions(t *testing.T) {
	tc := []struct {
		in       string
		hasError bool
		en, em   int
	}{
		{"10 20", false, 10, 20},
		{"-1 20", true, 0, 0},
		{"0 20", true, 0, 0},
		{"1 120", true, 0, 0},
		{"a b", true, 0, 0},
		{"1 2 3 4", true, 0, 0},
	}
	for i := range tc {
		n, m, err := ReadDimensions(strings.NewReader(tc[i].in))
		if err != nil {
			if !tc[i].hasError {
				t.Fatal("'%s' Unexpected error: %v", tc[i].in, err)
			}
			t.Log(err)
			continue
		}
		if n != tc[i].en {
			t.Errorf("'%s' Unexpected n: %d != %d", tc[i].in, tc[i].en, n)
		}
		if m != tc[i].em {
			t.Errorf("'%s' Unexpected m: %d != %d", tc[i].in, tc[i].em, m)
		}
	}
}

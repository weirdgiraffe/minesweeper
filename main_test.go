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

func TestReadField(t *testing.T) {
	tc := []struct {
		n, m     int
		in       string
		hasError bool
		eField   [][]bool
	}{
		{
			4, 4, "..*.\n*...\n....\n.**.\n", false,
			[][]bool{
				[]bool{false, false, true, false},
				[]bool{true, false, false, false},
				[]bool{false, false, false, false},
				[]bool{false, true, true, false},
			},
		},
		{
			4, 4, ".A*.\n*...\n....\n.**.\n", true,
			[][]bool{},
		},
		{
			4, 4, "..*.............\n*...\n....\n.**.\n", true,
			[][]bool{},
		},
		{
			4, 4, "....\n..*.\n*...\n....\n.**.\n", true,
			[][]bool{},
		},
	}
	for i := range tc {
		field, err := ReadField(strings.NewReader(tc[i].in), tc[i].n, tc[i].m)
		if err != nil {
			if !tc[i].hasError {
				t.Fatal("%dx%d '%s' Unexpected error: %v", tc[i].n, tc[i].m, tc[i].in, err)
			}
			continue
		}
		for j := range tc[i].eField {
			for k := range tc[i].eField[i] {
				if field[j][k] != tc[i].eField[j][k] {
					t.Errorf(
						"%dx%d '%s' Field [%d][%d] doesnt match: %v != %v",
						tc[i].n, tc[i].m, j, k, tc[i].eField[j][k], field[j][k],
					)
				}
			}
		}
	}

}

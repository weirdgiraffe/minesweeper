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
		eField   [][]byte
	}{
		{
			4, 4, "..*.\n*...\n....\n.**.\n", false,
			[][]byte{
				[]byte{0, 0, Bomb, 0},
				[]byte{Bomb, 0, 0, 0},
				[]byte{0, 0, 0, 0},
				[]byte{0, Bomb, Bomb, 0},
			},
		},
		{
			4, 4, ".A*.\n*...\n....\n.**.\n", true,
			[][]byte{},
		},
		{
			4, 4, "..*.............\n*...\n....\n.**.\n", true,
			[][]byte{},
		},
		{
			4, 4, "....\n..*.\n*...\n....\n.**.\n", true,
			[][]byte{},
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
						tc[i].n, tc[i].m, tc[i].in, j, k, tc[i].eField[j][k], field[j][k],
					)
				}
			}
		}
	}

}

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

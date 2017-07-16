//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"bufio"
	"fmt"
	"io"
)

func main() {}

const Bomb = 255

type Field struct {
	cell [][]byte
}

func NewField(r io.Reader) (*Field, error) {
	f := &Field{}
	err := f.readDimensions(r)
	if err != nil {
		return nil, err
	}
	err = f.readCells(r)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *Field) String() string {
	ret := fmt.Sprintf("%d %d\n", len(f.cell), len(f.cell[0]))
	for i := range f.cell {
		for j := range f.cell {
			if f.cell[i][j] == Bomb {
				ret += "*"
			} else {
				ret += "."
			}
		}
		ret += "\n"
	}
	return ret
}

func (f *Field) readDimensions(r io.Reader) (err error) {
	var n, m int
	_, err = fmt.Fscanf(r, "%d %d\n", &n, &m)
	if err != nil {
		return
	}
	if n <= 0 || m > 100 {
		err = fmt.Errorf("Wrong dimensions: allowed n > 0 and m <= 100")
		return
	}
	f.cell = make([][]byte, n)
	for i := range f.cell {
		f.cell[i] = make([]byte, m)
	}
	return nil
}

func (f *Field) readCells(r io.Reader) (err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() && i < len(f.cell) {
		line := scanner.Text()
		if len(line) != len(f.cell[0]) {
			return fmt.Errorf(
				"Bad field line '%s'. Must have exactly %d symbols",
				line, len(f.cell[0]),
			)
		}
		for j := range line {
			switch line[j] {
			case '.':
				f.cell[i][j] = 0
			case '*':
				f.cell[i][j] = Bomb
			default:
				return fmt.Errorf(
					"'%s' Bad symbol %c: '.', '*' are allowed",
					line, line[j],
				)
			}
		}
		i++
	}
	err = scanner.Err()
	if err != nil {
		return err
	}
	if i != len(f.cell) {
		return fmt.Errorf(
			"Row count doesn't match format: expected %d, got %d",
			len(f.cell), i,
		)
	}
	return nil
}

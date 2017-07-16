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

type Field struct {
	cell [][]byte
}

func NewField(r io.Reader) (*Field, error) {
	n, m, err := ReadDimensions(r)
	if err != nil {
		return nil, err
	}
	f := &Field{}
	f.cell, err = ReadField(r, n, m)
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

func ReadDimensions(r io.Reader) (n int, m int, err error) {
	_, err = fmt.Fscanf(r, "%d %d\n", &n, &m)
	if err != nil {
		return
	}
	if n <= 0 || m > 100 {
		err = fmt.Errorf("Wrong dimensions: allowed n > 0 and m <= 100")
		return
	}
	return n, m, nil
}

const Bomb = 255

func ReadField(r io.Reader, row, col int) (field [][]byte, err error) {
	field = make([][]byte, row)
	for i := 0; i < row; i++ {
		field[i] = make([]byte, col)
	}
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() && i < row {
		line := scanner.Text()
		if len(line) != col {
			err = fmt.Errorf(
				"Bad field line '%s'. Must have exactly %d symbols",
				line, col,
			)
			return
		}
		for j := range line {
			switch line[j] {
			case '.':
				field[i][j] = 0
			case '*':
				field[i][j] = Bomb
			default:
				err = fmt.Errorf("Bad symbol in line '%s'. Only '.' and '*' are allowed", line)
				return
			}
		}
		i++
	}
	err = scanner.Err()
	if err != nil {
		return
	}
	if i != row {
		err = fmt.Errorf("Row count doesn't match format: expected %d, got %d", row, i)
		return
	}
	return field, nil
}

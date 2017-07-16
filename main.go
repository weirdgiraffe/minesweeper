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
	"strconv"
)

func main() {}

const Bomb = 255

type Field struct {
	cell       [][]byte
	rows, cols int
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

func (f *Field) bombsAround(i, j int, cellInclusive bool) (count byte) {
	if j != 0 {
		if f.cell[i][j-1] == Bomb {
			count++
		}
	}
	if j != f.cols-1 {
		if f.cell[i][j+1] == Bomb {
			count++
		}
	}
	if cellInclusive {
		if f.cell[i][j] == Bomb {
			count++
		}
	}
	return count
}

func (f *Field) Enumerate() {
	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			cell := &f.cell[i][j]
			if *cell == Bomb {
				continue
			}
			*cell += f.bombsAround(i, j, false)
			if i != 0 {
				*cell += f.bombsAround(i-1, j, true)
			}
			if i != f.rows-1 {
				*cell += f.bombsAround(i+1, j, true)
			}
		}
	}
}

func (f *Field) String() string {
	ret := ""
	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			if f.cell[i][j] == Bomb {
				ret += "*"
			} else {
				ret += strconv.Itoa(int(f.cell[i][j]))
			}
		}
		ret += "\n"
	}
	return ret
}

func (f *Field) readDimensions(r io.Reader) (err error) {
	_, err = fmt.Fscanf(r, "%d %d\n", &f.rows, &f.cols)
	if err != nil {
		return
	}
	if f.rows <= 0 || f.cols > 100 {
		err = fmt.Errorf("Wrong dimensions: allowed n > 0 and m <= 100")
		return
	}
	f.cell = make([][]byte, f.rows)
	for i := range f.cell {
		f.cell[i] = make([]byte, f.cols)
	}
	return nil
}

func (f *Field) readCells(r io.Reader) (err error) {
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() && i < f.rows {
		line := scanner.Text()
		if len(line) != f.cols {
			return fmt.Errorf(
				"Bad field line '%s'. Must have exactly %d symbols",
				line, f.cols,
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
	if i != f.rows {
		return fmt.Errorf(
			"Row count doesn't match format: expected %d, got %d",
			f.rows, i,
		)
	}
	return nil
}

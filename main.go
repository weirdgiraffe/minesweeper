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
	"os"
)

func main() {
	DoMinesweeper(os.Stdin, os.Stdout)
}

func DoMinesweeper(r io.Reader, w io.Writer) (err error) {
	for i := 1; ; i++ {
		f := &Field{}
		err = f.readDimensions(r)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return
		}
		err = f.readCells(r)
		if err != nil {
			return
		}
		f.Enumerate()
		fmt.Fprintf(w, "Field #%d:\n%s\n", i, f.String())
	}
}

const Bomb = '*'
const Unknown = '.'

type Field struct {
	cell       [][]byte
	rows, cols int
}

func (f *Field) bombsAround(i, j int, cellInclusive bool) (count byte) {
	if i < 0 || i >= f.rows {
		return 0
	}
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
			if f.cell[i][j] == Bomb {
				continue
			}
			if f.cell[i][j] == Unknown {
				f.cell[i][j] = '0'
			}
			f.cell[i][j] += f.bombsAround(i, j, false)
			f.cell[i][j] += f.bombsAround(i-1, j, true)
			f.cell[i][j] += f.bombsAround(i+1, j, true)
		}
	}
}

func (f *Field) String() string {
	ret := ""
	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			ret += string(f.cell[i][j])
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
	if f.rows == 0 && f.cols == 0 {
		return io.EOF
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
				f.cell[i][j] = Unknown
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

//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

func main() {
	flag.Usage = func() {
		desc := "Minesweeper reads Field with BombCells from stdin,\n" +
			"find bombs, and prints field with number of bombs\n" +
			"in cells. Last field should have dimensions 0 0.\n" +
			"Example:\n\n" +
			"Input:\n3 5\n**...\n.....\n.*...\n0 0\n\n" +
			"Output:\nField #1:\n*100\n2210\n1*10\n1110\n"
		fmt.Fprintf(os.Stderr, desc)
	}
	flag.Parse() // to output help and usage
	err := DoMinesweeper(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

// DoMinesweeper reads input from r, find bombs and write results to w
func DoMinesweeper(r io.Reader, w io.Writer) (err error) {
	for i := 1; i < math.MaxInt32; i++ {
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
		f.FindBombs()
		fmt.Fprintf(w, "Field #%d:\n%s\n", i, f)
	}
	return fmt.Errorf("Too many inputs. Minesweeper can process up to %d inputs", math.MaxInt32)
}

const BombCell = '*'
const SafeCell = '.'

type Field struct {
	cell       [][]byte
	rows, cols int
}

// bombsAroundCell find all bombs in the same row around cell defined by i,j.
// If includeCellItself is true, then cell defined by i,j will be checked also.
// All boundaries are checked and this method produce neither panic nor errors
func (f *Field) bombsAroundCell(i, j int, includeCellItself bool) (count byte) {
	if i < 0 || i >= f.rows {
		return 0
	}
	if j != 0 {
		if f.cell[i][j-1] == BombCell {
			count++
		}
	}
	if j != f.cols-1 {
		if f.cell[i][j+1] == BombCell {
			count++
		}
	}
	if includeCellItself {
		if f.cell[i][j] == BombCell {
			count++
		}
	}
	return count
}

func (f *Field) FindBombs() {
	for i := 0; i < f.rows; i++ {
		for j := 0; j < f.cols; j++ {
			if f.cell[i][j] == BombCell {
				continue
			}
			if f.cell[i][j] == SafeCell {
				f.cell[i][j] = '0'
			}
			// find bombs on the current row
			f.cell[i][j] += f.bombsAroundCell(i, j, false)
			// find bombs on the row downside
			f.cell[i][j] += f.bombsAroundCell(i-1, j, true)
			// find bombs on the row upside
			f.cell[i][j] += f.bombsAroundCell(i+1, j, true)
		}
	}
}

// String implement Stringer interface
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

// readDimensions reads n and m integers from r and check input
func (f *Field) readDimensions(r io.Reader) (err error) {
	_, err = fmt.Fscanf(r, "%d %d\n", &f.rows, &f.cols)
	if err != nil {
		return
	}
	if f.rows == 0 && f.cols == 0 {
		return io.EOF
	}
	if f.rows < 0 || 100 < f.rows {
		err = fmt.Errorf("Wrong Field dimensions: allowed 0 < n <= 100")
		return
	}
	if f.cols < 0 || 100 < f.cols {
		err = fmt.Errorf("Wrong Field dimensions: allowed 0 < m <= 100")
		return
	}
	f.cell = make([][]byte, f.rows)
	for i := range f.cell {
		f.cell[i] = make([]byte, f.cols)
	}
	return nil
}

// readCells reads field cells from r and check input
func (f *Field) readCells(r io.Reader) (err error) {
	if f.rows == 0 && f.cols == 0 {
		return fmt.Errorf("Must define Field dimensions first")
	}
	scanner := bufio.NewScanner(r)
	i := 0
	for scanner.Scan() && i < f.rows {
		line := scanner.Text()
		if len(line) != f.cols {
			return fmt.Errorf(
				"Bad field line '%s'. Must be exactly %d symbols long",
				line, f.cols,
			)
		}
		for j := range line {
			switch line[j] {
			case '.':
				f.cell[i][j] = SafeCell
			case '*':
				f.cell[i][j] = BombCell
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

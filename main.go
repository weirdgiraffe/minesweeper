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

const maxRows = 100
const maxCols = 100
const maxFields = math.MaxInt32

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
	scanner := bufio.NewScanner(r)
	for i := 1; i < maxFields; i++ {
		f := &Field{}
		err = f.readDimensions(scanner)
		if err != nil {
			if _, ok := err.(*stopReadingError); ok {
				return nil
			}
			return
		}
		err = f.readCells(scanner)
		if err != nil {
			return
		}
		f.FindBombs()
		fmt.Fprintf(w, "Field #%d:\n%s\n", i, f)
	}
	return fmt.Errorf("Too many inputs. Minesweeper can process up to %d inputs", maxFields)
}

// StopReadingError is an error to indicate the end of fields processing
// Is needed to make difference between io.EOF and end of processing
type stopReadingError struct{}

func (e stopReadingError) Error() string {
	return "Minesweeper: stop reading"
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
func (f *Field) readDimensions(scanner *bufio.Scanner) (err error) {
	if scanner.Scan() == false {
		err = scanner.Err()
		if err != nil {
			return
		}
		return io.EOF
	}
	_, err = fmt.Sscanf(scanner.Text(), "%d %d\n", &f.rows, &f.cols)
	if err != nil {
		return
	}
	if f.rows == 0 && f.cols == 0 {
		return &stopReadingError{}
	}
	if f.rows < 0 || maxRows < f.rows {
		err = fmt.Errorf("Wrong Field dimensions: allowed 0 < n <= %d", maxRows)
		return
	}
	if f.cols < 0 || maxCols < f.cols {
		err = fmt.Errorf("Wrong Field dimensions: allowed 0 < m <= %d", maxCols)
		return
	}
	f.cell = make([][]byte, f.rows)
	for i := range f.cell {
		f.cell[i] = make([]byte, f.cols)
	}
	return nil
}

// readCells reads field cells from r and check input
func (f *Field) readCells(scanner *bufio.Scanner) (err error) {
	if f.rows == 0 && f.cols == 0 {
		return fmt.Errorf("Must define Field dimensions first")
	}
	for i := 0; i < f.rows; i++ {
		if scanner.Scan() == false {
			err = scanner.Err()
			if err != nil {
				return err
			}
			return io.EOF
		}
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
	}
	return nil
}

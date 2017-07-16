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

func ReadField(r io.Reader, row, col int) (field [][]bool, err error) {
	field = make([][]bool, row)
	for i := 0; i < row; i++ {
		field[i] = make([]bool, col)
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
				field[i][j] = false
			case '*':
				field[i][j] = true
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

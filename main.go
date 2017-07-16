//
// main.go
// Copyright (C) 2017 weirdgiraffe <giraffe@cyberzoo.xyz>
//
// Distributed under terms of the MIT license.
//

package main

import (
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

//
// Copyright (C) 2022 - 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fogfish/golem/pipe"
)

const (
	cap = 1
)

func main() {
	ctx, close := context.WithCancel(context.Background())

	// Generate sequence of integers
	ints := pipe.StdErr(pipe.Unfold(ctx, cap, 0,
		func(x int) (int, error) { return x + 1, nil },
	))

	// Limit sequence of integers
	ints10 := pipe.TakeWhile(ctx, ints,
		func(x int) bool { return x <= 10 },
	)

	// Calculate squares
	sqrt := pipe.StdErr(pipe.Map(ctx, ints10,
		func(x int) (int, error) { return x * x, nil },
	))

	// Numbers to string
	vals := pipe.StdErr(pipe.Map(ctx, sqrt,
		func(x int) (string, error) { return strconv.Itoa(x), nil },
	))

	// Output strings
	<-pipe.ForEach(ctx, vals,
		func(x string) { fmt.Printf("==> %s\n", x) },
	)

	close()
}

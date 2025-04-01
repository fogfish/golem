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

	"github.com/fogfish/golem/pipe/v2"
)

const (
	cap = 1
)

func main() {
	ctx, close := context.WithCancel(context.Background())

	// Generate sequence of integers
	ints := pipe.StdErr(pipe.Unfold(ctx, cap, 0,
		pipe.Pure(func(x int) int { return x + 1 }),
	))

	// Limit sequence of integers
	ints10 := pipe.TakeWhile(ctx, ints,
		pipe.Pure(func(x int) bool { return x <= 10 }),
	)

	// Calculate squares
	sqrt := pipe.StdErr(pipe.Map(ctx, ints10,
		pipe.Pure(func(x int) int { return x * x }),
	))

	// Numbers to string
	vals := pipe.StdErr(pipe.Map(ctx, sqrt,
		pipe.Pure(strconv.Itoa),
	))

	// Output strings
	<-pipe.ForEach(ctx, vals,
		pipe.Pure(
			func(x string) string {
				fmt.Printf("==> %s\n", x)
				return x
			},
		),
	)

	close()
}

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
	"time"

	"github.com/fogfish/golem/pipe"
)

const (
	fastProducer = 10000
	cap          = 1
)

func main() {
	ctx, close := context.WithCancel(context.Background())

	// Generate sequence of integers
	fast := pipe.StdErr(pipe.Unfold(ctx, fastProducer, 0,
		func(x int) (int, error) { return x + 1, nil },
	))

	// Throttle the "fast" pipe
	slow := pipe.Throttling(ctx, fast, 1, 100*time.Millisecond)

	// Numbers to string
	vals := pipe.StdErr(pipe.Map(ctx, slow,
		func(x int) (string, error) { return strconv.Itoa(x), nil },
	))

	// Output strings
	<-pipe.ForEach(ctx, vals,
		func(x string) {
			fmt.Printf("==> %s | %s\n", time.Now().Format(time.StampMilli), x)
		},
	)

	close()
}

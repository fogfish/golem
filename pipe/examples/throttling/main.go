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

	"github.com/fogfish/golem/pipe/v2"
)

const (
	fastProducer = 10000
	cap          = 1
)

func main() {
	ctx, close := context.WithCancel(context.Background())

	// Generate sequence of integers
	fast := pipe.StdErr(pipe.Unfold(ctx, fastProducer, 0,
		pipe.Pure(func(x int) int { return x + 1 }),
	))

	// Throttle the "fast" pipe
	slow := pipe.Throttling(ctx, fast, 1, 100*time.Millisecond)

	// Numbers to string
	vals := pipe.StdErr(pipe.Map(ctx, slow,
		pipe.Pure(strconv.Itoa),
	))

	// Output strings
	<-pipe.ForEach(ctx, vals,
		pipe.Pure(
			func(x string) string {
				fmt.Printf("==> %s | %s\n", time.Now().Format(time.StampMilli), x)
				return x
			},
		),
	)

	close()
}

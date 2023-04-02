//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pure_test

import (
	"testing"

	"github.com/fogfish/golem/pure"
	"github.com/fogfish/it"
)

func f(x int) int { return x + 1 }

func TestPipe(t *testing.T) {
	for n, f := range []func(int) int{
		pure.Pipe(f, f),
		pure.Pipe3(f, f, f),
		pure.Pipe4(f, f, f, f),
		pure.Pipe5(f, f, f, f, f),
		pure.Pipe6(f, f, f, f, f, f),
		pure.Pipe7(f, f, f, f, f, f, f),
		pure.Pipe8(f, f, f, f, f, f, f, f),
		pure.Pipe9(f, f, f, f, f, f, f, f, f),
		pure.Pipe10(f, f, f, f, f, f, f, f, f, f),
		pure.Pipe11(f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe12(f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe13(f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe14(f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe15(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe16(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe17(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe18(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe19(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
		pure.Pipe20(f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f, f),
	} {
		it.Ok(t).If(f(0)).Equal(n + 2)
	}
}

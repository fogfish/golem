//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package maptest

import (
	"testing"

	"github.com/fogfish/golem/maplike"
)

func TestMap[F_ maplike.Kind[A, B], A, B any](t *testing.T, mapT maplike.Map[F_, A, B], seed F_) {
	t.Run("Map.HKT", func(t *testing.T) {
		seed.HKT1(maplike.Type(nil))
		seed.HKT2(*new(A), *new(B))
	})

	t.Run("Map.Put", func(t *testing.T) {

	})

}

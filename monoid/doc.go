//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

//go:generate golem -lib -T int -generic github.com/fogfish/golem/monoid/monoid.go
//go:generate golem -lib -T string -generic github.com/fogfish/golem/monoid/monoid.go

// Package monoid implements an algebraic structure with a single associative
// binary operation and an identity element.
// See the post about monoid in Go:
// https://github.com/fogfish/golem/blob/master/doc/monoid.md
package monoid

//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

//go:generate golem -lib -T int -generic github.com/fogfish/golem/seq/seq.go
//go:generate golem -lib -T int -generic github.com/fogfish/golem/seq/seq_monoid.go
//go:generate golem -lib -T string -generic github.com/fogfish/golem/seq/seq.go
//go:generate golem -lib -T string -generic github.com/fogfish/golem/seq/seq_monoid.go

// Package seq implements a base trait for sequences.
//
// Sequence is a special case for slice. Unlike built-in slices,
// Sequence support convenient methods to operate with slices and
// modifying them in immutable manner.
//
// Sequence do not implement methods to append element or
// concatenate another sequence, a built-in Golang
// `append`, `len` and `cap` shall be used.
package seq

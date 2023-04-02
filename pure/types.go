//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pure

// ContraMap turn morphisms around f: B ⟼ A
type ContraMap[A, B any] func(B) A

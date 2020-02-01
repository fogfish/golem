//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package generic

// T is generic type parameter identifier.
// The parameter is replaced with reference to the specific type.
//
//    func f() generic.T {/* ... */}
type T interface{}

// L is generic labelled type parameter identifier.
// The parameter is replaced with reference to the specific type.
//
//    func f() generic.L {/* ... */}
type L map[string]interface{}

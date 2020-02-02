//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

// Package cipher implements crypto algorithms for crypto generic.
package cipher

// Default is global variable with reference to default cipher.
var Default = NewKMS()

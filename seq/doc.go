//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

/*

Package seq implements a base traits for sequential data types.
Sequences generally behave very much like lists. The package support
convenient methods to operate with sequence of elements and modifying
them in immutable manner. The definition of sequence type is inspired by
Haskell's Data.Sequence and Scala's Seq trait.

TODO: Explain following methods via FMap
 - scan (partial folds)
 - fold
 - map

*/
package seq

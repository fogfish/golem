//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package seq

/*

TODO:

* Construction
-[ ] Empty ~~> New
-[x] Cons <| :+
-[x] Concatenate, Join, ><
-[-] FromSlice, slice is not efficient due to interface cast
-[x] FromFunction, Build, Gen, <-
-[x] Unfold, seed + function

* View (Iterator)
-[x] Head
-[x] Tail

* Query
-[x] Length - The number of elements in the sequence.
-[x] IsEmpty - Is this the empty sequence.

* Sequential searches (Filters)
-[x] TakeWhile
-[x] DropWhile
-[x] Span
-[x] Partition
-[x] GroupBy
-[x] Filter
-[x] FilterNot

* Sorting
-[x] Sort
-[x] SortBy
-[x] Distinct

* Indexing / Lookups
- [x] Find
- [x] Take
- [x] Drop
- [x] Span

* Indexing with Predicate
- Exists (at least one element holds)
- ForAll (all elements)
- Find

* Zip and UnZip
- Zip
- ZipWith
- UnZip
- UnZipWith

*/

// type Seq interface {
// }

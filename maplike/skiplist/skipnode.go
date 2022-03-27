//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package skiplist

import "fmt"

/*

Each element is represented by a tSkipNode in a skip list. Each node has
a height or level (length of fingers array), which corresponds to the number
of forward pointers the node has. When a new element is inserted into the list,
a node with a random level is inserted to represent the element. Random levels
are generated with a simple pattern: 50% are level 1, 25% are level 2, 12.5% are
level 3 and so on.
*/
type tSkipNode[K, V any] struct {
	key     K
	val     V
	fingers []*tSkipNode[K, V]
}

func (node *tSkipNode[K, V]) String() string {
	fingers := ""
	for _, x := range node.fingers {
		if x != nil {
			fingers = fingers + fmt.Sprintf("%v ", x.key)
		} else {
			fingers = fingers + fmt.Sprintf("nil ")
		}
	}

	return fmt.Sprintf("{%v\t| %s}", node.key, fingers)
}

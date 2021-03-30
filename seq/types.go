package seq

import "github.com/fogfish/golem"

/*

List is a linked list
*/
type List struct {
	head golem.Ord
	tail *List
}

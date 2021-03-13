/*

Package skiplist implements ...
http://citeseerx.ist.psu.edu/viewdoc/summary?doi=10.1.1.17.524

*/
package skiplist

import (
	"math/rand"

	"github.com/fogfish/golem"
)

/*

SkipList ...
*/
type SkipList interface {
	golem.MapLike
	Config(int, float64)
	Length() int
}

/*

tSkipNode ...
*/
type tSkipNode struct {
	key  golem.Ord
	val  golem.Data
	next []*tSkipNode
}

/*

tSkipList ...
*/
type tSkipList struct {
	head   *tSkipNode
	length int

	//
	//
	random rand.Source
	levels int
	p      []float64
}

var (
	_ SkipList = &tSkipList{}
	_ SkipList = (*tSkipList)(nil)
)

/*

tVoid is no element, holder of skip-list pointers
*/
type tVoid int

func (tVoid) Eq(golem.Eq) bool  { return false }
func (tVoid) Ne(golem.Eq) bool  { return false }
func (tVoid) Lt(golem.Ord) bool { return false }

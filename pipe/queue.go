//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pipe

import "sync"

type queue[A any] struct {
	head *q[A]
	tail *q[A]
	pool sync.Pool
}

type q[A any] struct {
	value *A
	next  *q[A]
}

func newq[A any]() *queue[A] {
	queue := &queue[A]{}
	queue.pool.New = func() interface{} { return &q[A]{} }
	return queue
}

func enq[A any](x *A, queue *queue[A]) {
	val := queue.pool.Get().(*q[A])
	val.value = x
	val.next = nil

	if queue.tail != nil {
		queue.tail.next = val
	}
	queue.tail = val

	if queue.head == nil {
		queue.head = val
	}
}

func deq[A any](queue *queue[A]) *A {
	val := queue.head
	queue.head = val.next

	if val == queue.tail {
		queue.tail = nil
	}

	queue.pool.Put(val)
	return val.value
}

func head[A any](queue *queue[A]) A {
	if queue.head == nil {
		return *new(A)
	}
	return *queue.head.value
}

func emit[A any](ch chan<- A, queue *queue[A]) chan<- A {
	if queue.head == nil {
		return nil
	}
	return ch
}

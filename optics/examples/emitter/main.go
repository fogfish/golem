//
// Copyright (C) 2022 - 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/fogfish/golem/optics"
	"github.com/fogfish/golem/pure"
	"github.com/fogfish/guid/v2"
)

//
// The most advanced usages of optics library to abstract structure.
//

// Let's assume the envelop to supports development of event-driven solutions.
// Primary purpose to build type safe approach on handling these events without
// too much definition of boilerplate code. Key element of this example is
// automatic handling of envelops
type Event[T any] struct {
	//
	// Unique identity for event
	// It is automatically defined by the library upon the transmission
	ID guid.K `json:"@id,omitempty"`

	//
	// Canonical IRI that defines a type of action.
	// It is automatically defined by the library upon the transmission
	Type string `json:"@type,omitempty"`

	//
	// ISO8601 timestamps when action has been created
	// It is automatically defined by the library upon the transmission
	Created time.Time `json:"created,omitempty"`

	//
	// The object upon which the event is carried out.
	Object T `json:"object,omitempty"`
}

// Our computation is abstraction F[Event[_]], HKT is required for type safeness
// See doc/higher-kinded-polymorphism.md
type EventType any

type EventKind[A any] pure.HKT[EventType, A]

func (Event[T]) HKT1(EventType) {}
func (Event[T]) HKT2(T)         {}

//
// Generic
//

// Writer is generic computation that deals with group of Event[T].
type Writer[T any, E EventKind[T]] struct {
	typeOf string
	shape  optics.Lens3[E, guid.K, string, time.Time]
}

func NewWriter[T any, E EventKind[T]]() Writer[T, E] {
	return Writer[T, E]{
		typeOf: reflect.TypeOf(new(E)).Elem().Name(),
		shape:  optics.ForShape3[E, guid.K, string, time.Time](),
	}
}

func (w Writer[T, E]) Write(evt *E) {
	// Note: This is a key feature of optics library.
	//       Processing of struct fields in generic way.
	w.shape.Put(evt, guid.G(guid.Clock), w.typeOf, time.Now())

	b, _ := json.MarshalIndent(evt, "", "  ")
	fmt.Println(string(b))
}

//
// Application only defines types and instances of generic algorithms
//

type User struct {
	Name      string
	Followers int
}

type EventUser Event[User]

func (EventUser) HKT1(EventType) {}
func (EventUser) HKT2(User)      {}

var eventUser = NewWriter[User, EventUser]()

type City struct {
	Name       string
	Population int
}

type EventCity Event[City]

func (EventCity) HKT1(EventType) {}
func (EventCity) HKT2(City)      {}

var eventCity = NewWriter[City, EventCity]()

func main() {
	eventUser.Write(&EventUser{Object: User{Name: "user", Followers: 10}})
	eventCity.Write(&EventCity{Object: City{Name: "user", Population: 1000}})
}

//
// FAQ:
//
// Q: What prevents to define interface { Put(...) } as abstraction on Event[T] and use embedding?
//
//   func (evt *Event[T]) Put(...)
//
//   type EventUser { Event[User] }
//
//   func (w Writer) Write(evt interface{ Put(...) })
//
// A: As discussed by doc/Abstract over Golang structure fields using optics, this is most of idiomatic
//    approach. It only suffers from type safety, complier would not be able to distinguish if writer is
//    fed by correct type. Any type that implement an interface accepted by writer.
//
//    busUser.Write(&EventUser{ ... })
//    busCity.Write(&EventCity{ ... }) <- not able to find error
//
//    Also, amount of boilerplate code would not be small.
//
//
// Q: Why writer elevate envelop to api? Usage of Write(obj *T) with assembly of envelop within
//    the writer function reduces needs of abstracting access to envelop field. It reduces needs of
//    optics abstraction.
//
// A: True, for this this simple example. The complex domain (e.g. Action as it is defined by
//    schema.org) require an envelop to be shared through generic middleware and application.
//    In this case, optics abstraction is a solution.
//

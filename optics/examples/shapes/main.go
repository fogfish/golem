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
	"time"

	"github.com/fogfish/golem/optics"
)

// Declare type and its shapes (getters & setter)
type User struct {
	Name      string
	Followers int
	Updated   time.Time
}

var user = optics.ForShape3[User, string, int, time.Time]()

type City struct {
	Name       string
	Population int
	Updated    time.Time
}

var city = optics.ForShape3[City, string, int, time.Time]()

// Generic algorithm that modifies struct fields
func show[T any](shape optics.Lens3[T, string, int, time.Time], v *T) {
	if s, i, t := shape.Get(v); t.IsZero() {
		shape.Put(v, s, i, time.Now())
	}

	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	show(user, &User{Name: "user", Followers: 10})
	show(city, &City{Name: "city", Population: 1000})
}

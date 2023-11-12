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

// Declare type and its lenses (getters & setter)
type User struct {
	Name    string
	Updated time.Time
}

var userT = optics.ForProduct1[User, time.Time]()

type City struct {
	Name    string
	Updated time.Time
}

var cityT = optics.ForProduct1[City, time.Time]()

// Generic algorithm that modifies struct fields
func show[T any](updated optics.Lens[T, time.Time], v *T) {
	if t := updated.Get(v); t.IsZero() {
		updated.Put(v, time.Now())
	}

	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}

func main() {
	show(userT, &User{Name: "user", Updated: time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)})
	show(cityT, &City{Name: "city"})
}

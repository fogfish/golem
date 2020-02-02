//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

//go:generate golem -T Identity -generic github.com/fogfish/golem/crypto/crypto.go

// Package identity is an example of custom ADT
package identity

// Identity is example data type
type Identity struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PinCode  int    `json:"pincode"`
}

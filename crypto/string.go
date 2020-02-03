//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package crypto

import (
	"encoding/json"
	"github.com/fogfish/golem/crypto/cipher"
)

// String is an alias built-in `string`. It shall be used as a container for sensitive
// data. Its sensitive value is not assignable to variable of type `string`. You have
// to either use helper method `PlainText` or cast it to string. This makes a simple
// protection against accidental leakage.
type String string

// UnmarshalJSON implements automatic decryption of data
func (value *String) UnmarshalJSON(b []byte) (err error) {
	var cryptotext string
	if err = json.Unmarshal(b, &cryptotext); err != nil {
		return
	}

	text, err := cipher.Default.Decrypt(cryptotext)
	if err != nil {
		return
	}

	*value = String(text)
	return
}

// MarshalJSON implements automatic encryption of sensitive strings during data marshalling.
func (value String) MarshalJSON() (bytes []byte, err error) {
	text, err := cipher.Default.Encrypt([]byte(value))
	if err != nil {
		return
	}

	return json.Marshal(text)
}

// PlainText returns plain text value
func (value String) PlainText() string {
	return string(value)
}

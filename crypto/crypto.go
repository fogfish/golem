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
	"github.com/fogfish/golem/generic"
)

// AnyT is generic definition of sensitive structures
type AnyT generic.L

// UnmarshalJSON implements automatic decryption of data
func (value *AnyT) UnmarshalJSON(b []byte) (err error) {
	type Referable AnyT

	var cryptotext string
	if err = json.Unmarshal(b, &cryptotext); err != nil {
		return
	}

	text, err := cipher.Default.Decrypt(cryptotext)
	if err != nil {
		return
	}

	var gen Referable
	if err = json.Unmarshal(text, &gen); err != nil {
		return
	}

	*value = AnyT(gen)
	return
}

// MarshalJSON implements automatic encryption of sensitive strings during data marshalling.
func (value AnyT) MarshalJSON() (bytes []byte, err error) {
	type Referable AnyT

	binary, err := json.Marshal(Referable(value))
	if err != nil {
		return
	}

	text, err := cipher.Default.Encrypt(binary)
	if err != nil {
		return
	}

	return json.Marshal(text)
}

// PlainText returns plain text value
func (value AnyT) PlainText() generic.L {
	return generic.L(value)
}

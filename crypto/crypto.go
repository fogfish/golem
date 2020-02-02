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

//
type AnyT generic.L

//
func (value *AnyT) UnmarshalJSON(b []byte) (err error) {
	type Referable AnyT

	var cryptotext string
	if err = json.Unmarshal(b, &cryptotext); err != nil {
		return err
	}
	text, err := cipher.Default.Decrypt(cryptotext)

	var gen Referable
	if err = json.Unmarshal(text, &gen); err != nil {
		return err
	}
	*value = AnyT(gen)

	return
}

//
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

//
func (value AnyT) PlainText() generic.L {
	return generic.L(value)
}

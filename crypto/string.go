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

//
type String string

//
func (value *String) UnmarshalJSON(b []byte) (err error) {
	var cryptotext string
	if err = json.Unmarshal(b, &cryptotext); err != nil {
		return err
	}
	text, err := cipher.Default.Decrypt(cryptotext)
	*value = String(text)

	return
}

//
func (value String) MarshalJSON() (bytes []byte, err error) {
	text, err := cipher.Default.Encrypt([]byte(value))
	if err != nil {
		return
	}

	return json.Marshal(text)
}

//
func (value String) PlainText() string {
	return string(value)
}

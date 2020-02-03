//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fogfish/golem/crypto"
	"github.com/fogfish/golem/crypto/cipher"
	"github.com/fogfish/golem/crypto/example/identity"
)

// User example usage of data type
type User struct {
	Identity identity.Crypto `json:"identity"`
	About    crypto.String   `json:"about"`
}

func main() {
	cipher.Default.UseKey(os.Args[1])

	user := User{
		Identity: identity.Crypto{
			Email:    "any@example.com",
			Password: "sensitive data",
			PinCode:  1234,
		},
		About: "Lorem ipsum dolor sit amet",
	}

	log.Println("==> encrypting...")
	bytes, err := json.Marshal(&user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bytes))

	log.Println("==> decrypting...")
	err = json.Unmarshal(bytes, &user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}

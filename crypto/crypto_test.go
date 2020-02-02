//
// Copyright (C) 2019 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package crypto_test

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/fogfish/golem/crypto"
	"github.com/fogfish/golem/crypto/cipher"
	"github.com/fogfish/golem/generic"
	"github.com/fogfish/it"
)

type MyString struct {
	Secret crypto.String `json:"secret"`
}

type MyJSON struct {
	Secret crypto.AnyT `json:"secret"`
}

func TestStringUnmarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyString{}
	input := []byte("{\"secret\":\"cGxhaW50ZXh0\"}")

	it.Ok(t).
		If(json.Unmarshal(input, &value)).Should().Equal(nil).
		If(value.Secret).Should().Equal(crypto.String("plaintext")).
		If(value.Secret.Value()).Should().Equal("plaintext")
}

func TestStringMarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyString{crypto.String("plaintext")}
	bytes, err := json.Marshal(value)

	it.Ok(t).
		If(err).Should().Equal(nil).
		If(bytes).Should().Equal([]byte("{\"secret\":\"cGxhaW50ZXh0\"}"))
}

func TestAnyTUnmarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyJSON{}
	input := "{\"secret\":\"eyJ0ZXh0IjoicGxhaW50ZXh0In0=\"}"

	it.Ok(t).
		If(json.Unmarshal([]byte(input), &value)).Should().Equal(nil).
		If(value.Secret).Should().Equal(crypto.AnyT{"text": "plaintext"}).
		If(value.Secret.Value()).Should().Equal(generic.L{"text": "plaintext"})
}

func TestCryptoMarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyJSON{crypto.AnyT{"text": "plaintext"}}
	bytes, err := json.Marshal(value)

	it.Ok(t).
		If(err).Should().Equal(nil).
		If(bytes).Should().Equal([]byte("{\"secret\":\"eyJ0ZXh0IjoicGxhaW50ZXh0In0=\"}"))
}

//
//
type mock struct {
	kmsiface.KMSAPI
}

func (mock) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return &kms.DecryptOutput{
		Plaintext: input.CiphertextBlob,
	}, nil
}

func (mock) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return &kms.EncryptOutput{
		CiphertextBlob: input.Plaintext,
	}, nil
}

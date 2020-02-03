//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package crypto_test

import (
	"encoding/json"
	"errors"
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
		If(value.Secret.PlainText()).Should().Equal("plaintext")
}

func TestStringUnmarshalFail(t *testing.T) {
	cipher.Default.Mock(fail{})

	value := MyString{}
	input := []byte("{\"secret\":\"cGxhaW50ZXh0\"}")

	it.Ok(t).
		If(
			func() error {
				return json.Unmarshal(input, &value)
			},
		).Should().Intercept(ErrDecrypt)
}

func TestStringMarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyString{crypto.String("plaintext")}
	bytes, err := json.Marshal(value)

	it.Ok(t).
		If(err).Should().Equal(nil).
		If(bytes).Should().Equal([]byte("{\"secret\":\"cGxhaW50ZXh0\"}"))
}

func TestStringMarshalFail(t *testing.T) {
	cipher.Default.Mock(fail{})

	value := MyString{crypto.String("plaintext")}

	it.Ok(t).
		If(
			func() error {
				_, err := json.Marshal(value)
				return err
			},
		).Should().Intercept(ErrEncrypt)
}

func TestAnyTUnmarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyJSON{}
	input := []byte("{\"secret\":\"eyJ0ZXh0IjoicGxhaW50ZXh0In0=\"}")

	it.Ok(t).
		If(json.Unmarshal(input, &value)).Should().Equal(nil).
		If(value.Secret).Should().Equal(crypto.AnyT{"text": "plaintext"}).
		If(value.Secret.PlainText()).Should().Equal(generic.L{"text": "plaintext"})
}

func TestAnyTUnmarshalFail(t *testing.T) {
	cipher.Default.Mock(fail{})

	value := MyJSON{}
	input := []byte("{\"secret\":\"eyJ0ZXh0IjoicGxhaW50ZXh0In0=\"}")

	it.Ok(t).
		If(
			func() error {
				return json.Unmarshal(input, &value)
			},
		).Should().Intercept(ErrDecrypt)
}

func TestAnyTMarshalJSON(t *testing.T) {
	cipher.Default.Mock(mock{})

	value := MyJSON{crypto.AnyT{"text": "plaintext"}}
	bytes, err := json.Marshal(value)

	it.Ok(t).
		If(err).Should().Equal(nil).
		If(bytes).Should().Equal([]byte("{\"secret\":\"eyJ0ZXh0IjoicGxhaW50ZXh0In0=\"}"))
}

func TestAnyTMarshalFail(t *testing.T) {
	cipher.Default.Mock(fail{})

	value := MyJSON{crypto.AnyT{"text": "plaintext"}}

	it.Ok(t).
		If(
			func() error {
				_, err := json.Marshal(value)
				return err
			},
		).Should().Intercept(ErrEncrypt)
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

//
//
type fail struct {
	kmsiface.KMSAPI
}

var ErrDecrypt = errors.New("Unable to decrypt")
var ErrEncrypt = errors.New("Unable to encrypt")

func (fail) Decrypt(input *kms.DecryptInput) (*kms.DecryptOutput, error) {
	return nil, ErrDecrypt
}

func (fail) Encrypt(input *kms.EncryptInput) (*kms.EncryptOutput, error) {
	return nil, ErrEncrypt
}

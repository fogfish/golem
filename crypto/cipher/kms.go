//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package cipher

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
)

//
type KMS struct {
	api kmsiface.KMSAPI
	key string
}

//
func NewKMS() *KMS {
	return &KMS{
		kms.New(session.Must(session.NewSession())),
		"",
	}
}

//
func (c *KMS) Mock(api kmsiface.KMSAPI) {
	c.api = api
}

//
func (c *KMS) UseKey(key string) {
	c.key = key
}

//
func (c *KMS) Decrypt(cryptotext string) (plaintext []byte, err error) {
	bytes, err := base64.StdEncoding.DecodeString(cryptotext)
	if err != nil {
		return
	}

	input := &kms.DecryptInput{
		CiphertextBlob: []byte(bytes),
	}

	result, err := c.api.Decrypt(input)
	if err != nil {
		return
	}

	return result.Plaintext, nil
}

//
func (c *KMS) Encrypt(plaintext []byte) (cryptotext string, err error) {
	input := &kms.EncryptInput{
		KeyId:     aws.String(c.key),
		Plaintext: plaintext,
	}

	result, err := c.api.Encrypt(input)
	if err != nil {
		return
	}

	return base64.StdEncoding.EncodeToString(result.CiphertextBlob), nil
}

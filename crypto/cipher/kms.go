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

// KMS is the context of AWS KMS cipher
type KMS struct {
	api kmsiface.KMSAPI
	key string
}

// NewKMS returns AWS KMS context
func NewKMS() *KMS {
	return &KMS{
		kms.New(session.Must(session.NewSession())),
		"",
	}
}

// Mock replaces instances of AWS KMS API
func (c *KMS) Mock(api kmsiface.KMSAPI) {
	c.api = api
}

// UseKey defines encryption key
func (c *KMS) UseKey(key string) {
	c.key = key
}

// Decrypt uses AWS KMS API to decrypt cryptotext.
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

// Encrypt uses AWS KMS API to encrypt plaintext.
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

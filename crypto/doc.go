//
// Copyright (C) 2020 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

// Package crypto implements semi-auto cipher codec of textual content and custom
// Algebraic Data Types. Encryption/Decryption process is transparent for developers.
// It is embedded into `json.Marshal` and `json.Unmarshal` routines so that protection
// of sensitive data happens during the process of data serialization.
//
// The package offers the developer friendly solution to apply Application/Record Level
// Encryption with help of AWS KMS for sensitive structured data (e.g. JSON).
// The design aims to address few requirements:
//
// ↣ transparent for developers - encryption/decryption is built with semi-auto codec.
// It makes a "magic" of switching representation between crypto/plain texts.
// Developer just declares the intent to protect sensitive data.
//
// ↣ compile time type-safeness - the sensitive data is modelled with algebraic data types.
// The type tagging (annotation) is used to declare the the intent to protect sensitive
// data. Golang compiler discover and prevents type errors or other glitches at the time
// it assembles binaries.
//
// ↣  generic - encryption/decryption are generic algorithms applicable to any algebraic
// data types (not only to strings). The library provides ability to apply algorithms for
// any product type in developer's application context.
//
// ↣ data in use is not supported. Developers have to combine this library with other
//solutions.
//
//
// The package implements final type to encrypt/decrypt strings `crypto.String` and
// generic `crypto.AnyT` type, which allows to handle any application specific ADTs.
package crypto

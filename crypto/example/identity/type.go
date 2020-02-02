//go:generate golem -T Identity -generic github.com/fogfish/golem/crypto/crypto.go

// Package identity is an example of custom ADT
package identity

// Identity is example data type
type Identity struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	PinCode  int    `json:"pincode"`
}

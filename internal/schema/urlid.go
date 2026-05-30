package schema

import (
	"crypto/rand"
)

const urlIDCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const urlIDLength = 10

// A short randomized ID for use with document URLs.
type URLID string

// Generates a new URLID.
func NewURLID() URLID {
	var out [urlIDLength]byte
	rand.Read(out[:])

	for i, n := range out {
		// Replace the random number with a byte from the charset. Modulo is necessary to prevent out-of-bounds.
		out[i] = urlIDCharset[n%urlIDLength]
	}

	return URLID(out[:])
}

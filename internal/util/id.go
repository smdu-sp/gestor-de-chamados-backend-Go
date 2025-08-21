package util

import (
	"crypto/rand"
	"encoding/hex"
)

// Gera um ID único de 16 bytes e retorna como string hexadecimal 
func NewID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

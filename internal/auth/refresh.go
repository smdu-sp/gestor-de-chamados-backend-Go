package auth

import (
	"crypto/rand"
	"encoding/hex"
)

type RefreshManager struct{}

func NewRefreshManager() *RefreshManager {
	return &RefreshManager{}
}

func (m *RefreshManager) Generate() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	
	return hex.EncodeToString(b)
}

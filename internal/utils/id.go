package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"time"
)

var ErrUUIDv7Generation = errors.New("falha ao gerar bytes aleatórios para UUID v7")

// UUIDv7 representa um UUID v7 de 16 bytes
type UUIDv7 [16]byte

// NewUUIDv7Bytes gera um UUID v7 retornando os 16 bytes
func NewUUIDv7Bytes() (UUIDv7, error) {
	var uuid UUIDv7

	// Timestamp em milissegundos (48 bits)
	now := uint64(time.Now().UnixMilli())
	uuid[0] = byte(now >> 40)
	uuid[1] = byte(now >> 32)
	uuid[2] = byte(now >> 24)
	uuid[3] = byte(now >> 16)
	uuid[4] = byte(now >> 8)
	uuid[5] = byte(now)

	// Preenche os 10 bytes restantes com aleatório
	if _, err := rand.Read(uuid[6:]); err != nil {
		return uuid, NewAppError(
			"[utils.NewUUIDv7Bytes]",
			LevelError,
			"erro ao gerar UUID v7",
			fmt.Errorf("%w: %v", ErrUUIDv7Generation, err),
		)
	}

	// Define versão v7 (4 bits)
	uuid[6] = (uuid[6] & 0x0f) | 0x70

	// Define variante RFC 4122 (2 bits)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

// NewUUIDv7String gera um UUID v7 diretamente como string
func NewUUIDv7String() (string, error) {
	uuid, err := NewUUIDv7Bytes()
	if err != nil {
		return "", fmt.Errorf("[utils.NewUUIDv7String]: %w", err)
	}
	return uuid.String(), nil
}

// String converte o UUIDv7 para o formato padrão com hífens
func (u UUIDv7) String() string {
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		uint32(u[0])<<24|uint32(u[1])<<16|uint32(u[2])<<8|uint32(u[3]),
		uint16(u[4])<<8|uint16(u[5]),
		uint16(u[6])<<8|uint16(u[7]),
		uint16(u[8])<<8|uint16(u[9]),
		u[10:16],
	)
}

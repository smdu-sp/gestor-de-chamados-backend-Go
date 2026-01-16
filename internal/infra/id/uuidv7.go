package id

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"sync"
	"time"
)

// UUIDv7 representa um UUID versão 7 conforme RFC 9562.
// Estrutura: 48 bits timestamp + 4 bits version + 12 bits rand_a + 2 bits variant + 62 bits rand_b
type UUIDv7 [16]byte

// Estado mantém o estado necessário para geração monotônica de UUIDs.
type estado struct {
	mu            sync.Mutex
	lastTimestamp uint64 // Unix timestamp em milissegundos
	lastRandA     uint16 // 12 bits menos significativos (0x0000-0x0FFF)
}

// estadoGlobal é o estado compartilhado para geração de UUIDs.
var estadoGlobal = &estado{}

const (
	maskRandA      = 0x0FFF // Máscara para os 12 bits de rand_a
	version        = 0x7    // Versão do UUID (v7 = 0111 = 0x7)
	versionShift   = 4      // Posição do nibble de versão no byte 6
	variantRFC4122 = 0x80   // Variante RFC 4122 (10xxxxxx = 0x80)
	variantMask    = 0x3F   // Máscara para limpar os 2 bits de variante
)

// GerarUUIDv7Bytes gera um UUID v7 monotônico em conformidade com a RFC 9562.
//
// Características:
//   - Ordenação lexicográfica garantida (monotonic)
//   - Proteção contra clock rollback
//   - Incremento de rand_a para UUIDs no mesmo milissegundo
//   - Espera ativa em caso de overflow de rand_a
//
// Retorna erro apenas em caso de falha na geração de bytes aleatórios.
func GerarUUIDv7Bytes() (UUIDv7, error) {
	// 1- Preparar variável para o UUID
	var uuid UUIDv7

	// 2- Bloquear estado global para geração segura
	estadoGlobal.mu.Lock()
	defer estadoGlobal.mu.Unlock()

	// 3- Obter timestamp e rand_a com lógica monotônica
	timestamp, randA, err := estadoGlobal.obterTimestampERandA()
	if err != nil {
		return uuid, err
	}

	// 4- Preencher os campos do UUID
	estadoGlobal.escreverTimestamp(&uuid, timestamp)
	estadoGlobal.escreverVersaoERandA(&uuid, randA)

	// 5- Preencher rand_b aleatório
	if err := estadoGlobal.escreverRandB(&uuid); err != nil {
		return uuid, err
	}

	// 6- Definir os bits de variante
	estadoGlobal.escreverVariante(&uuid)

	return uuid, nil
}

// obterTimestampERandA gerencia a lógica de geração monotônica.
func (e *estado) obterTimestampERandA() (timestamp uint64, randA uint16, err error) {
	// 1- Obter timestamp atual em milissegundos
	timestamp = uint64(time.Now().UnixMilli())

	// 2- Garantir monotonicidade
	timestamp = max(timestamp, e.lastTimestamp)

	// 3- Lógica para rand_a
	switch {
	// caso o timestamp seja maior que o último registrado
	case timestamp > e.lastTimestamp:
		randA, err = e.gerarNovoRandA()
		if err != nil {
			return 0, 0, err
		}
		e.lastTimestamp = timestamp
		e.lastRandA = randA

	// caso o timestamp seja igual ao último registrado
	case e.lastRandA < maskRandA:
		e.lastRandA++
		randA = e.lastRandA

	// caso o rand_a atingiu o máximo
	default:
		timestamp, randA, err = e.aguardarProximoMillisegundo()
		if err != nil {
			return 0, 0, err
		}
		e.lastTimestamp = timestamp
		e.lastRandA = randA
	}

	return timestamp, randA, nil
}

// gerarNovoRandA gera um valor aleatório de 12 bits.
func (e *estado) gerarNovoRandA() (uint16, error) {
	// 1- Gerar 2 bytes aleatórios
	var buf [2]byte

	// 2- Aplicar máscara para obter 12 bits
	if _, err := rand.Read(buf[:]); err != nil {
		return 0, fmt.Errorf("falha ao gerar rand_a: %w", err)
	}

	// 3- Retornar valor mascarado
	return binary.BigEndian.Uint16(buf[:]) & maskRandA, nil
}

// aguardarProximoMillisegundo espera ativamente até o próximo milissegundo.
func (e *estado) aguardarProximoMillisegundo() (uint64, uint16, error) {
	for {
		timestamp := uint64(time.Now().UnixMilli())
		if timestamp > e.lastTimestamp {
			randA, err := e.gerarNovoRandA()
			if err != nil {
				return 0, 0, err
			}
			return timestamp, randA, nil
		}
		// Sleep mínimo para evitar busy-wait excessivo
		time.Sleep(100 * time.Microsecond)
	}
}

// escreverTimestamp escreve os 48 bits de timestamp no UUID.
func (e *estado) escreverTimestamp(uuid *UUIDv7, timestamp uint64) {
	uuid[0] = byte(timestamp >> 40)
	uuid[1] = byte(timestamp >> 32)
	uuid[2] = byte(timestamp >> 24)
	uuid[3] = byte(timestamp >> 16)
	uuid[4] = byte(timestamp >> 8)
	uuid[5] = byte(timestamp)
}

// escreverVersaoERandA escreve os 4 bits de versão e 12 bits de rand_a.
func (e *estado) escreverVersaoERandA(uuid *UUIDv7, randA uint16) {
	// Byte 6: versão (4 bits) + rand_a high (4 bits)
	uuid[6] = (version << versionShift) | byte(randA>>8)
	
	// Byte 7: rand_a low (8 bits)
	uuid[7] = byte(randA)
}

// escreverRandB escreve os 64 bits aleatórios rand_b.
func (e *estado) escreverRandB(uuid *UUIDv7) error {
	if _, err := rand.Read(uuid[8:16]); err != nil {
		return fmt.Errorf("falha ao gerar rand_b: %w", err)
	}
	return nil
}

// escreverVariante escreve os 2 bits de variante RFC 4122.
func (e *estado) escreverVariante(uuid *UUIDv7) {
	// Byte 8: variante (2 bits) + rand_b (6 bits)
	uuid[8] = (uuid[8] & variantMask) | variantRFC4122
}

// GerarUUIDv7String gera um UUID v7 no formato string canônico.
//
// Formato: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
// Exemplo: 018e8c6a-73f8-7001-8b3a-d5f3469a2e0f
func GerarUUIDv7String() (string, error) {
	uuid, err := GerarUUIDv7Bytes()
	if err != nil {
		return "", err
	}
	return uuid.String(), nil
}

// String converte o UUIDv7 para o formato canônico 8-4-4-4-12.
// Otimizado para performance com buffer pré-alocado e lookup table.
func (u UUIDv7) String() string {
	const hexDigits = "0123456789abcdef"
	var buf [36]byte

	// Posições dos hífens: 8, 13, 18, 23
	encodeHex(buf[0:8], u[0:4], hexDigits)
	buf[8] = '-'
	encodeHex(buf[9:13], u[4:6], hexDigits)
	buf[13] = '-'
	encodeHex(buf[14:18], u[6:8], hexDigits)
	buf[18] = '-'
	encodeHex(buf[19:23], u[8:10], hexDigits)
	buf[23] = '-'
	encodeHex(buf[24:36], u[10:16], hexDigits)

	return string(buf[:])
}

// encodeHex converte bytes para representação hexadecimal no buffer.
func encodeHex(dst []byte, src []byte, hexDigits string) {
	j := 0
	for _, b := range src {
		dst[j] = hexDigits[b>>4]
		dst[j+1] = hexDigits[b&0x0F]
		j += 2
	}
}

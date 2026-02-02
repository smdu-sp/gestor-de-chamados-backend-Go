package usuario

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// Regex para validação de email.
const emailRegexPattern  = `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`

// Compila o regex apenas uma vez para melhorar a performance.
var emailRegex = regexp.MustCompile(emailRegexPattern)


// Email representa o objeto de valor Email.
type Email struct {
	endereco string
}

// NovoEmail cria uma nova instância de Email.
func NovoEmail(endereco string) Email {
	endereco = strings.ToLower(strings.TrimSpace(endereco))
	return Email{endereco: endereco}
}

// ValidarEmail valida o formato do email
func (e *Email) ValidarEmail() error {
	e.endereco = strings.ToLower(strings.TrimSpace(e.endereco))

	if e.endereco == "" {
		return errors.New("o email não pode estar vazio")
	}
	if !emailRegex.MatchString(e.endereco) {
		return errors.New("o email é inválido, formato esperado: user@example.com")
	}
	return nil
}

// String retorna o valor do email como string.
func (e Email) String() string {
	return e.endereco
}

// MarshalJSON implementa a interface json.Marshaler.
func (e Email) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.endereco)
}

// UnmarshalJSON implementa a interface json.Unmarshaler.
func (e *Email) UnmarshalJSON(data []byte) error {
    if err := json.Unmarshal(data, &e.endereco); err != nil {
        return fmt.Errorf("erro ao desserializar JSON: %w", err)
    }
    e.endereco = strings.ToLower(strings.TrimSpace(e.endereco))
    return nil
}

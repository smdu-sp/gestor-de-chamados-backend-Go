package config

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ErrosConfig é um erro estruturado especializado para validação
// de variáveis de ambiente (.env / config).
//
// Ele pode ser logado como string (via Error())
// ou como JSON estruturado (via MarshalJSON()).
type ErrosConfig struct {
	erros map[string]string
}

// NovoErrosConfig cria uma instância vazia.
func NovoErrosConfig() *ErrosConfig {
	return &ErrosConfig{
		erros: make(map[string]string),
	}
}

// Add registra um erro relacionado a um campo de configuração.
func (c *ErrosConfig) Add(campo, mensagem string) {
	c.erros[campo] = mensagem
}

// HaErros informa se algum erro foi registrado.
func (c *ErrosConfig) HaErros() bool {
	return len(c.erros) > 0
}

// Error retorna um texto simples e amigável para logs humanos.
// Esse texto aparece quando você usa: logger.Error("msg", "error", err)
func (c ErrosConfig) Error() string {
	if len(c.erros) == 0 {
		return ""
	}

	// Extrai e ordena as chaves
	chaves := make([]string, 0, len(c.erros))
	for chave := range c.erros {
		chaves = append(chaves, chave)
	}

	sort.Strings(chaves)

	// Monta a mensagem
	partes := make([]string, 0, len(chaves))
	for _, k := range chaves {
		partes = append(partes, fmt.Sprintf("%s: %s", k, c.erros[k]))
	}

	return strings.Join(partes, "; ")
}

// Erros retorna o mapa cru — útil para testes ou inspeções internas.
func (c ErrosConfig) Erros() map[string]string {
	return c.erros
}

// MarshalJSON fornece uma estrutura amigável para ELK.
// slog usará automaticamente este método se você passar o erro
// como valor (ex: "config.errors", err).
func (c ErrosConfig) MarshalJSON() ([]byte, error) {
	if len(c.erros) == 0 {
		return []byte("{}"), nil
	}

	// Constrói uma lista ordenada
	type variavelErro struct {
		Variavel string `json:"variavel"`
		Erro     string `json:"erro"`
	}

	// Extrai e ordena as chaves
	chaves := make([]string, 0, len(c.erros))
	for chave := range c.erros {
		chaves = append(chaves, chave)
	}
	sort.Strings(chaves)

	// Popula a lista ordenada
	ordenados := make([]variavelErro, 0, len(chaves))
	for _, k := range chaves {
		ordenados = append(ordenados, variavelErro{
			Variavel: k,
			Erro:     c.erros[k],
		})
	}

	// Serializa como array de objetos ordenado
	return json.Marshal(ordenados)
}

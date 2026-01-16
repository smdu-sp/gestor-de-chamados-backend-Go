package domain

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// ErrosValidacao é um erro estruturado especializado para validação de dados de domínio.
//
// Ele pode ser logado como string (via Error())
// ou como JSON estruturado (via MarshalJSON()).
type ErrosValidacao struct {
	erros map[string]string
}

// NovoErrosValidacao cria uma nova instância de ErrosValidacao.
func NovoErrosValidacao() *ErrosValidacao {
	return &ErrosValidacao{
		erros: make(map[string]string),
	}
}

// Add adiciona um erro a um campo específico.
func (v *ErrosValidacao) Add(campo, mensagem string) {
	v.erros[campo] = mensagem
}

// HaErros indica se há erros acumulados.
func (v *ErrosValidacao) HaErros() bool {
	return len(v.erros) > 0
}

// Error implementa a interface error.
func (v ErrosValidacao) Error() string {
	if len(v.erros) == 0 {
		return ""
	}

	// extrai as chaves
	chaves := make([]string, 0, len(v.erros))
	for valor := range v.erros {
		chaves = append(chaves, valor)
	}

	sort.Strings(chaves)

	// monta a mensagem
	partes := make([]string, 0, len(chaves))
	for _, k := range chaves {
		partes = append(partes, fmt.Sprintf("%s: %s", k, v.erros[k]))
	}

	return strings.Join(partes, "; ")
}

// Erros retorna o mapa cru, útil para casos onde se quer inspecionar os erros.
func (v ErrosValidacao) Erros() map[string]string {
	return v.erros
}

// MarshalJSON permite serializar o erro como {"Campo": "mensagem"}.
func (v ErrosValidacao) MarshalJSON() ([]byte, error) {
	if len(v.erros) == 0 {
		return []byte("{}"), nil
	}

	// Constrói uma lista ordenada
	type campoErro struct {
		Campo string `json:"campo"`
		Erro  string `json:"erro"`
	}

	// Extrai e ordena as chaves
	chaves := make([]string, 0, len(v.erros))
	for valor := range v.erros {
		chaves = append(chaves, valor)
	}
	sort.Strings(chaves)

	// Popula a lista ordenada
	ordenados := make([]campoErro, 0, len(chaves))
	for _, k := range chaves {
		ordenados = append(ordenados, campoErro{
			Campo: k,
			Erro:  v.erros[k],
		})
	}

	// Serializa como um array de objetos ordenado
	return json.Marshal(ordenados)
}

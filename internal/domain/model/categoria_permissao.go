package model

import (
	"errors"
	"fmt"
	"time"
)

// CategoriaPermissao representa a permissão de um usuário para uma categoria específica.
type CategoriaPermissao struct {
	CategoriaID  string    `json:"categoriaId"`
	UsuarioID    string    `json:"usuarioId"`
	Permissao    Permissao `json:"permissao"`
	CriadoEm     time.Time `json:"criadoEm"`
	AtualizadoEm time.Time `json:"atualizadoEm"`
}

var (
	ErrCategoriaPermissaoIDInvalido = errors.New("ID de categoriaPermissao inválido")
)

// NewCategoriaPermissao cria uma nova instância de CategoriaPermissao com os dados fornecidos.
func NewCategoriaPermissao(categoriaID, usuarioID string, permissao Permissao) (*CategoriaPermissao, error) {
	now := time.Now()
	categoriaPermissao := &CategoriaPermissao{
		CategoriaID:  categoriaID,
		UsuarioID:    usuarioID,
		Permissao:    permissao,
		CriadoEm:     now,
		AtualizadoEm: now,
	}
	if err := ValidarCategoriaPermissao(categoriaPermissao); err != nil {
		return nil, err
	}
	return categoriaPermissao, nil
}

// ValidarCategoriaPermissao valida os campos da permissão de categoria.
func ValidarCategoriaPermissao(c *CategoriaPermissao) error {
	var erros []error

	if c.CategoriaID == "" {
		erros = append(erros, ErrCategoriaPermissaoIDInvalido)
	}
	if c.UsuarioID == "" {
		erros = append(erros, ErrIDInvalido)
	}
	if err := ValidarPermissao(c.Permissao); err != nil {
		erros = append(erros, err)
	}
	if len(erros) > 0 {
		return fmt.Errorf("[model.ValidarCategoriaPermissao] erros de validação: %v", erros)
	}
	return nil
}

// CategoriaPermissaoFiltro representa os critérios de filtro para buscar permissões de categoria.
type CategoriaPermissaoFiltro struct {
	Pagina      int
	Limite      int
	CategoriaID *string
	UsuarioID   *string
	Permissao   *string
}

// String retorna uma representação em string da CategoriaPermissao para fins de logging.
func (c *CategoriaPermissao) String() string {
	return fmt.Sprintf(
		"[CategoriaID: %s | UsuarioID: %s | Permissao: %v]",
		c.CategoriaID, c.UsuarioID, c.Permissao,
	)
}
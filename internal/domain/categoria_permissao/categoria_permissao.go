package categoria_permissao

import (
	"fmt"
	"strings"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// CategoriaPermissao representa a permissão de um usuário para uma categoria específica.
type CategoriaPermissao struct {
	categoriaID  string
	usuarioID    string
	permissao    usr.Permissao
	criadoEm     time.Time
	atualizadoEm time.Time
}

// Novo cria uma nova instância de CategoriaPermissao com os dados fornecidos.
func Novo(categoriaID, usuarioID string, permissao usr.Permissao) (*CategoriaPermissao, error) {
	now := time.Now()
	categoriaPermissao := CategoriaPermissao{
		categoriaID:  strings.TrimSpace(categoriaID),
		usuarioID:    strings.TrimSpace(usuarioID),
		permissao:    permissao,
		criadoEm:     now,
		atualizadoEm: now,
	}

	if err := categoriaPermissao.Validar(); err != nil {
		return nil, err
	}

	return &categoriaPermissao, nil
}

// CarregarDoDB recebe uma estrutura CategoriaPermissaoDB e retorna uma instância de CategoriaPermissao.
func CarregarDoDB(c CategoriaPermissaoDB) *CategoriaPermissao {
	return &CategoriaPermissao{
		categoriaID:  c.CategoriaID,
		usuarioID:    c.UsuarioID,
		permissao:    usr.Permissao(c.Permissao),
		criadoEm:     c.CriadoEm,
		atualizadoEm: c.AtualizadoEm,
	}
}

// Validar valida os campos da permissão de categoria.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (c CategoriaPermissao) Validar() error {
	erros := domain.NovoErrosValidacao()

	if c.categoriaID == "" {
		erros.Add("CategoriaID", "o ID da categoria não pode estar vazio")
	}
	if c.usuarioID == "" {
		erros.Add("UsuarioID", "o ID do usuário não pode estar vazio")
	}
	if err := usr.ValidarPermissao(c.permissao); err != nil {
		erros.Add("Permissao", err.Error())
	}

	if erros.HaErros() {
		return erros
	}

	return nil
}

// ComDadosAtualizados recebe os parâmetros para atualizar a permissão de categoria 
// e retorna uma nova instância de CategoriaPermissao com os dados atualizados.
//
// Em caso de erros de validação, retorna uma instância de domain.ErrosValidacao.
func (c CategoriaPermissao) ComDadosAtualizados(params AtualizarParams) (CategoriaPermissao, error) {
	c.permissao = params.Permissao
	c.atualizadoEm = time.Now()

	if err := c.Validar(); err != nil {
		return CategoriaPermissao{}, err
	}

	return c, nil
}

// Metódos de acesso aos campos da CategoriaPermissao.

func (c CategoriaPermissao) CategoriaID() string     { return c.categoriaID }
func (c CategoriaPermissao) UsuarioID() string       { return c.usuarioID }
func (c CategoriaPermissao) Permissao() string       { return c.permissao.String() }
func (c CategoriaPermissao) CriadoEm() time.Time     { return c.criadoEm }
func (c CategoriaPermissao) AtualizadoEm() time.Time { return c.atualizadoEm }

// String retorna uma representação em string da permissão de categoria.
func (c CategoriaPermissao) String() string {
	return fmt.Sprintf(
		"[CategoriaID: %s, UsuarioID: %s, Permissao: %s, CriadoEm: %s, AtualizadoEm: %s]",
		c.categoriaID,
		c.usuarioID,
		c.permissao.String(),
		c.CriadoEm().Format(time.RFC3339),
		c.AtualizadoEm().Format(time.RFC3339),
	)
}

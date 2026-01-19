package subcategoria

import (
	"strings"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// TestNovoSubcategoria testa a criação de novas subcategorias com várias combinações de entradas válidas e inválidas.
func TestNovoSubcategoria(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		nome        string
		categoriaID string
		wantErr     bool
	}{
		{
			name:        "criar subcategoria válida",
			id:          "  sub-123  ",
			nome:        "  Hardware  ",
			categoriaID: "  cat-456  ",
			wantErr:     false,
		},
		{
			name:        "erro quando id é vazio",
			id:          "",
			nome:        "Hardware",
			categoriaID: "cat-1",
			wantErr:     true,
		},
		{
			name:        "erro quando nome é muito curto",
			id:          "sub-1",
			nome:        "AB",
			categoriaID: "cat-1",
			wantErr:     true,
		},
		{
			name:        "erro quando nome é muito longo",
			id:          "sub-1",
			nome:        strings.Repeat("a", tamanhoMaxNome+1),
			categoriaID: "cat-1",
			wantErr:     true,
		},
		{
			name:        "erro quando categoriaID é vazio",
			id:          "sub-1",
			nome:        "Hardware",
			categoriaID: "",
			wantErr:     true,
		},
		{
			name:        "erro com múltiplos campos inválidos",
			id:          "",
			nome:        "",
			categoriaID: "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			subc, err := Novo(tt.id, tt.nome, tt.categoriaID)

			if tt.wantErr {
				if err == nil {
					t.Fatal("esperava erro, mas recebeu nil")
				}

				// Garantia importante: erro de domínio estruturado
				if _, ok := err.(*dmn.ErrosValidacao); !ok {
					t.Fatalf("esperava ErrosValidacao, mas recebeu %T", err)
				}

				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			// ---- asserts do estado da entidade ----

			if subc.ID() != strings.TrimSpace(tt.id) {
				t.Errorf("ID não foi normalizado corretamente")
			}

			if subc.Nome() != strings.TrimSpace(tt.nome) {
				t.Errorf("Nome não foi normalizado corretamente")
			}

			if subc.CategoriaID() != strings.TrimSpace(tt.categoriaID) {
				t.Errorf("CategoriaID não foi normalizado corretamente")
			}

			if !subc.Status() {
				t.Errorf("status esperado true")
			}

			if subc.CriadoEm().IsZero() {
				t.Errorf("CriadoEm não deveria ser zero")
			}

			if subc.AtualizadoEm().IsZero() {
				t.Errorf("AtualizadoEm não deveria ser zero")
			}
		})
	}
}

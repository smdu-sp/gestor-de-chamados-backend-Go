package subcategoria

import (
	"strings"
	"testing"
	"time"

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

// TestCarregarDoDB testa a criação de uma subcategoria a partir dos dados do banco de dados.
func TestCarregarDoDB(t *testing.T) {
	now := time.Now()

	subcDB := SubcategoriaDB{
		ID:          "sub-123",
		Nome:        "Software",
		Status:      false,
		CategoriaID: "cat-456",
		CriadoEm:    now,
		AtualizadoEm: now,
	}
	
	subc := CarregarDoBD(subcDB)

	if subc.ID() != subcDB.ID {
		t.Errorf("esperava ID %s, recebeu %s", subcDB.ID, subc.ID())
	}

	if subc.Nome() != subcDB.Nome {
		t.Errorf("esperava Nome %s, recebeu %s", subcDB.Nome, subc.Nome())
	}

	if subc.Status() != subcDB.Status {
		t.Errorf("esperava Status %v, recebeu %v", subcDB.Status, subc.Status())
	}

	if subc.CategoriaID() != subcDB.CategoriaID {
		t.Errorf("esperava CategoriaID %s, recebeu %s", subcDB.CategoriaID, subc.CategoriaID())
	}

	if !subc.CriadoEm().Equal(subcDB.CriadoEm) {
		t.Errorf("esperava CriadoEm %v, recebeu %v", subcDB.CriadoEm, subc.CriadoEm())
	}

	if !subc.AtualizadoEm().Equal(subcDB.AtualizadoEm) {
		t.Errorf("esperava AtualizadoEm %v, recebeu %v", subcDB.AtualizadoEm, subc.AtualizadoEm())
	}
}

// TestValidarSubcategoria testa o método Validar da subcategoria.
func TestValidarSubcategoria(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		subc    Subcategoria
		wantErr bool
	}{
		{
			name: "subcategoria válida",
			subc: Subcategoria{
				id:          "sub-123",
				nome:        "Redes",
				status:      true,
				categoriaID: "cat-456",
				criadoEm:    now,
				atualizadoEm: now,
			},
			wantErr: false,
		},
		{
			name: "erro quando id é vazio",
			subc: Subcategoria{
				id:          "",
				nome:        "Redes",
				status:      true,
				categoriaID: "cat-456",
				criadoEm:    now,
				atualizadoEm: now,
			},
			wantErr: true,
		},
		{
			name: "erro quando nome é muito curto",
			subc: Subcategoria{
				id:          "sub-1",
				nome:        "A",
				status:      true,
				categoriaID: "cat-456",
				criadoEm:    now,
				atualizadoEm: now,
			},
			wantErr: true,
		},
		{
			name: "erro quando categoriaID é vazio",
			subc: Subcategoria{
				id:          "sub-1",
				nome:        "Redes",
				status:      true,
				categoriaID: "",
				criadoEm:    now,
				atualizadoEm: now,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.subc.Validar()

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
		})
	}
}

// TestAtivarDesativarSubcategoria testa os métodos Ativar e Desativar da subcategoria.
func TestAtivarDesativarSubcategoria(t *testing.T) {
	now := time.Now()

	subc := Subcategoria{
		id:          "sub-123",
		nome:        "Periféricos",
		status:      true,
		categoriaID: "cat-456",
		criadoEm:    now,
		atualizadoEm: now,
	}

	// Desativar
	subc.Desativar()
	if subc.Status() {
		t.Errorf("esperava status false após desativar, mas recebeu true")
	}

	// Ativar
	subc.Ativar()
	if !subc.Status() {
		t.Errorf("esperava status true após ativar, mas recebeu false")
	}
}

// TestComDadosAtualizadosSubcategoria testa o método ComDadosAtualizados da subcategoria.
func TestComDadosAtualizadosSubcategoria(t *testing.T) {
	subcategoria, _ := Novo("sub-123", "Antivírus", "cat-456")

	time.Sleep(1 * time.Second) // garante que o tempo de atualizadoEm será diferente

	params := AtualizarParams{
		Nome:        ptrString("Firewall"),
		CategoriaID: ptrString("cat-789"),
		Status:      ptrBool(false),
	}

	subcategoriaAtualizada, err := subcategoria.ComDadosAtualizados(params)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar dados: %v", err)
	}

	if subcategoriaAtualizada.Nome() != *params.Nome {
		t.Errorf("esperava Nome atualizado para %s, recebeu %s", *params.Nome, subcategoriaAtualizada.Nome())
	}

	if subcategoriaAtualizada.CategoriaID() != *params.CategoriaID {
		t.Errorf("esperava CategoriaID atualizado para %s, recebeu %s", *params.CategoriaID, subcategoriaAtualizada.CategoriaID())
	}

	if subcategoriaAtualizada.Status() != *params.Status {
		t.Errorf("esperava Status atualizado para %v, recebeu %v", *params.Status, subcategoriaAtualizada.Status())
	}

	if !subcategoriaAtualizada.AtualizadoEm().After(subcategoriaAtualizada.CriadoEm()) {
		t.Errorf("esperava AtualizadoEm após CriadoEm")
	}
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// ptrString é uma função auxiliar para obter um ponteiro para uma string.
func ptrString(s string) *string {
	return &s
}

// ptrBool é uma função auxiliar para obter um ponteiro para um bool.
func ptrBool(b bool) *bool {
	return &b
}
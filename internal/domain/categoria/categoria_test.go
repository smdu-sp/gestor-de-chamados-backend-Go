package categoria

import (
	"testing"
	"time"
)

// TestCategoria_Novo testa a função Novo da categoria.
func TestCategoria_Novo(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		nome      string
		wantError bool
	}{
		{
			name:      "criar categoria válida",
			id:        "cat-123",
			nome:      "Suporte Técnico",
			wantError: false,
		},
		{
			name:      "falha ao criar categoria com ID vazio",
			id:        "",
			nome:      "Suporte Técnico",
			wantError: true,
		},
		{
			name:      "falha ao criar categoria com nome muito curto",
			id:        "cat-124",
			nome:      "AB",
			wantError: true,
		},
		{
			name:      "falha ao criar categoria com nome muito longo",
			id:        "cat-125",
			nome:      string(make([]byte, 101)), // 101 caracteres
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Novo(tt.id, tt.nome)
			if (err != nil) != tt.wantError {
				t.Errorf("Novo() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestCategoria_CarregarDoBD testa a função CarregarDoBD da categoria.
func TestCategoria_CarregarDoBD(t *testing.T) {
	now := time.Now()
	categoriaDB := CategoriaDB{
		ID:        "cat-123",
		Nome:      "Suporte Técnico",
		Status:    true,
		CriadoEm:  now,
		AtualizadoEm: now,
	}

	categoria := CarregarDoBD(categoriaDB)

	if categoria.id != categoriaDB.ID {
		t.Errorf("ID incorreto: got %v, want %v", categoria.id, categoriaDB.ID)
	}

	if categoria.nome != categoriaDB.Nome {
		t.Errorf("Nome incorreto: got %v, want %v", categoria.nome, categoriaDB.Nome)
	}

	if categoria.status != categoriaDB.Status {
		t.Errorf("Status incorreto: got %v, want %v", categoria.status, categoriaDB.Status)
	}

	if !categoria.criadoEm.Equal(categoriaDB.CriadoEm) {
		t.Errorf("CriadoEm incorreto: got %v, want %v", categoria.criadoEm, categoriaDB.CriadoEm)
	}

	if !categoria.atualizadoEm.Equal(categoriaDB.AtualizadoEm) {
		t.Errorf("AtualizadoEm incorreto: got %v, want %v", categoria.atualizadoEm, categoriaDB.AtualizadoEm)
	}
}

// TestCategoria_Validar testa o método Validar da categoria.
func TestCategoria_Validar(t *testing.T) {
	tests := []struct {
		name      string
		categoria Categoria
		wantError bool
	}{
		{
			name: "categoria válida",
			categoria: Categoria{
				id:   "cat-123",
				nome: "Suporte Técnico",
			},
			wantError: false,
		},
		{
			name: "categoria com ID vazio",
			categoria: Categoria{
				id:   "",
				nome: "Suporte Técnico",
			},
			wantError: true,
		},
		{
			name: "categoria com nome muito curto",
			categoria: Categoria{
				id:   "cat-124",
				nome: "AB",
			},
			wantError: true,
		},
		{
			name: "categoria com nome muito longo",
			categoria: Categoria{
				id:   "cat-125",
				nome: string(make([]byte, 101)), // 101 caracteres
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.categoria.Validar()
			if (err != nil) != tt.wantError {
				t.Errorf("Validar() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestCategoria_AtivarDesativar testa os métodos Ativar e Desativar da categoria.
func TestCategoria_AtivarDesativar(t *testing.T) {
	categoria := Categoria{
		id:     "cat-123",
		nome:   "Suporte Técnico",
		status: false,
	}

	categoria.Ativar()
	if !categoria.status {
		t.Errorf("Ativar() falhou: status esperado true, got false")
	}

	categoria.Desativar()
	if categoria.status {
		t.Errorf("Desativar() falhou: status esperado false, got true")
	}
}

// TestCategoria_AtualizarDados testa o método AtualizarDados da categoria.
func TestCategoria_AtualizarDados(t *testing.T) {
	categoria := Categoria{
		id:     "cat-123",
		nome:   "Suporte Técnico",
		status: true,
	}

	time.Sleep(1 * time.Second) // garantir que o timestamp de atualizadoEm mude

	params := AtualizarParams{
		Nome:   ptrString("Suporte Avançado"),
		Status: nil,
	}

	err := categoria.AtualizarDados(params)
	if err != nil {
		t.Fatalf("AtualizarDados() retornou erro inesperado: %v", err)
	}

	if categoria.nome != "Suporte Avançado" {
		t.Errorf("AtualizarDados() falhou: nome esperado 'Suporte Avançado', got %v", categoria.nome)
	}

	if categoria.atualizadoEm.Equal(categoria.CriadoEm()) {
		t.Errorf("AtualizarDados() falhou: atualizadoEm não foi atualizado")
	}
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// ptrString é uma função auxiliar para obter um ponteiro para uma string.
func ptrString(s string) *string {
	return &s
}
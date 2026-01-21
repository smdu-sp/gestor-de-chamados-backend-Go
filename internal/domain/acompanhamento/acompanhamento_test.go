package acompanhamento

import (
	"testing"
	"time"

	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// TestAcompanhamento_Novo testa a função Novo do acompanhamento.
func TestAcompanhamento_Novo(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		chamadoID   string
		usuarioID   string
		conteudo    string
		remetente   usr.Permissao
		expectError bool
	}{
		{
			name:        "criação válida de acompanhamento",
			id:          "acomp-123",
			chamadoID:   "chamado-456",
			usuarioID:   "user-789",
			conteudo:    "Este é um acompanhamento válido.",
			remetente:   usr.PermTEC,
			expectError: false,
		},
		{
			name:        "conteúdo vazio",
			id:          "acomp-124",
			chamadoID:   "chamado-457",
			usuarioID:   "user-790",
			conteudo:    "",
			remetente:   usr.PermADM,
			expectError: true,
		},
		{
			name:        "remetente inválido",
			id:          "acomp-125",
			chamadoID:   "chamado-458",
			usuarioID:   "user-791",
			conteudo:    "Conteúdo válido.",
			remetente:   usr.Permissao("INVALIDO"), // Remetente inválido
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Novo(tt.id, tt.chamadoID, tt.usuarioID, tt.conteudo, tt.remetente)
			if (err != nil) != tt.expectError {
				t.Errorf("esperado erro: %v, obtido: %v", tt.expectError, err)
			}
		})
	}
}

// TestAcompanhamento_CarregarDoDB testa a função CarregarDoDB do acompanhamento.
func TestAcompanhamento_CarregarDoDB(t *testing.T) {
	now := time.Now()
	acompDB := AcompanhamentoDB{
		ID:           "acomp-123",
		ChamadoID:    "chamado-456",
		UsuarioID:    "user-789",
		Conteudo:     "Conteúdo do acompanhamento.",
		Remetente:    string(usr.PermTEC),
		CriadoEm:     now,
		AtualizadoEm: now,
	}

	acomp := CarregarDoDB(acompDB)

	if acomp.id != acompDB.ID {
		t.Errorf("esperado id: %s, obtido: %s", acompDB.ID, acomp.id)
	}
	if acomp.chamadoID != acompDB.ChamadoID {
		t.Errorf("esperado chamadoID: %s, obtido: %s", acompDB.ChamadoID, acomp.chamadoID)
	}
	if acomp.usuarioID != acompDB.UsuarioID {
		t.Errorf("esperado usuarioID: %s, obtido: %s", acompDB.UsuarioID, acomp.usuarioID)
	}
	if acomp.conteudo != acompDB.Conteudo {
		t.Errorf("esperado conteudo: %s, obtido: %s", acompDB.Conteudo, acomp.conteudo)
	}
	if string(acomp.remetente) != acompDB.Remetente {
		t.Errorf("esperado remetente: %s, obtido: %s", acompDB.Remetente, acomp.remetente)
	}
	if !acomp.criadoEm.Equal(acompDB.CriadoEm) {
		t.Errorf("esperado criadoEm: %v, obtido: %v", acompDB.CriadoEm, acomp.criadoEm)
	}
	if !acomp.atualizadoEm.Equal(acompDB.AtualizadoEm) {
		t.Errorf("esperado atualizadoEm: %v, obtido: %v", acompDB.AtualizadoEm, acomp.atualizadoEm)
	}
}

// TestAcompanhamento_Validar testa o método Validar do acompanhamento.
func TestAcompanhamento_Validar(t *testing.T) {
	tests := []struct {
		name        string
		acomp       Acompanhamento
		expectError bool
	}{
		{
			name: "acompanhamento válido",
			acomp: Acompanhamento{
				id:        "acomp-123",
				chamadoID: "chamado-456",
				usuarioID: "user-789",
				conteudo:  "Conteúdo válido do acompanhamento.",
				remetente: usr.PermTEC,
			},
			expectError: false,
		},
		{
			name: "acompanhamento com id vazio",
			acomp: Acompanhamento{
				id:        "",
				chamadoID: "chamado-456",
				usuarioID: "user-789",
				conteudo:  "Conteúdo válido do acompanhamento.",
				remetente: usr.PermTEC,
			},
			expectError: true,
		},
		{
			name: "acompanhamento com conteúdo muito curto",
			acomp: Acompanhamento{
				id:        "acomp-124",
				chamadoID: "chamado-457",
				usuarioID: "user-790",
				conteudo:  "Curto",
				remetente: usr.PermADM,
			},
			expectError: true,
		},
		{
			name: "acompanhamento com remetente inválido",
			acomp: Acompanhamento{
				id:        "acomp-125",
				chamadoID: "chamado-458",
				usuarioID: "user-791",
				conteudo:  "Conteúdo válido do acompanhamento.",
				remetente: usr.Permissao("INVALIDO"),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.acomp.Validar()
			if (err != nil) != tt.expectError {
				t.Errorf("esperado erro: %v, obtido: %v", tt.expectError, err)
			}
		})
	}
}

// TestValidarRemetente testa a função ValidarRemetente.
func TestAcompanhamentoValidarRemetente(t *testing.T) {
	tests := []struct {
		remetente   usr.Permissao
		expectError bool
	}{
		{remetente: usr.PermTEC, expectError: false},
		{remetente: usr.PermADM, expectError: false},
		{remetente: usr.PermDEV, expectError: false},
		{remetente: usr.PermUSR, expectError: false},
		{remetente: usr.Permissao("INVALIDO"), expectError: true},
	}

	for _, tt := range tests {
		err := ValidarRemetente(tt.remetente)
		if (err != nil) != tt.expectError {
			t.Errorf("remetente: %s, esperado erro: %v, obtido: %v", tt.remetente, tt.expectError, err)
		}
	}
}

// TestAcompanhamento_AtualizarDados testa a função AtualizarDados do acompanhamento.
func TestAcompanhamento_AtualizarDados(t *testing.T) {
	acomp, _ := Novo("acomp-123", "chamado-456", "user-789", "Conteúdo inicial do acompanhamento.", usr.PermTEC)

	time.Sleep(1 * time.Second) // garante que o tempo de atualizadoEm será diferente

	params := AtualizarParams{
		Conteudo: "Conteúdo atualizado do acompanhamento.",
	}
	
	err := acomp.AtualizarDados(params)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar dados: %v", err)
	}

	if acomp.conteudo != params.Conteudo {
		t.Errorf("esperado conteudo: %s, obtido: %s", params.Conteudo, acomp.conteudo)
	}

	if acomp.atualizadoEm.Equal(acomp.criadoEm) {
		t.Errorf("esperado atualizadoEm diferente de criadoEm")
	}
}

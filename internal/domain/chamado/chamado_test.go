package chamado

import (
	"testing"
	"time"
)

// TestNovoChamado testa a criação de novos chamados.
func TestNovoChamado(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		categoriaID string
		subcatID    string
		criadorID   string
		titulo      string
		descricao   string
		wantErr     bool
	}{
		{
			name:        "cria chamado valido",
			id:          "chamado-123",
			categoriaID: "cat-123",
			subcatID:    "subcat-123",
			criadorID:   "user-123",
			titulo:      "Problema com o sistema",
			descricao:   "O sistema está apresentando erros ao tentar salvar.",
			wantErr:     false,
		},
		{
			name:        "erro com titulo curto",
			id:          "chamado-123",
			categoriaID: "cat-123",
			subcatID:    "subcat-123",
			criadorID:   "user-123",
			titulo:      "Curto",
			descricao:   "O sistema está apresentando erros ao tentar salvar.",
			wantErr:     true,
		},
		{
			name:        "erro com descricao curta",
			id:          "chamado-123",
			categoriaID: "cat-123",
			subcatID:    "subcat-123",
			criadorID:   "user-123",
			titulo:      "Problema com o sistema",
			descricao:   "Curto",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := Novo(tt.id, tt.categoriaID, tt.subcatID, tt.criadorID, tt.titulo, tt.descricao)

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if !tt.wantErr && c == nil {
				t.Fatalf("esperava chamado, mas recebeu nil")
			}

			if !tt.wantErr && c.ID() != tt.id {
				t.Fatalf("esperava ID %s, mas recebeu %s", tt.id, c.ID())
			}
		})
	}
}

// TestCarregarDoDB testa a função CarregarDoDB.
func TestCarregarDoDB(t *testing.T) {
	now := time.Now()
	solucao := "Reinicie o sistema."

	chamadoDB := ChamadoDB{
		ID:             "chamado-123",
		CategoriaID:    "cat-123",
		SubcategoriaID: "subcat-123",
		CriadorID:      "user-123",
		Titulo:         "Problema com o sistema",
		Descricao:      "O sistema está apresentando erros ao tentar salvar.",
		Status:         "Fechado",
		Arquivado:      false,
		Solucao:        &solucao,
		SolucionadoEm:  &now,
		FechadoEm:      &now,
		CriadoEm:       now,
		AtualizadoEm:   now,
	}

	chamado := CarregarDoDB(chamadoDB)

	if chamado.ID() != chamadoDB.ID {
		t.Errorf("esperava ID %s, mas recebeu %s", chamadoDB.ID, chamado.ID())
	}

	if chamado.Titulo() != chamadoDB.Titulo {
		t.Errorf("esperava Titulo %s, mas recebeu %s", chamadoDB.Titulo, chamado.Titulo())
	}

	if chamado.Descricao() != chamadoDB.Descricao {
		t.Errorf("esperava Descricao %s, mas recebeu %s", chamadoDB.Descricao, chamado.Descricao())
	}

	if chamado.Status() != StatusChamado(chamadoDB.Status) {
		t.Errorf("esperava Status %s, mas recebeu %s", chamadoDB.Status, chamado.Status())
	}

	if chamado.Solucao() == nil || *chamado.Solucao() != *chamadoDB.Solucao {
		t.Errorf("esperava Solucao %s, mas recebeu %v", *chamadoDB.Solucao, chamado.Solucao())
	}

	if chamado.SolucionadoEm() == nil || !chamado.SolucionadoEm().Equal(*chamadoDB.SolucionadoEm) {
		t.Errorf("esperava SolucionadoEm %v, mas recebeu %v", *chamadoDB.SolucionadoEm, chamado.SolucionadoEm())
	}

	if chamado.FechadoEm() == nil || !chamado.FechadoEm().Equal(*chamadoDB.FechadoEm) {
		t.Errorf("esperava FechadoEm %v, mas recebeu %v", *chamadoDB.FechadoEm, chamado.FechadoEm())
	}

	if !chamado.CriadoEm().Equal(chamadoDB.CriadoEm) {
		t.Errorf("esperava CriadoEm %v, mas recebeu %v", chamadoDB.CriadoEm, chamado.CriadoEm())
	}

	if !chamado.AtualizadoEm().Equal(chamadoDB.AtualizadoEm) {
		t.Errorf("esperava AtualizadoEm %v, mas recebeu %v", chamadoDB.AtualizadoEm, chamado.AtualizadoEm())
	}
}

// TestChamado_Validar testa o método Validar do Chamado.
func TestChamado_Validar(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		chamado Chamado
		wantErr bool
	}{
		{
			name: "chamado valido",
			chamado: Chamado{
				id:          "chamado-123",
				categoriaID: "cat-123",
				subcategoriaID: "subcat-123",
				criadorID:   "user-123",
				titulo:      "Problema com o sistema",
				descricao:   "O sistema está apresentando erros ao tentar salvar.",
				status: StatusAberto,
				arquivado: false,
				solucao: nil,
				solucionadoEm: nil,
				fechadoEm: nil,
				criadoEm: now,
				atualizadoEm: now,
			},
			wantErr: false,
		},
		{
			name: "erro com titulo curto",
			chamado: Chamado{
				id:          "chamado-123",
				categoriaID: "cat-123",
				subcategoriaID: "subcat-123",
				criadorID:	 "user-123",
				titulo:      "Curto",
				descricao:   "O sistema está apresentando erros ao tentar salvar.",
				status: StatusAberto,
				arquivado: false,
				solucao: nil,
				solucionadoEm: nil,
				fechadoEm: nil,
				criadoEm: now,
				atualizadoEm: now,
			},
			wantErr: true,
		},
		{
			name: "erro com descricao curta",
			chamado: Chamado{
				id:          "chamado-123",
				categoriaID: "cat-123",
				subcategoriaID: "subcat-123",
				criadorID:   "user-123",
				titulo:      "Problema com o sistema",
				descricao:   "Curto",
				status: StatusAberto,
				arquivado: false,
				solucao: nil,
				solucionadoEm: nil,
				fechadoEm: nil,
				criadoEm: now,
				atualizadoEm: now,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.chamado.Validar()

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestChamado_AtualizarSolucao testa o método AtualizarSolucao do Chamado.
func TestChamado_AtualizarSolucao(t *testing.T) {
	now := time.Now()
	chamado, _ := Novo("chamado-123", "cat-123", "subcat-123", "user-123", "Problema com o sistema", "O sistema está apresentando erros ao tentar salvar.")

	solucao := "Reinicie o sistema."
	err := chamado.AtualizarSolucao(solucao)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar solução: %v", err)
	}

	if chamado.Solucao() == nil || *chamado.Solucao() != solucao {
		t.Fatalf("esperava Solucao %s, mas recebeu %v", solucao, chamado.Solucao())
	}

	if chamado.SolucionadoEm() == nil || !chamado.SolucionadoEm().Equal(now) {
		t.Fatalf("esperava SolucionadoEm %v, mas recebeu %v", now, chamado.SolucionadoEm())
	}
}

// TestChamado_AtualizarStatus testa o método AtualizarStatus do Chamado.
func TestChamado_AtualizarStatus(t *testing.T) {
	chamado, _ := Novo("chamado-123", "cat-123", "subcat-123", "user-123", "Problema com o sistema", "O sistema está apresentando erros ao tentar salvar.")

	err := chamado.AtualizarStatus(StatusAberto, nil)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar status: %v", err)
	}

	if chamado.Status() != StatusAberto {
		t.Fatalf("esperava Status %s, mas recebeu %s", StatusAberto, chamado.Status())
	}

	err = chamado.AtualizarStatus(StatusFechado, nil)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar status: %v", err)
	}

	if chamado.Status() != StatusFechado {
		t.Fatalf("esperava Status %s, mas recebeu %s", StatusFechado, chamado.Status())
	}
}

// TestChamado_AtualizarDados testa o método AtualizarDados do Chamado.
func TestChamado_AtualizarDados(t *testing.T) {
	chamado, _ := Novo("chamado-123", "cat-123", "subcat-123", "user-123", "Problema com o sistema", "O sistema está apresentando erros ao tentar salvar.")

	time.Sleep(1 * time.Second) // garante que o tempo de atualizadoEm será diferente

	params := AtualizarParams{
		Titulo:    ptrString("Novo título do chamado"),
		Descricao: ptrString("Nova descrição detalhada do problema."),
		Arquivado: ptrBool(false),
		CategoriaID: ptrString("cat-456"),
		SubcategoriaID: ptrString("subcat-456"),
	}
	err := chamado.AtualizarDados(params)
	if err != nil {
		t.Fatalf("erro inesperado ao atualizar dados: %v", err)
	}

	if chamado.Titulo() != *params.Titulo {
		t.Fatalf("esperava Titulo %s, mas recebeu %s", *params.Titulo, chamado.Titulo())
	}

	if chamado.Descricao() != *params.Descricao {
		t.Fatalf("esperava Descricao %s, mas recebeu %s", *params.Descricao, chamado.Descricao())
	}

	if chamado.CategoriaID() != *params.CategoriaID {
		t.Fatalf("esperava CategoriaID %s, mas recebeu %s", *params.CategoriaID, chamado.CategoriaID())
	}

	if chamado.SubcategoriaID() != *params.SubcategoriaID {
		t.Fatalf("esperava SubcategoriaID %s, mas recebeu %s", *params.SubcategoriaID, chamado.SubcategoriaID())
	}

	if !chamado.AtualizadoEm().After(*chamado.CriadoEm()) {
		t.Fatalf("esperava AtualizadoEm após CriadoEm")
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
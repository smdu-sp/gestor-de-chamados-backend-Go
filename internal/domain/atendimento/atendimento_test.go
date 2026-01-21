package atendimento

import (
	"testing"
	"time"
)

// TestAtendimento_Novo testa a função Novo do Atendimento.
func TestAtendimento_Novo(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		atribuidoID string
		chamadoID   string
		wantErr     bool
	}{
		{
			name:        "criação válida de atendimento",
			id:          "atendimento-123",
			atribuidoID: "tecnico-456",
			chamadoID:   "chamado-789",
			wantErr:     false,
		},
		{
			name:        "erro ao criar atendimento com ID vazio",
			id:          "",
			atribuidoID: "tecnico-456",
			chamadoID:   "chamado-789",
			wantErr:     true,
		},
		{
			name:        "erro ao criar atendimento com AtribuidoID vazio",
			id:          "atendimento-123",
			atribuidoID: "",
			chamadoID:   "chamado-789",
			wantErr:     true,
		},
		{
			name:        "erro ao criar atendimento com ChamadoID vazio",
			id:          "atendimento-123",
			atribuidoID: "tecnico-456",
			chamadoID:   "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Novo(tt.id, tt.atribuidoID, tt.chamadoID)

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestAtendimento_CarregarDoBD testa a função CarregarDoBD do Atendimento.
func TestAtendimento_CarregarDoBD(t *testing.T) {
	now := time.Now()
	atendimentoDB := AtendimentoDB{
		ID:          "atendimento-123",
		AtribuidoID: "tecnico-456",
		ChamadoID:   "chamado-789",
		CriadoEm:    now,
		AtualizadoEm: now,
	}

	atendimento := CarregarDoBD(atendimentoDB)
	
	if atendimento.id != atendimentoDB.ID {
		t.Errorf("esperava ID %s, mas recebeu %s", atendimentoDB.ID, atendimento.id)
	}

	if atendimento.atribuidoID != atendimentoDB.AtribuidoID {
		t.Errorf("esperava AtribuidoID %s, mas recebeu %s", atendimentoDB.AtribuidoID, atendimento.atribuidoID)
	}

	if atendimento.chamadoID != atendimentoDB.ChamadoID {
		t.Errorf("esperava ChamadoID %s, mas recebeu %s", atendimentoDB.ChamadoID, atendimento.chamadoID)
	}

	if !atendimento.criadoEm.Equal(atendimentoDB.CriadoEm) {
		t.Errorf("esperava CriadoEm %v, mas recebeu %v", atendimentoDB.CriadoEm, atendimento.criadoEm)
	}
	
	if !atendimento.atualizadoEm.Equal(atendimentoDB.AtualizadoEm) {
		t.Errorf("esperava AtualizadoEm %v, mas recebeu %v", atendimentoDB.AtualizadoEm, atendimento.atualizadoEm)
	}
}

// TestAtendimento_Validar testa o método Validar do Atendimento.
func TestAtendimento_Validar(t *testing.T) {
	tests := []struct {
		name        string
		atendimento Atendimento
		wantErr     bool
	}{
		{
			name: "atendimento válido",
			atendimento: Atendimento{
				id:          "atendimento-123",
				atribuidoID: "tecnico-456",
				chamadoID:   "chamado-789",
			},
			wantErr: false,
		},
		{
			name: "erro ao validar atendimento com ID vazio",
			atendimento: Atendimento{
				id:          "",
				atribuidoID: "tecnico-456",
				chamadoID:   "chamado-789",
			},
			wantErr: true,
		},
		{
			name: "erro ao validar atendimento com AtribuidoID vazio",
			atendimento: Atendimento{
				id:          "atendimento-123",
				atribuidoID: "",
				chamadoID:   "chamado-789",
			},
			wantErr: true,
		},
		{
			name: "erro ao validar atendimento com ChamadoID vazio",
			atendimento: Atendimento{
				id:          "atendimento-123",
				atribuidoID: "tecnico-456",
				chamadoID:   "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.atendimento.Validar()

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestAtendimento_AtualizarDados testa o método AtualizarDados do Atendimento.
func TestAtendimento_AtualizarDados(t *testing.T) {
	atendimento, _ := Novo("atendimento-123", "tecnico-456", "chamado-789")

	time.Sleep(1 * time.Second) // garantir que o timestamp de atualizadoEm será diferente

	params := AtualizarParams{
		AtribuidoID: "tecnico-999",
		ChamadoID:   "chamado-888",
	}

	err := atendimento.AtualizarDados(params)
	if err != nil {
		t.Fatalf("erro ao atualizar dados: %v", err)
	}

	if atendimento.atribuidoID != params.AtribuidoID {
		t.Errorf("esperava AtribuidoID %s, mas recebeu %s", params.AtribuidoID, atendimento.atribuidoID)
	}

	if atendimento.chamadoID != params.ChamadoID {
		t.Errorf("esperava ChamadoID %s, mas recebeu %s", params.ChamadoID, atendimento.chamadoID)
	}

	if !atendimento.atualizadoEm.After(atendimento.criadoEm) {
		t.Errorf("esperava AtualizadoEm após CriadoEm, mas não foi o caso")
	}
}
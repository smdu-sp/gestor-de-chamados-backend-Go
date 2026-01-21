package chamado

import (
	"strings"
	"testing"
)

// TestValidarStatus testa a função ValidarStatus.
func TestValidarStatus(t *testing.T) {
	tests := []struct {
		name    string
		status  StatusChamado
		wantErr bool
	}{
		{
			name:    "status válido: ABERTO",
			status:  StatusAberto,
			wantErr: false,
		},
		{
			name:    "status válido: ATRIBUIDO",
			status:  StatusAtribuido,
			wantErr: false,
		},
		{
			name:    "status válido: RESOLVIDO",
			status:  StatusResolvido,
			wantErr: false,
		},
		{
			name:    "status válido: REJEITADO",
			status:  StatusRejeitado,
			wantErr: false,
		},
		{
			name:    "status válido: FECHADO",
			status:  StatusFechado,
			wantErr: false,
		},
		{
			name:    "status inválido",
			status:  "INVALIDO",
			wantErr: true,
		},
		{
			name:    "status vazio",
			status:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidarStatus(tt.status)
			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestStatusValidos_DeveRetornarTodosStatus testa se StatusValidos retorna todos os status definidos.
func TestStatusValidos_DeveRetornarTodosStatus(t *testing.T) {
	statuses := StatusValidos()

	esperados := []string{
		string(StatusAberto),
		string(StatusAtribuido),
		string(StatusResolvido),
		string(StatusRejeitado),
		string(StatusFechado),
	}

	for _, esperado := range esperados {
		if !strings.Contains(statuses, esperado) {
			t.Fatalf("esperava que StatusValidos contivesse %q, mas não contém", esperado)
		}
	}
}

// TestStatusChamad0String testa o método String dos status de chamado.
func TestStatusChamado_String(t *testing.T) {
	tests := []struct {
		nome   string
		status StatusChamado
	}{
		{"StatusAberto", StatusAberto},
		{"StatusAtribuido", StatusAtribuido},
		{"StatusResolvido", StatusResolvido},
		{"StatusRejeitado", StatusRejeitado},
		{"StatusFechado", StatusFechado},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			if tt.status.String() != string(tt.status) {
				t.Fatalf("esperava %q, mas recebeu %q", tt.status, tt.status.String())
			}
		})
	}
}
package log

import (
	"strings"
	"testing"
)

// TestValidarAcao_Sucesso testa a função ValidarAcao com ações válidas.
func TestValidarAcao_Sucesso(t *testing.T) {
	testes := []struct {
		nome string
		acao Acao
	}{
		{"Criar", Criar},
		{"Atualizar", Atualizar},
		{"Desativar", Desativar},
		{"Ativar", Ativar},
		{"Arquivar", Arquivar},
		{"Desarquivar", Desarquivar},
		{"Deletar", Deletar},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			err := ValidarAcao(tt.acao)
			if err != nil {
				t.Fatalf("não esperava erro, recebeu: %v", err)
			}
		})
	}
}

// TestValidarAcao_AcaoVazia testa a função ValidarAcao com ação vazia.
func TestValidarAcao_AcaoVazia(t *testing.T) {
	err := ValidarAcao("")

	if err == nil {
		t.Fatalf("esperava erro para ação vazia")
	}
}

// TestValidarAcao_AcaoInvalida testa a função ValidarAcao com ação inválida.
func TestValidarAcao_AcaoInvalida(t *testing.T) {
	acaoInvalida := Acao("INVALIDA")
	err := ValidarAcao(acaoInvalida)

	if err == nil {
		t.Fatalf("esperava erro para ação inválida")
	}
}

// TestAcoesValidas_DeveConterTodasAcoes testa se AcoesValidas retorna todas as ações definidas.
func TestAcoesValidas_DeveConterTodasAcoes(t *testing.T) {
	acoes := AcoesValidas()

	esperadas := []Acao{
		Criar,
		Atualizar,
		Desativar,
		Ativar,
		Arquivar,
		Desarquivar,
		Deletar,
	}

	for _, acao := range esperadas {
		if !strings.Contains(acoes, acao.String()) {
			t.Fatalf("esperava que AcoesValidas contivesse %q", acao)
		}
	}
}

// TestAcao_String testa o método String do tipo Acao.
func TestAcao_String(t *testing.T) {
	testes := []struct {
		nome string
		acao Acao
	}{
		{"Criar", Criar},
		{"Atualizar", Atualizar},
		{"Desativar", Desativar},
		{"Ativar", Ativar},
		{"Arquivar", Arquivar},
		{"Desarquivar", Desarquivar},
		{"Deletar", Deletar},
	}

	for _, tt := range testes {
		t.Run(tt.nome, func(t *testing.T) {
			resultado := tt.acao.String()
			if resultado != string(tt.acao) {
				t.Fatalf("esperava %q, recebeu %q", tt.acao, resultado)
			}
		})
	}
}

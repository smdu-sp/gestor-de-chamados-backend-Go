package service

import (
	"errors"
	"testing"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// --- Variáveis de erro para testes de serviço ---------------------------------------------------------------------------------------

var (
	errFakeGeradorID = errors.New("fake: erro gerador id")
	errFakeRepo      = errors.New("fake: erro repositório")
)

// =====================================================================================================================
//  FAKES
// =====================================================================================================================

// --- Fake GeradorID --------------------------------------------------------------------------------------------------

// fakeGeradorID é uma implementação falsa de GeradorID para testes.
type fakeGeradorID struct {
	id  string
	err error
}

// NovoID retorna um ID falso ou um erro, dependendo da configuração do fakeGeradorID.
func (f *fakeGeradorID) NovoID() (string, error) {
	return f.id, f.err
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// assertError é um helper que verifica se o erro está conforme esperado.
func assertError(t *testing.T, erroRecebido error, wantErr bool) {
	t.Helper()
	
	if wantErr && erroRecebido == nil {
		t.Fatal("esperava erro, mas recebeu nil")
	}
	
	if !wantErr && erroRecebido != nil {
		t.Fatalf("erro inesperado: %v", erroRecebido)
	}
}

// verificarErro é um helper que verifica se o erro recebido corresponde ao erro esperado.
func verificarErro(t *testing.T, erroRecebido error, erroEsperado error) {
	t.Helper()

	if erroEsperado == nil {
		if erroRecebido != nil {
			t.Fatalf("erro inesperado: %v", erroRecebido)
		}
		return
	}

	if erroRecebido == nil {
		t.Fatalf("esperava erro %v, recebeu nil", erroEsperado)
	}

	// Se for erro de validação → compara tipo
	var erroValidacao *dmn.ErrosValidacao
	if errors.As(erroEsperado, &erroValidacao) {
		if !errors.As(erroRecebido, &erroValidacao) {
			t.Fatalf("esperava erro de validação, recebeu %v", erroRecebido)
		}
		return
	}

	// Senão, erro sentinela
	if !errors.Is(erroRecebido, erroEsperado) {
		t.Fatalf("erro esperado %v, mas recebeu %v", erroEsperado, erroRecebido)
	}
}

// ptr é um helper que retorna um ponteiro para o valor fornecido.
func ptr[T any](v T) *T {
	return &v
}

// testLogWriter é um writer personalizado para capturar logs durante os testes.
type testLogWriter struct {
	writeFn func(msg string)
}

// Write implementa a interface io.Writer.
func (w *testLogWriter) Write(p []byte) (n int, err error) {
	w.writeFn(string(p))
	return len(p), nil
}

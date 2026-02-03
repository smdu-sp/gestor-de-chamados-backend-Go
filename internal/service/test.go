package service

import (
	"errors"
	"testing"
	"time"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// =====================================================================================================================
//  ERROS FALSOS
// =====================================================================================================================

var (
	errFakeGeradorID = errors.New("fake: erro gerador id") // Erro falso para gerador de ID
	errFakeRepo      = errors.New("fake: erro repositório") // Erro falso para repositório
)

// =====================================================================================================================
//  TIPOS AUXILIARES
// =====================================================================================================================

// comparadorFn é um tipo de função que compara dois valores do mesmo tipo e retorna true se forem iguais.
type comparadorFn[T any] func(a, b T) bool

// =====================================================================================================================
//  FAKES
// =====================================================================================================================

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
// Usado para validar erros retornados por funções em testes.
func verificarErro(t *testing.T, recebido, esperado error) {
	t.Helper()

	// Se não espera erro, mas recebeu
	if esperado == nil && recebido != nil {
		t.Fatalf("não esperava erro, mas recebeu %v", recebido)
	}

	// Se espera erro, mas não recebeu
	if esperado != nil && recebido == nil {
		t.Fatalf("esperava erro %v, mas recebeu nil", esperado)
	}

	// Se for erro de validação → compara tipo
	var erroValidacao *dmn.ErrosValidacao
	if errors.As(esperado, &erroValidacao) {
		if !errors.As(recebido, &erroValidacao) {
			t.Fatalf("esperava erro de validação %v, mas recebeu %v", esperado, recebido)
		}
		return
	}

	// Compara erros diretamente
	if !errors.Is(recebido, esperado) {
		t.Fatalf("esperava erro %v, mas recebeu %v", esperado, recebido)
	}
}

// verificarNaoNulo é um helper que verifica se o valor fornecido não é nulo.
// Usado para garantir que uma entidade foi criada ou retornada.
func verificarNaoNulo[T any](t *testing.T, valor *T, mensagem string) {
	t.Helper()

	if valor == nil {
		t.Fatal(mensagem)
	}
}

// verificarNulo é um helper que verifica se o valor fornecido é nulo.
// Usado para garantir que nenhuma entidade foi criada ou retornada.
func verificarNulo[T any](t *testing.T, valor *T, mensagem string) {
	t.Helper()

	if valor != nil {
		t.Error(mensagem)
	}
}

// verificarLimite é um helper que verifica se o tamanho do slice corresponde ao esperado.
// Usado para garantir que a quantidade de itens retornados por uma consulta está correta.
func verificarLimite[T any](t *testing.T, slice []T, esperado int, mensagem string) {
	t.Helper()

	if len(slice) != esperado {
		t.Fatalf("%s: esperado %d, obtido %d", mensagem, esperado, len(slice))
	}
}

// verificarTotal é um helper que verifica se o total retornado corresponde ao esperado.
// Usado para garantir que o total de itens retornados por uma consulta está correto.
func verificarTotal(t *testing.T, total, esperado int, mensagem string) {
	t.Helper()

	if total != esperado {
		t.Fatalf("%s: esperado total %d, obtido %d", mensagem, esperado, total)
	}
}

// verificarNaoAlterado é um helper que verifica se o valor antes e depois são iguais, usando a função de comparação fornecida.
// Usado para garantir que certas operações não alterem o estado de um objeto.
func verificarNaoAlterado[T any](t *testing.T, antes T, depois T, igual comparadorFn[T], mensagem string) {
	t.Helper()

	if !igual(antes, depois) {
		t.Fatal(mensagem)
	}
}

// verificarDatas é um helper que verifica se as datas de criação e atualização estão corretas.
// Garante que não são zero e que atualizadoEm não é anterior a criadoEm.
func verificarDatas(t *testing.T, criadoEm, atualizadoEm time.Time) {
	t.Helper()

	if criadoEm.IsZero() {
		t.Fatal("CriadoEm não deveria ser zero")
	}
	if atualizadoEm.IsZero() {
		t.Fatal("AtualizadoEm não deveria ser zero")
	}
	if atualizadoEm.Before(criadoEm) {
		t.Fatalf("AtualizadoEm(%v) não pode ser anterior a CriadoEm(%v)", atualizadoEm, criadoEm)
	}
}

// verificarNaoPersistido é um helper que verifica se nenhuma entidade foi persistida.
// Usado para garantir que operações que falharam não deixem efeitos colaterais.
func verificarNaoPersistido(t *testing.T, quantidade int) {
	t.Helper()

	if quantidade != 0 {
		t.Fatal("nenhuma entidade deveria ter sido persistida")
	}
}

// verificarPersistido é um helper que verifica se uma entidade foi persistida.
// Usado para garantir que operações bem-sucedidas resultem em persistência.
func verificarPersistido(t *testing.T, quantidade int) {
	t.Helper()

	if quantidade == 0 {
		t.Fatal("esperava entidade persistida, mas nenhuma foi")
	}
}

// verificarOk é um helper que verifica se ok é true, caso contrário falha o teste com a mensagem fornecida.
// Usado para validar condições que devem ser verdadeiras.
func verificarOk(t *testing.T, ok bool, mensagem string) {
	t.Helper()

	if !ok {
		t.Fatal(mensagem)
	}
}

// verificarIDs é um helper que compara dois IDs e falha o teste se eles não forem iguais.
// Usado para garantir que entidades esperadas e recebidas correspondam.
func verificarIDs(t *testing.T, esperado, recebido string) {
	if esperado != recebido {
		t.Errorf("ID esperado '%s', recebido '%s'", esperado, recebido)
	}
}

// verificarLogins é um helper que compara dois logins e falha o teste se eles não forem iguais.
// Usado para garantir que logins esperados e recebidos correspondam.
func verificarLogins(t *testing.T, esperado, recebido string) {
	if esperado != recebido {
		t.Errorf("Login esperado '%s', recebido '%s'", esperado, recebido)
	}
}

// verificarPaginacao é um helper que verifica se a página e o limite do filtro retornado correspondem aos esperados.
// Usado para garantir que a paginação esteja funcionando corretamente.
func verificarPaginacao(t *testing.T, filtroRetornado dmn.Paginacao, paginaEsperada, limiteEsperado int) {
	if filtroRetornado.Pagina() != paginaEsperada || filtroRetornado.Limite() != limiteEsperado {
		t.Errorf("esperava filtro retornado com página %d e limite %d, mas recebeu página %d e limite %d", 
		paginaEsperada, limiteEsperado, filtroRetornado.Pagina(), filtroRetornado.Limite())
	}
}

// ptr é um helper que retorna um ponteiro para o valor fornecido.
// Usado para facilitar a criação de ponteiros em testes.
func ptr[T any](valor T) *T {
	return &valor
}

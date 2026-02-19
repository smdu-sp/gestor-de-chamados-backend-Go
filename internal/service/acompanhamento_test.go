package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que AcompanhamentoRepositoryFake implementa AcompanhamentoRepository
var _ acp.Repository = (*fakeAcompanhamentoRepository)(nil)

// fakeAcompanhamentoRepository é um repositório falso para testes de AcompanhamentoService. Ele armazena os acompanhamentos em memória.
type fakeAcompanhamentoRepository struct {
	acompanhamentos map[string]*acp.Acompanhamento
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoDeletar   error
	erroAoListar    error
}

// novoAcompanhamentoRepositoryFake cria uma nova instância de AcompanhamentoRepositoryFake.
func novoFakeAcompanhamentoRepository() *fakeAcompanhamentoRepository {
	return &fakeAcompanhamentoRepository{
		acompanhamentos: make(map[string]*acp.Acompanhamento),
	}
}

// Criar armazena um novo acompanhamento no repositório falso.
func (r *fakeAcompanhamentoRepository) Criar(ctx context.Context, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	if r.erroAoCriar != nil {
		return nil, r.erroAoCriar
	}
	acompanhamento := a
	r.acompanhamentos[a.ID()] = &acompanhamento
	return &acompanhamento, nil
}

// Atualizar atualiza um acompanhamento existente no repositório falso.
func (r *fakeAcompanhamentoRepository) Atualizar(ctx context.Context, id string, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	if r.erroAoAtualizar != nil {
		return nil, r.erroAoAtualizar
	}
	acompanhamento := a
	r.acompanhamentos[a.ID()] = &acompanhamento
	return &acompanhamento, nil
}

// BuscarPorID busca um acompanhamento pelo seu ID no repositório falso.
func (r *fakeAcompanhamentoRepository) BuscarPorID(ctx context.Context, id string) (*acp.Acompanhamento, error) {
	if r.erroAoBuscar != nil {
		return nil, r.erroAoBuscar
	}
	a, existe := r.acompanhamentos[id]
	if !existe {
		return nil, mysql.ErrAcompanhamentoNaoEncontrado
	}
	return a, nil
}

// BuscarPorChamadoID busca os acompanhamentos de um chamado pelo ID do chamado no repositório falso.
func (r *fakeAcompanhamentoRepository) BuscarPorChamadoID(ctx context.Context, chamadoID string) ([]acp.Acompanhamento, error) {
	if r.erroAoListar != nil {
		return nil, r.erroAoListar
	}
	var acompanhamentos []acp.Acompanhamento
	for _, a := range r.acompanhamentos {
		if a.ChamadoID() == chamadoID {
			acompanhamentos = append(acompanhamentos, *a)
		}
	}
	// Ordenar por data de criação (do mais antigo para o mais recente)
	sort.Slice(acompanhamentos, func(i, j int) bool {
		return acompanhamentos[i].CriadoEm().Before(acompanhamentos[j].CriadoEm())
	})
	return acompanhamentos, nil
}

// Deletar remove um acompanhamento do repositório falso.
func (r *fakeAcompanhamentoRepository) Deletar(ctx context.Context, id string) error {
	if r.erroAoDeletar != nil {
		return r.erroAoDeletar
	}
	delete(r.acompanhamentos, id)
	return nil
}

// Listar lista todos os acompanhamentos do repositório falso.
func (r *fakeAcompanhamentoRepository) Listar(ctx context.Context, filtro acp.Filtro) ([]acp.Acompanhamento, int, error) {
	if r.erroAoListar != nil {
		return nil, 0, r.erroAoListar
	}

	// Converter o mapa para um slice
	acompanhamentos := make([]acp.Acompanhamento, 0, len(r.acompanhamentos))
	for _, a := range r.acompanhamentos {
		acompanhamentos = append(acompanhamentos, *a)
	}

	// Ordenação ORDER BY criado_em ASC
	sort.Slice(acompanhamentos, func(i, j int) bool {
		return acompanhamentos[i].CriadoEm().Before(acompanhamentos[j].CriadoEm())
	})

	// Contar o total antes da paginação
	total := len(acompanhamentos)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio > total {
			return []acp.Acompanhamento{}, total, nil
		}

		if fim > total {
			fim = total
		}

		acompanhamentos = acompanhamentos[inicio:fim]
	}

	return acompanhamentos, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const acompanhamentoTesteID = "acompanhamento-123"

// novoAcompanhamentoTeste cria uma nova instância de Acompanhamento para ser usada em testes.
func novoAcompanhamentoTeste() *acp.Acompanhamento {
	a, _ := acp.Novo(
		acompanhamentoTesteID,
		chamadoTesteID,
		tecnicoTesteID,
		"Conteúdo do acompanhamento de teste",
		usr.PermUSR,
	)
	return a
}

// popularAcompanhamentoTeste insere um acompanhamento de teste no repositório falso.
func popularAcompanhamentoTeste(repo *fakeAcompanhamentoRepository, quantidade int) {
	for i := 0; i < quantidade; i++ {
		id := fmt.Sprintf("acompanhamento-%02d", i)
		chamadoTesteID := fmt.Sprintf("chamado-%02d", i)
		tecnicoTesteID := fmt.Sprintf("tecnico-%02d", i)

		a, _ := acp.Novo(
			id,
			chamadoTesteID,
			tecnicoTesteID,
			fmt.Sprintf("Conteúdo do acompanhamento de teste %d", i),
			usr.PermUSR,
		)
		repo.acompanhamentos[id] = a
	}
}

// novoCriarAcompanhamentoParamsTeste cria um conjunto de parâmetros válidos para testar a criação de um acompanhamento.
func novoCriarAcompanhamentoParamsTeste() acp.CriarParams {
	return acp.CriarParams{
		Conteudo:  "Conteúdo do acompanhamento de teste",
		ChamadoID: chamadoTesteID,
	}
}

// novoCriarAcompanhamentoParamsInvalidosTeste cria um conjunto de parâmetros inválidos para testar a criação de um acompanhamento.
func novoCriarAcompanhamentoParamsInvalidosTeste() acp.CriarParams {
	return acp.CriarParams{
		Conteudo:  "",
		ChamadoID: chamadoTesteID,
	}
}

// novoAtualizarAcompanhamentoParamsTeste cria um conjunto de parâmetros válidos para testar a atualização de um acompanhamento.
func novoAtualizarAcompanhamentoParamsTeste() acp.AtualizarParams {
	return acp.AtualizarParams{
		Conteudo: "Conteúdo atualizado do acompanhamento de teste",
	}
}

// novoAtualizarAcompanhamentoParamsInvalidosTeste cria um conjunto de parâmetros inválidos para testar a atualização de um acompanhamento.
func novoAtualizarAcompanhamentoParamsInvalidosTeste() acp.AtualizarParams {
	return acp.AtualizarParams{
		Conteudo: "",
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_AcompanhamentoService_Criar testa o caso de uso de criação de acompanhamento.
func Test_AcompanhamentoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() acp.Repository
		params             acp.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento)
	}{
		{
			nome: "criar acompanhamento com sucesso",
			geradorID: &fakeGeradorID{
				id: acompanhamentoTesteID,
			},
			prepararRepo: func() acp.Repository {
				return novoFakeAcompanhamentoRepository()
			},
			params:       novoCriarAcompanhamentoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoNulo(t, acompanhamento)
				verificarPersistido(t, len(repo.(*fakeAcompanhamentoRepository).acompanhamentos))

				if acompanhamento.Conteudo() != novoCriarAcompanhamentoParamsTeste().Conteudo {
					t.Errorf("conteúdo esperado '%s', mas obteve '%s'",
						novoCriarAcompanhamentoParamsTeste().Conteudo, acompanhamento.Conteudo())
				}

				verificarIDs(t, acompanhamentoTesteID, acompanhamento.ID())
				verificarIDs(t, chamadoTesteID, acompanhamento.ChamadoID())
				verificarIDs(t, tecnicoTesteID, acompanhamento.UsuarioID())
				verificarDatas(t, acompanhamento.CriadoEm(), acompanhamento.AtualizadoEm())
			},
		},
		{
			nome: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() acp.Repository {
				return novoFakeAcompanhamentoRepository()
			},
			params:       novoCriarAcompanhamentoParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeAcompanhamentoRepository).acompanhamentos))
				verificarNulo(t, acompanhamento)
			},
		},
		{
			nome: "erro ao salvar acompanhamento no repositório",
			geradorID: &fakeGeradorID{
				id: acompanhamentoTesteID,
			},
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarAcompanhamentoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeAcompanhamentoRepository).acompanhamentos))
				verificarNulo(t, acompanhamento)
			},
		},
		{
			nome: "erro de validação ao criar acompanhamento",
			geradorID: &fakeGeradorID{
				id: acompanhamentoTesteID,
			},
			prepararRepo: func() acp.Repository {
				return novoFakeAcompanhamentoRepository()
			},
			params:       novoCriarAcompanhamentoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeAcompanhamentoRepository).acompanhamentos))
				verificarNulo(t, acompanhamento)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			novoChamado := novoChamadoTeste()
			chamadoRepo.chamados[chamadoTesteID] = novoChamado

			atendimentoRepo := novoFakeAtendimentoRepository()
			novoAtendimento := novoAtendimentoTeste()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimento
			service := NovoAcompanhamentoService(tt.geradorID, repo, chamadoRepo, atendimentoRepo)
			acompanhamentoCriado, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, acompanhamentoCriado)
			}
		})
	}
}

// Test_AcompanhamentoService_Atualizar testa o caso de uso de atualização de acompanhamento.
func Test_AcompanhamentoService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() acp.Repository
		id                 string
		params             acp.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento)
	}{
		{
			nome: "atualizar acompanhamento com sucesso",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				acompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = acompanhamento
				return repo
			},
			id:           acompanhamentoTesteID,
			params:       novoAtualizarAcompanhamentoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoNulo(t, acompanhamento)

				if acompanhamento.Conteudo() != novoAtualizarAcompanhamentoParamsTeste().Conteudo {
					t.Errorf("conteúdo esperado '%s', mas obteve '%s'",
						novoAtualizarAcompanhamentoParamsTeste().Conteudo, acompanhamento.Conteudo())
				}

				verificarIDs(t, acompanhamentoTesteID, acompanhamento.ID())
				verificarIDs(t, chamadoTesteID, acompanhamento.ChamadoID())
				verificarIDs(t, tecnicoTesteID, acompanhamento.UsuarioID())
				verificarDatas(t, acompanhamento.CriadoEm(), acompanhamento.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar acompanhamento no repositório para atualizar",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           acompanhamentoTesteID,
			params:       novoAtualizarAcompanhamentoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				acompanhamentoRepo, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
				verificarNaoAlterado(t, novoAcompanhamentoTeste(), acompanhamentoRepo, compararAcompanhamentos)
				verificarDatas(t, acompanhamentoRepo.CriadoEm(), acompanhamentoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir acompanhamento atualizado no repositório",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           acompanhamentoTesteID,
			params:       novoAtualizarAcompanhamentoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				acompanhamentoRepo, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
				verificarNaoAlterado(t, novoAcompanhamentoTeste(), acompanhamentoRepo, compararAcompanhamentos)
				verificarDatas(t, acompanhamentoRepo.CriadoEm(), acompanhamentoRepo.AtualizadoEm())
			},
		},
		{
			nome: "acompanhamento não encontrado para atualizar",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoBuscar = mysql.ErrAcompanhamentoNaoEncontrado
				return repo
			},
			id:           acompanhamentoTesteID,
			params:       novoAtualizarAcompanhamentoParamsTeste(),
			erroEsperado: mysql.ErrAcompanhamentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				acompanhamentoRepo, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
				verificarNaoAlterado(t, novoAcompanhamentoTeste(), acompanhamentoRepo, compararAcompanhamentos)
				verificarDatas(t, acompanhamentoRepo.CriadoEm(), acompanhamentoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro de validação ao atualizar acompanhamento",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				return repo
			},
			id:           acompanhamentoTesteID,
			params:       novoAtualizarAcompanhamentoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				acompanhamentoRepo, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
				verificarNaoAlterado(t, novoAcompanhamentoTeste(), acompanhamentoRepo, compararAcompanhamentos)
				verificarDatas(t, acompanhamentoRepo.CriadoEm(), acompanhamentoRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			novoChamado := novoChamadoTeste()
			chamadoRepo.chamados[chamadoTesteID] = novoChamado
			atendimentoRepo := novoFakeAtendimentoRepository()
			novoAtendimento := novoAtendimentoTeste()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimento
			service := NovoAcompanhamentoService(nil, repo, chamadoRepo, atendimentoRepo)
			acompanhamentoAtualizado, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, acompanhamentoAtualizado)
			}
		})
	}
}

// Test_AcompanhamentoService_BuscarPorID testa o caso de uso de busca de acompanhamento por ID.
func Test_AcompanhamentoService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() acp.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento)
	}{
		{
			nome: "buscar acompanhamento por ID com sucesso",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				acompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = acompanhamento
				return repo
			},
			id:           acompanhamentoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				verificarNaoNulo(t, acompanhamento)
				verificarIDs(t, acompanhamentoTesteID, acompanhamento.ID())
			},
		},
		{
			nome: "erro ao buscar acompanhamento por id no repositório",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           acompanhamentoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				_, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
			},
		},
		{
			nome: "acompanhamento não encontrado ao buscar por id",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoBuscar = mysql.ErrAcompanhamentoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrAcompanhamentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamento *acp.Acompanhamento) {
				t.Helper()

				_, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
				verificarNulo(t, acompanhamento)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			novoChamado := novoChamadoTeste()
			chamadoRepo.chamados[chamadoTesteID] = novoChamado
			atendimentoRepo := novoFakeAtendimentoRepository()
			novoAtendimento := novoAtendimentoTeste()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimento
			service := NovoAcompanhamentoService(nil, repo, chamadoRepo, atendimentoRepo)
			acompanhamento, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, acompanhamento)
			}
		})
	}
}

// Test_AcompanhamentoService_Deletar testa o caso de uso de deleção de acompanhamento por ID.
func Test_AcompanhamentoService_Deletar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome         string
		prepararRepo func() acp.Repository
		id           string
		erroEsperado error
		verificarResultado func(t *testing.T, repo acp.Repository)
	}{
		{
			nome: "deletar acompanhamento com sucesso",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				acompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = acompanhamento
				return repo
			},
			id:           acompanhamentoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo acp.Repository) {
				t.Helper()

				_, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				if ok {
					t.Errorf("acompanhamento deveria ter sido deletado do repositório, mas ainda existe")
				}
			},
		},
		{
			nome: "acompanhamento não encontrado para deletar",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoBuscar = mysql.ErrAcompanhamentoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrAcompanhamentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo acp.Repository) {
				t.Helper()

				_, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
			},
		},
		{
			nome: "erro ao deletar acompanhamento no repositório",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				novoAcompanhamento := novoAcompanhamentoTeste()
				repo.acompanhamentos[acompanhamentoTesteID] = novoAcompanhamento
				repo.erroAoDeletar = errFakeRepo
				return repo
			},
			id:           acompanhamentoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository) {
				t.Helper()

				_, ok := repo.(*fakeAcompanhamentoRepository).acompanhamentos[acompanhamentoTesteID]
				verificarOk(t, ok, "deveria existir acompanhamento no repositório")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			novoChamado := novoChamadoTeste()
			chamadoRepo.chamados[chamadoTesteID] = novoChamado
			atendimentoRepo := novoFakeAtendimentoRepository()
			novoAtendimento := novoAtendimentoTeste()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimento
			service := NovoAcompanhamentoService(nil, repo, chamadoRepo, atendimentoRepo)
			erroRecebido := service.Deletar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo)
			}
		})
	}
}

// Test_AcompanhamentoService_Listar testa o caso de uso de listagem de acompanhamentos por ID do chamado.
func Test_AcompanhamentoService_Listar(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome         string
		prepararRepo func() acp.Repository
		filtro 		 acp.Filtro
		erroEsperado error
		verificarResultado func(t *testing.T, repo acp.Repository, acompanhamentos []acp.Acompanhamento, total int, filtroRetornado acp.Filtro)
	}{
		{
			nome: "listar 10 acompanhamentos paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				popularAcompanhamentoTeste(repo, 10)
				return repo
			},
			filtro: acp.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamentos []acp.Acompanhamento, total int, filtroRetornado acp.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeAcompanhamentoRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.acompanhamentos))
				verificarLimite(t, acompanhamentos, 5)
			},
		},
		{
			nome: "erro ao listar acompanhamentos paginados no repositório",
			prepararRepo: func() acp.Repository {
				repo := novoFakeAcompanhamentoRepository()
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro: acp.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo acp.Repository, acompanhamentos []acp.Acompanhamento, total int, filtroRetornado acp.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, acompanhamentos, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			novoChamado := novoChamadoTeste()
			chamadoRepo.chamados[chamadoTesteID] = novoChamado
			atendimentoRepo := novoFakeAtendimentoRepository()
			novoAtendimento := novoAtendimentoTeste()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimento
			service := NovoAcompanhamentoService(nil, repo, chamadoRepo, atendimentoRepo)
			acompanhamentos, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, acompanhamentos, total, filtroRetornado)
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

// compararAcompanhamentos compara dois acompanhamentos e retorna uma mensagem de erro se eles forem diferentes.
func compararAcompanhamentos(t *testing.T, antes, depois *acp.Acompanhamento) {
	t.Helper()

	if antes == nil || depois == nil {
		t.Fatalf("acompanhamentos para comparação não podem ser nil: antes=%v depois=%v", antes, depois)
	}

	if antes.ID() != depois.ID() ||
		antes.ChamadoID() != depois.ChamadoID() ||
		antes.UsuarioID() != depois.UsuarioID() ||
		antes.Conteudo() != depois.Conteudo() ||
		antes.Remetente() != depois.Remetente() {

		t.Fatalf("acompanhamentos não deveriam ser diferentes:\nantes=%+v\ndepois=%+v",
			antes, depois,
		)
	}
}

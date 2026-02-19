package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que fakeAtendimentoRepository implementa atd.Repository
var _ atd.Repository = (*fakeAtendimentoRepository)(nil)

// fakeAtendimentoRepository é uma implementação fake do repositório de atendimentos para testes.
type fakeAtendimentoRepository struct {
	atendimentos    map[string]*atd.Atendimento
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
}

func novoFakeAtendimentoRepository() *fakeAtendimentoRepository {
	return &fakeAtendimentoRepository{
		atendimentos: make(map[string]*atd.Atendimento),
	}
}

// Criar adiciona um novo atendimento ao repositório fake.
func (f *fakeAtendimentoRepository) Criar(ctx context.Context, a atd.Atendimento) (*atd.Atendimento, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	atendimento := a
	f.atendimentos[a.ID()] = &atendimento
	return &atendimento, nil
}

// Atualizar atualiza um atendimento existente no repositório fake.
func (f *fakeAtendimentoRepository) Atualizar(ctx context.Context, id string, a atd.Atendimento) (*atd.Atendimento, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	atendimento := a
	f.atendimentos[id] = &atendimento
	return &atendimento, nil
}

// BuscarPorID busca um atendimento pelo ID no repositório fake.
func (f *fakeAtendimentoRepository) BuscarPorID(ctx context.Context, id string) (*atd.Atendimento, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	a, existe := f.atendimentos[id]
	if !existe {
		return nil, mysql.ErrAtendimentoNaoEncontrado
	}
	return a, nil
}

// BuscarPorChamadoEAtribuidoID busca um atendimento pelo ID do chamado e do atribuído no repositório fake.
func (f *fakeAtendimentoRepository) BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	for _, a := range f.atendimentos {
		if a.ChamadoID() == chamadoID && a.AtribuidoID() == atribuidoID {
			return a, nil
		}
	}
	return nil, mysql.ErrAtendimentoNaoEncontrado
}

// Listar lista atendimentos com base no filtro fornecido no repositório fake.
func (f *fakeAtendimentoRepository) Listar(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para uma slice
	atendimentos := make([]atd.Atendimento, 0, len(f.atendimentos))
	for _, a := range f.atendimentos {
		atendimentos = append(atendimentos, *a)
	}

	// Ordenação ORDER BY criadoEm DESC
	sort.Slice(atendimentos, func(i, j int) bool {
		return atendimentos[i].CriadoEm().After(atendimentos[j].CriadoEm())
	})

	// Contar o total antes da paginação
	total := len(atendimentos)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio > total {
			return []atd.Atendimento{}, total, nil
		}

		if fim > total {
			fim = total
		}

		atendimentos = atendimentos[inicio:fim]
	}

	return atendimentos, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const atendimentoTesteID = "atendimento-123"

// novoAtendimentoTeste cria um novo atendimento para testes.
func novoAtendimentoTeste() *atd.Atendimento {
	a, _ := atd.Novo(
		atendimentoTesteID,
		tecnicoTesteID,
		chamadoTesteID,
	)
	return a
}

// popuplarAtendimentosTeste popula o repositório fake com atendimentos para testes.
func popuplarAtendimentosTeste(repo *fakeAtendimentoRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("atendimento-%02d", i)
		tecnicoID := fmt.Sprintf("tecnico-%02d", i)
		chamadoID := fmt.Sprintf("chamado-%02d", i)

		a, _ := atd.Novo(
			id,
			tecnicoID,
			chamadoID,
		)
		repo.atendimentos[id] = a
	}
}

// novoCriarAtendimentoParamsTeste cria parâmetros de criação de atendimento para testes.
func novoCriarAtendimentoParamsTeste() atd.CriarParams {
	return atd.CriarParams{
		AtribuidoID: tecnicoTesteID,
		ChamadoID:   chamadoTesteID,
	}
}

// novoAtualizarAtendimentoParamsTeste cria parâmetros de atualização de atendimento para testes.
func novoAtualizarAtendimentoParamsTeste() atd.AtualizarParams {
	return atd.AtualizarParams{
		AtribuidoID: tecnicoTesteID,
		ChamadoID:   chamadoTesteID,
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_AtendimentoService_Criar testa o método Criar do serviço de atendimento.
func Test_AtendimentoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() atd.Repository
		params             atd.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento)
	}{
		{
			nome: "criar atendimento com sucesso",
			geradorID: &fakeGeradorID{
				id: atendimentoTesteID,
			},
			prepararRepo: func() atd.Repository {
				return novoFakeAtendimentoRepository()
			},
			params:       novoCriarAtendimentoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoNulo(t, atendimento)
				verificarPersistido(t, len(repo.(*fakeAtendimentoRepository).atendimentos))
				verificarIDs(t, atendimentoTesteID, atendimento.ID())
				verificarIDs(t, tecnicoTesteID, atendimento.AtribuidoID())
				verificarIDs(t, chamadoTesteID, atendimento.ChamadoID())
				verificarDatas(t, atendimento.CriadoEm(), atendimento.AtualizadoEm())
			},
		},
		{
			nome: "erro ao salvar atendimento no repositório",
			geradorID: &fakeGeradorID{
				id: atendimentoTesteID,
			},
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarAtendimentoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeAtendimentoRepository).atendimentos))
				verificarNulo(t, atendimento)
			},
		},
		{
			nome: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() atd.Repository {
				return novoFakeAtendimentoRepository()
			},
			params:       novoCriarAtendimentoParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeAtendimentoRepository).atendimentos))
				verificarNulo(t, atendimento)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			chamadoRepo.chamados[chamadoTesteID] = novoChamadoTeste()

			categoriaPermissaoRepo := novoFakeCategoriaPermissaoRepository()
			categoriaPermissaoRepo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()

			usuarioRepo := novoFakeUsuarioRepository()
			novoTecnico := novoTecnicoTeste()
			usuarioRepo.usuarios[tecnicoTesteID] = novoTecnico
			service := NovoAtendimentoService(tt.geradorID, repo, chamadoRepo, categoriaPermissaoRepo, usuarioRepo)
			atendimentoCriado, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, atendimentoCriado)
			}
		})
	}
}

// Test_AtendimentoService_BuscarPorID testa o método BuscarPorID do serviço de atendimento.
func Test_AtendimentoService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() atd.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento)
	}{
		{
			nome: "buscar atendimento por id com sucesso",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				return repo
			},
			id:           atendimentoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoNulo(t, atendimento)
				verificarIDs(t, atendimentoTesteID, atendimento.ID())
				verificarIDs(t, tecnicoTesteID, atendimento.AtribuidoID())
				verificarIDs(t, chamadoTesteID, atendimento.ChamadoID())
			},
		},
		{
			nome: "erro ao buscar atendimento por id no repositório",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           atendimentoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
		{
			nome: "atendimento não encontrado por id",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoBuscar = mysql.ErrAtendimentoNaoEncontrado
				return repo
			},
			id:           atendimentoTesteID,
			erroEsperado: mysql.ErrAtendimentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoAtendimentoService(nil, repo, nil, nil, nil)
			atendimentoEncontrado, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, atendimentoEncontrado)
			}
		})
	}
}

// Test_AtendimentoService_BuscarPorChamadoEAtribuidoID testa o método BuscarPorChamadoEAtribuidoID do serviço de atendimento.
func Test_AtendimentoService_BuscarPorChamadoEAtribuidoID(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() atd.Repository
		chamadoID         string
		atribuidoID       string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento)
	}{
		{
			nome: "buscar atendimento por chamado e atribuído com sucesso",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				return repo
			},
			chamadoID:   chamadoTesteID,
			atribuidoID: tecnicoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoNulo(t, atendimento)
				verificarIDs(t, atendimentoTesteID, atendimento.ID())
				verificarIDs(t, tecnicoTesteID, atendimento.AtribuidoID())
				verificarIDs(t, chamadoTesteID, atendimento.ChamadoID())
			},
		},
		{
			nome: "erro ao buscar atendimento por chamado e atribuído no repositório",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			chamadoID:   chamadoTesteID,
			atribuidoID: tecnicoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
		{
			nome: "atendimento não encontrado por chamado e atribuído",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoBuscar = mysql.ErrAtendimentoNaoEncontrado
				return repo
			},
			chamadoID:   chamadoTesteID,
			atribuidoID: tecnicoTesteID,
			erroEsperado: mysql.ErrAtendimentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoAtendimentoService(nil, repo, nil, nil, nil)
			atendimentoEncontrado, erroRecebido := service.BuscarPorChamadoEAtribuidoID(ctx, tt.chamadoID, tt.atribuidoID)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, atendimentoEncontrado)
			}
		})
	}
}

// Test_AtendimentoService_Atualizar testa o método Atualizar do serviço de atendimento.
func Test_AtendimentoService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() atd.Repository
		id                 string
		params             atd.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento)
	}{
		{
			nome: "atualizar atendimento com sucesso",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				return repo
			},
			id:     atendimentoTesteID,
			params: novoAtualizarAtendimentoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				verificarNaoNulo(t, atendimento)
				verificarIDs(t, atendimentoTesteID, atendimento.ID())
				verificarIDs(t, tecnicoTesteID, atendimento.AtribuidoID())
				verificarIDs(t, chamadoTesteID, atendimento.ChamadoID())
				verificarDatas(t, atendimento.CriadoEm(), atendimento.AtualizadoEm())
			},
		},
		{
			nome: "erro ao atualizar atendimento no repositório",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:     atendimentoTesteID,
			params: novoAtualizarAtendimentoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
		{
			nome: "atendimento não encontrado para atualizar",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				novoAtendimento := novoAtendimentoTeste()
				repo.atendimentos[atendimentoTesteID] = novoAtendimento
				repo.erroAoAtualizar = mysql.ErrAtendimentoNaoEncontrado
				return repo
			},
			id:     atendimentoTesteID,
			params: novoAtualizarAtendimentoParamsTeste(),
			erroEsperado: mysql.ErrAtendimentoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimento *atd.Atendimento) {
				t.Helper()

				_, ok := repo.(*fakeAtendimentoRepository).atendimentos[atendimentoTesteID]
				verificarOk(t, ok, "deveria existir atendimento no repositório")
				verificarNulo(t, atendimento)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			chamadoRepo := novoFakeChamadoRepository()
			chamadoRepo.chamados[chamadoTesteID] = novoChamadoTeste()

			categoriaPermissaoRepo := novoFakeCategoriaPermissaoRepository()
			categoriaPermissaoRepo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()

			usuarioRepo := novoFakeUsuarioRepository()
			novoTecnico := novoTecnicoTeste()
			usuarioRepo.usuarios[tecnicoTesteID] = novoTecnico
			service := NovoAtendimentoService(nil, repo, chamadoRepo, categoriaPermissaoRepo, usuarioRepo)
			atendimentoAtualizado, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, atendimentoAtualizado)
			}
		})
	}
}

// Test_AtendimentoService_Listar testa o método Listar do serviço de atendimento.
func Test_AtendimentoService_Listar(t *testing.T) {
	t.Parallel()

	ctx := contextoTecnicoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() atd.Repository
		filtro             atd.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo atd.Repository, atendimentos []atd.Atendimento, total int, filtroRetornado atd.Filtro)
	}{
		{
			nome: "listar atendimentos paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				popuplarAtendimentosTeste(repo, 10)
				return repo
			},
			filtro: atd.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimentos []atd.Atendimento, total int, filtroRetornado atd.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeAtendimentoRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.atendimentos))
				verificarLimite(t, atendimentos, 5)
			},
		},
		{
			nome: "erro ao listar atendimentos no repositório",
			prepararRepo: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				popuplarAtendimentosTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro: atd.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo atd.Repository, atendimentos []atd.Atendimento, total int, filtroRetornado atd.Filtro) {
				t.Helper()
				
				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, atendimentos, 0)
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoAtendimentoService(nil, repo, nil, nil, nil)
			atendimentos, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, atendimentos, total, filtroRetornado)
			}
		})
	}
}
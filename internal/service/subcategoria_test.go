package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que fakeSubcategoriaRepository implementa subc.Repository
var _ subc.Repository = (*fakeSubcategoriaRepository)(nil)

// fakeSubcategoriaRepository é um repositório fake para testes, com estado interno.
type fakeSubcategoriaRepository struct {
	subcategorias   map[string]*subc.Subcategoria
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
}

// novoFakeSubcategoriaRepository cria uma nova instância do repositório fake de subcategorias.
func novoFakeSubcategoriaRepository() *fakeSubcategoriaRepository {
	return &fakeSubcategoriaRepository{
		subcategorias: make(map[string]*subc.Subcategoria),
	}
}

// Implementação dos métodos da interface sub.Repository
func (r *fakeSubcategoriaRepository) Criar(ctx context.Context, s subc.Subcategoria) (*subc.Subcategoria, error) {
	if r.erroAoCriar != nil {
		return nil, r.erroAoCriar
	}
	subcategoria := s
	r.subcategorias[s.ID()] = &subcategoria
	return &subcategoria, nil
}

// Atualizar atualiza uma subcategoria existente no repositório fake.
func (r *fakeSubcategoriaRepository) Atualizar(ctx context.Context, id string, s subc.Subcategoria) (*subc.Subcategoria, error) {
	if r.erroAoAtualizar != nil {
		return nil, r.erroAoAtualizar
	}
	subcategoria := s
	r.subcategorias[id] = &subcategoria
	return &subcategoria, nil
}

// BuscarPorID busca uma subcategoria por ID no repositório fake.
func (r *fakeSubcategoriaRepository) BuscarPorID(ctx context.Context, id string) (*subc.Subcategoria, error) {
	if r.erroAoBuscar != nil {
		return nil, r.erroAoBuscar
	}
	s, existe := r.subcategorias[id]
	if !existe {
		return nil, mysql.ErrSubcategoriaNaoEncontrada
	}
	return s, nil
}

// BuscarPorID busca uma subcategoria por ID no repositório fake.
func (r *fakeSubcategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*subc.Subcategoria, error) {
	if r.erroAoBuscar != nil {
		return nil, r.erroAoBuscar
	}
	for _, s := range r.subcategorias {
		if s.Nome() == nome {
			return s, nil
		}
	}
	return nil, mysql.ErrSubcategoriaNaoEncontrada
}

// Listar lista subcategorias no repositório fake com base no filtro fornecido.
func (r *fakeSubcategoriaRepository) Listar(ctx context.Context, filtro subc.Filtro) ([]subc.Subcategoria, int, error) {
	if r.erroAoListar != nil {
		return nil, 0, r.erroAoListar
	}

	// Converter o mapa para uma slice
	subcategorias := make([]subc.Subcategoria, 0, len(r.subcategorias))
	for _, s := range r.subcategorias {
		subcategorias = append(subcategorias, *s)
	}

	// Ordenação ORDER BY nome ASC
	sort.Slice(subcategorias, func(i, j int) bool {
		return subcategorias[i].Nome() < subcategorias[j].Nome()
	})

	// Contar o total antes da paginação
	total := len(subcategorias)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio >= total {
			return []subc.Subcategoria{}, total, nil
		}

		if fim > total {
			fim = total
		}

		subcategorias = subcategorias[inicio:fim]
	}

	return subcategorias, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const (
	subcategoriaTesteID   = "subcategoria-123"
	subcategoriaTesteNome = "Suporte" 
)

// novoSubcategoriaTeste cria uma subcategoria de teste com os valores fornecidos.
func novoSubcategoriaTeste() *subc.Subcategoria {
	s, _ := subc.Novo(
		subcategoriaTesteID,
		subcategoriaTesteNome,
		categoriaTesteID,
	)
	return s
}

// popularSubcategoriasTeste popula o repositório fake com subcategorias de teste.
func popularSubcategoriasTeste(repo *fakeSubcategoriaRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("subcategoria-%02d", i)

		s, _ := subc.Novo(
			id,
			fmt.Sprintf("Subcategoria %02d", i),
			categoriaTesteID,
		)
		repo.subcategorias[id] = s
	}
}

// novoCriarSubcategoriaParamsTeste cria parâmetros de criação de subcategoria para testes.
func novoCriarSubcategoriaParamsTeste() subc.CriarParams {
	return subc.CriarParams{
		Nome:        subcategoriaTesteNome,
		CategoriaID: categoriaTesteID,
	}
}

// novoCriarSubcategoriaParamsInvalidosTeste cria parâmetros inválidos de criação de subcategoria para testes.
func novoCriarSubcategoriaParamsInvalidosTeste() subc.CriarParams {
	return subc.CriarParams{
		Nome:        "",
		CategoriaID: "",
	}
}

// novoAtualizarSubcategoriaParamsTeste cria parâmetros de atualização de subcategoria para testes.
func novoAtualizarSubcategoriaParamsTeste() subc.AtualizarParams {
	return subc.AtualizarParams{
		Nome: ptr("Suporte Avançado"),
	}
}

// novoAtualizarSubcategoriaParamsInvalidosTeste cria parâmetros inválidos de atualização de subcategoria para testes.
func novoAtualizarSubcategoriaParamsInvalidosTeste() subc.AtualizarParams {
	return subc.AtualizarParams{
		Nome: ptr(""),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_SubcategoriaService_Criar testa o método Criar do serviço de subcategorias.
func Test_SubcategoriaService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() subc.Repository
		params             subc.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "criar subcategoria com sucesso",
			geradorID: &fakeGeradorID{
				id: subcategoriaTesteID,
			},
			prepararRepo: func() subc.Repository {
				return novoFakeSubcategoriaRepository()
			},
			params:       novoCriarSubcategoriaParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				if subcategoria.Nome() != novoCriarSubcategoriaParamsTeste().Nome {
					t.Errorf("Nome esperado '%s', recebido '%s'",
						novoCriarSubcategoriaParamsTeste().Nome, subcategoria.Nome())
				}

				if subcategoria.CategoriaID() != novoCriarSubcategoriaParamsTeste().CategoriaID {
					t.Errorf("CategoriaID esperado '%s', recebido '%s'",
						novoCriarSubcategoriaParamsTeste().CategoriaID, subcategoria.CategoriaID())
				}

				verificarIDs(t, subcategoriaTesteID, subcategoria.ID())
				verificarDatas(t, subcategoria.CriadoEm(), subcategoria.AtualizadoEm())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria criada, recebeu nil")
				verificarPersistido(t, len(repo.(*fakeSubcategoriaRepository).subcategorias))
			},
		},
		{
			nome: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() subc.Repository {
				return novoFakeSubcategoriaRepository()
			},
			params:       novoCriarSubcategoriaParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeSubcategoriaRepository).subcategorias))
				verificarNulo(t, subcategoria, "não deveria retornar subcategoria quando gerador de ID falha")
			},
		},
		{
			nome: "erro ao salvar no repositorio",
			geradorID: &fakeGeradorID{
				id: subcategoriaTesteID,
			},
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarSubcategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeSubcategoriaRepository).subcategorias))
				verificarNulo(t, subcategoria, "não deveria retornar subcategoria quando repositório falha")
			},
		},
		{
			nome: "erro de validação ao criar subcategoria",
			geradorID: &fakeGeradorID{
				id: subcategoriaTesteID,
			},
			prepararRepo: func() subc.Repository {
				return novoFakeSubcategoriaRepository()
			},
			params:       novoCriarSubcategoriaParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeSubcategoriaRepository).subcategorias))
				verificarNulo(t, subcategoria, "não deveria retornar subcategoria quando parâmetros são inválidos")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			servico := NovoSubcategoriaService(tt.geradorID, repo)
			subcategoriaCriada, erroRecebido := servico.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategoriaCriada)
			}
		})
	}
}

// Test_SubcategoriaService_Atualizar testa o método Atualizar do serviço de subcategorias.
func Test_SubcategoriaService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		id                 string
		params             subc.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "atualizar subcategoria com sucesso",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			id:           subcategoriaTesteID,
			params:       novoAtualizarSubcategoriaParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				if novoAtualizarSubcategoriaParamsTeste().Nome != nil {
					if subcategoria.Nome() != *novoAtualizarSubcategoriaParamsTeste().Nome {
						t.Errorf("Nome esperado '%s', recebido '%s'",
							*novoAtualizarSubcategoriaParamsTeste().Nome, subcategoria.Nome())
					}
				}

				if novoAtualizarSubcategoriaParamsTeste().CategoriaID != nil {
					if subcategoria.CategoriaID() != *novoAtualizarSubcategoriaParamsTeste().CategoriaID {
						t.Errorf("CategoriaID esperado '%s', recebido '%s'",
							*novoAtualizarSubcategoriaParamsTeste().CategoriaID, subcategoria.CategoriaID())
					}
				}

				verificarIDs(t, subcategoriaTesteID, subcategoria.ID())
				verificarDatas(t, subcategoria.CriadoEm(), subcategoria.AtualizadoEm())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria atualizada, recebeu nil")
			},
		},
		{
			nome: "erro ao buscar subcategoria no repositório para atualizar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           subcategoriaTesteID,
			params:       novoAtualizarSubcategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao buscar")
			},
		},
		{
			nome: "erro ao persistir subcategoria atualizada no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           subcategoriaTesteID,
			params:       novoAtualizarSubcategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao atualizar")
			},
		},
		{
			nome: "subcategoria não encontrada para atualizar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = mysql.ErrSubcategoriaNaoEncontrada
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarSubcategoriaParamsTeste(),
			erroEsperado: mysql.ErrSubcategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando subcategoria não é encontrada")
			},
		},
		{
			nome: "erro de validação ao atualizar subcategoria",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			id:           subcategoriaTesteID,
			params:       novoAtualizarSubcategoriaParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro de validação")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)
			subcategoriaAtualizada, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategoriaAtualizada)
			}
		})
	}
}

// Test_SubcategoriaService_BuscarPorID testa o método BuscarPorID do serviço de subcategorias.
func Test_SubcategoriaService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "buscar subcategoria por id com sucesso",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			id:           subcategoriaTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				verificarIDs(t, subcategoriaTesteID, subcategoria.ID())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria encontrada, recebeu nil")
			},
		},
		{
			nome: "erro ao buscar subcategoria por id no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           subcategoriaTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				_, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao buscar")
			},
		},
		{
			nome: "subcategoria não encontrada por id",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = mysql.ErrSubcategoriaNaoEncontrada
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrSubcategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				_, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando subcategoria não é encontrada")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)
			s, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, s)
			}
		})
	}
}

// Test_SubcategoriaService_BuscarPorNome testa o método BuscarPorNome do serviço de subcategorias.
func Test_SubcategoriaService_BuscarPorNome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		nomeSubcategoria   string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "buscar subcategoria por nome com sucesso",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			nomeSubcategoria: subcategoriaTesteNome,
			erroEsperado:     nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				verificarNomes(t, subcategoriaTesteNome, subcategoria.Nome())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria encontrada, recebeu nil")
			},
		},
		{
			nome: "erro ao buscar subcategoria por nome no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			nomeSubcategoria: subcategoriaTesteNome,
			erroEsperado:     errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				_, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao buscar")
			},
		},
		{
			nome: "subcategoria não encontrada por nome",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = mysql.ErrSubcategoriaNaoEncontrada
				return repo
			},
			nomeSubcategoria: "nome-inexistente",
			erroEsperado:     mysql.ErrSubcategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				_, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando subcategoria não é encontrada")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)

			subcategoriaEncontrada, erroRecebido := service.BuscarPorNome(ctx, tt.nomeSubcategoria)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategoriaEncontrada)
			}
		})
	}
}

// Test_SubcategoriaService_Desativar testa o método Desativar do serviço de subcategorias.
func Test_SubcategoriaService_Desativar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "desativar subcategoria com sucesso",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Ativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			id: 				 subcategoriaTesteID,
			erroEsperado:       nil,
			verificarResultado: func (t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				if subcategoria.Status() {
					t.Errorf("Status esperado desativado, recebeu Status=%v", subcategoria.Status())
				}

				verificarIDs(t, subcategoriaTesteID, subcategoria.ID())
				verificarDatas(t, subcategoria.CriadoEm(), subcategoria.AtualizadoEm())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria desativada, recebeu nil")
			},
		},
		{
			nome: "erro ao buscar subcategoria no repositório para desativar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Ativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id: 				 subcategoriaTesteID,
			erroEsperado:       errFakeRepo,
			verificarResultado: func (t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao buscar")
			},
		},
		{
			nome: "erro ao persistir subcategoria desativada no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Ativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id: 				 subcategoriaTesteID,
			erroEsperado:       errFakeRepo,
			verificarResultado: func (t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao atualizar")
			},
		},
		{
			nome: "subcategoria não encontrada para desativar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Ativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = mysql.ErrSubcategoriaNaoEncontrada
				return repo
			},
			id: 				 "id-inexistente",
			erroEsperado:       mysql.ErrSubcategoriaNaoEncontrada,
			verificarResultado: func (t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando subcategoria não é encontrada")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)
			subcategoriaDesativada, erroRecebido := service.Desativar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategoriaDesativada)
			}
		})
	}
}

// Test_SubcategoriaService_Ativar testa o método Ativar do serviço de subcategorias.
func Test_SubcategoriaService_Ativar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	
	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria)
	}{
		{
			nome: "ativar subcategoria com sucesso",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Desativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				return repo
			},
			id:           subcategoriaTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				if !subcategoria.Status() {
					t.Errorf("Status esperado ativado, recebeu Status=%v", subcategoria.Status())
				}

				verificarIDs(t, subcategoriaTesteID, subcategoria.ID())
				verificarDatas(t, subcategoria.CriadoEm(), subcategoria.AtualizadoEm())
				verificarNaoNulo(t, subcategoria, "esperava subcategoria ativada, recebeu nil")
			},
		},
		{
			nome: "erro ao buscar subcategoria no repositório para ativar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Desativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           subcategoriaTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao buscar")
			},
		},
		{
			nome: "erro ao persistir subcategoria ativada no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Desativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           subcategoriaTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando há erro ao atualizar")
			},
		},
		{
			nome: "subcategoria não encontrada para ativar",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				novaSubcategoria := novoSubcategoriaTeste()
				novaSubcategoria.Desativar()
				repo.subcategorias[novaSubcategoria.ID()] = novaSubcategoria
				repo.erroAoBuscar = mysql.ErrSubcategoriaNaoEncontrada
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrSubcategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategoria *subc.Subcategoria) {
				t.Helper()

				subcategoriaRepo, ok := repo.(*fakeSubcategoriaRepository).subcategorias[subcategoriaTesteID]
				verificarOk(t, ok, "deveria existir subcategoria no repositório")
				verificarNaoAlterado(t, novoSubcategoriaTeste(), subcategoriaRepo, compararSubcategorias)
				verificarDatas(t, subcategoriaRepo.CriadoEm(), subcategoriaRepo.AtualizadoEm())
				verificarNulo(t, subcategoria, "não esperava subcategoria retornada quando subcategoria não é encontrada")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)
			subcategoriaAtivada, erroRecebido := service.Ativar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategoriaAtivada)
			}
		})
	}
}

// Test_SubcategoriaService_Listar testa o método Listar do serviço de subcategorias.
func Test_SubcategoriaService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() subc.Repository
		filtro             subc.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo subc.Repository, subcategorias []subc.Subcategoria, total int, filtroRetornado subc.Filtro)
	}{
		{
			nome: "listar 10 subcategorias com sucesso - pagina 1 limite 5",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				popularSubcategoriasTeste(repo, 10)
				return repo
			},
			filtro:       subc.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategorias []subc.Subcategoria, total int, filtroRetornado subc.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeSubcategoriaRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.subcategorias))
				verificarLimite(t, subcategorias, 5)
			},
		},
		{
			nome: "erro ao listar subcategorias paginadas no repositório",
			prepararRepo: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository()
				popularSubcategoriasTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       subc.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo subc.Repository, subcategorias []subc.Subcategoria, total int, filtroRetornado subc.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, subcategorias, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoSubcategoriaService(nil, repo)

			subcategorias, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, subcategorias, total, filtroRetornado)
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

// compararSubcategorias verifica se duas subcategorias são iguais comparando seus campos.
func compararSubcategorias(t *testing.T, antes, depois *subc.Subcategoria) {
	t.Helper()

	if antes == nil && depois == nil {
		t.Fatalf("subcategoria não deveria nil: antes=%v depois=%v", antes, depois)
	}

	if antes.ID() != depois.ID() ||
		antes.Nome() != depois.Nome() ||
		antes.CategoriaID() != depois.CategoriaID() {

		t.Fatalf("subcategorias não deveriam ser diferentes:\nantes: %+v\ndepois: %+v",
			antes, depois,
		)
	}
}

package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// fakeCategoriaRepository é uma implementação fake do repositório de categorias para testes.
type fakeCategoriaRepository struct {
	categorias      map[string]*ctg.Categoria
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
}

// novoFakeCategoriaRepository cria uma nova instância do repositório fake.
func novoFakeCategoriaRepository() *fakeCategoriaRepository {
	return &fakeCategoriaRepository{
		categorias: make(map[string]*ctg.Categoria),
	}
}

// Criar adiciona uma nova categoria ao repositório fake.
func (f *fakeCategoriaRepository) Criar(ctx context.Context, c ctg.Categoria) (*ctg.Categoria, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	categoria := c
	f.categorias[c.ID()] = &categoria
	return &categoria, nil
}

// Atualizar atualiza uma categoria existente no repositório fake.
func (f *fakeCategoriaRepository) Atualizar(ctx context.Context, id string, c ctg.Categoria) (*ctg.Categoria, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	categoria := c
	f.categorias[id] = &categoria
	return &categoria, nil
}

// BuscarPorID recupera uma categoria pelo ID do repositório fake.
func (f *fakeCategoriaRepository) BuscarPorID(ctx context.Context, id string) (*ctg.Categoria, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	c, existe := f.categorias[id]
	if !existe {
		return nil, mysql.ErrCategoriaNaoEncontrada
	}

	return c, nil
}

// BuscarPorNome recupera uma categoria pelo nome do repositório fake.
func (f *fakeCategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*ctg.Categoria, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}

	c, existe := f.categorias[nome]
	if !existe {
		return nil, mysql.ErrCategoriaNaoEncontrada
	}

	return c, nil
}

// Listar lista categorias do repositório fake com paginação.
func (f *fakeCategoriaRepository) Listar(ctx context.Context, filtro ctg.Filtro) ([]ctg.Categoria, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para uma slice
	categorias := make([]ctg.Categoria, 0, len(f.categorias))
	for _, c := range f.categorias {
		categorias = append(categorias, *c)
	}

	// Ordenação ORDER BY nome ASC
	sort.Slice(categorias, func(i, j int) bool {
		return categorias[i].Nome() < categorias[j].Nome()
	})

	// Contar o total antes da paginação
	total := len(categorias)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio > total {
			return []ctg.Categoria{}, total, nil
		}

		if fim > total {
			fim = total
		}

		categorias = categorias[inicio:fim]
	}

	return categorias, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const (
	categoriaTesteID   = "categoria-123"
	categoriaTesteNome = "Categoria de Teste"
)

// novoCategoriaTeste cria uma nova categoria para testes.
func novoCategoriaTeste() *ctg.Categoria {
	c, _ := ctg.Novo(categoriaTesteID, categoriaTesteNome)
	return c
}

// popularCategoriasTeste adiciona múltiplas categorias ao repositório falso para testes.
func popularCategoriasTeste(repo *fakeCategoriaRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("categoria-%02d", i)

		c, _ := ctg.Novo(id, fmt.Sprintf("Categoria %02d", i))
		repo.categorias[id] = c
	}
}

// novoCriarCategoriaParamsTeste cria parâmetros de criação de categoria para testes.
func novoCriarCategoriaParamsTeste() ctg.CriarParams {
	return ctg.CriarParams{
		Nome: categoriaTesteNome,
	}
}

// novoCriarCategoriaParamsInvalidosTeste cria parâmetros inválidos de criação de categoria para testes.
func novoCriarCategoriaParamsInvalidosTeste() ctg.CriarParams {
	return ctg.CriarParams{
		Nome: "",
	}
}

// novoAtualizarCategoriaParamsTeste cria parâmetros de atualização de categoria para testes.
func novoAtualizarCategoriaParamsTeste() ctg.AtualizarParams {
	return ctg.AtualizarParams{
		Nome:   ptr("Categoria Atualizada"),
		Status: ptr(true),
	}
}

// novoAtualizarCategoriaParamsInvalidosTeste cria parâmetros inválidos de atualização de categoria para testes.
func novoAtualizarCategoriaParamsInvalidosTeste() ctg.AtualizarParams {
	return ctg.AtualizarParams{
		Nome:   ptr(""),
		Status: ptr(true),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_CategoriaService_Criar testa o método Criar do CategoriaService.
func Test_CategoriaService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() ctg.Repository
		params             ctg.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			nome: "criar categoria com sucesso",
			geradorID: &fakeGeradorID{
				id: categoriaTesteID,
			},
			prepararRepo: func() ctg.Repository {
				return novoFakeCategoriaRepository()
			},
			params:       novoCriarCategoriaParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoNulo(t, categoria)
				verificarPersistido(t, len(repo.(*fakeCategoriaRepository).categorias))

				if categoria.Nome() != novoCriarCategoriaParamsTeste().Nome {
					t.Errorf("Nome esperado '%s', recebeu '%s'", categoriaTesteNome, categoria.Nome())
				}

				verificarIDs(t, categoriaTesteID, categoria.ID())
				verificarDatas(t, categoria.CriadoEm(), categoria.AtualizadoEm())
			},
		},
		{
			nome: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() ctg.Repository {
				return novoFakeCategoriaRepository()
			},
			params:       novoCriarCategoriaParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeCategoriaRepository).categorias))
				verificarNulo(t, categoria)
			},
		},
		{
			nome: "erro ao salvar categoria no repositório",
			geradorID: &fakeGeradorID{
				id: categoriaTesteID,
			},
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarCategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeCategoriaRepository).categorias))
				verificarNulo(t, categoria)
			},
		},
		{
			nome: "erro de validação ao criar categoria",
			geradorID: &fakeGeradorID{
				id: categoriaTesteID,
			},
			prepararRepo: func() ctg.Repository {
				return novoFakeCategoriaRepository()
			},
			params:       novoCriarCategoriaParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeCategoriaRepository).categorias))
				verificarNulo(t, categoria)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaService(tt.geradorID, repo)
			categoriaCriada, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaCriada)
			}
		})
	}
}

// Test_CategoriaService_Atualizar testa o método Atualizar do CategoriaService.
func Test_CategoriaService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() ctg.Repository
		id                 string
		params             ctg.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			nome: "atualizar categoria com sucesso",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				return repo
			},
			id:           categoriaTesteID,
			params:       novoAtualizarCategoriaParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoNulo(t, categoria)

				if novoAtualizarCategoriaParamsTeste().Nome != nil {
					if categoria.Nome() != *novoAtualizarCategoriaParamsTeste().Nome {
						t.Errorf("Nome esperado '%s', recebeu '%s'", *novoAtualizarCategoriaParamsTeste().Nome, categoria.Nome())
					}
				}

				if novoAtualizarCategoriaParamsTeste().Status != nil {
					if categoria.Status() != *novoAtualizarCategoriaParamsTeste().Status {
						t.Errorf("Status esperado '%t', recebeu '%t'", *novoAtualizarCategoriaParamsTeste().Status, categoria.Status())
					}
				}

				verificarIDs(t, categoriaTesteID, categoria.ID())
				verificarDatas(t, categoria.CriadoEm(), categoria.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar categoria no repositório para atualizar",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           categoriaTesteID,
			params:       novoAtualizarCategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				categoriaRepo, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria no repositório")
				verificarNulo(t, categoria)
				verificarNaoAlterado(t, novoCategoriaTeste(), categoriaRepo, compararCategorias)
			},
		},
		{
			nome: "erro ao persistir categoria atualizada no repositório",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           categoriaTesteID,
			params:       novoAtualizarCategoriaParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				categoriaRepo, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria no repositório")
				verificarNulo(t, categoria)
				verificarNaoAlterado(t, novoCategoriaTeste(), categoriaRepo, compararCategorias)
				verificarDatas(t, categoriaRepo.CriadoEm(), categoriaRepo.AtualizadoEm())
			},
		},
		{
			nome: "categoria não encontrada para atualizar",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				repo.erroAoBuscar = mysql.ErrCategoriaNaoEncontrada
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarCategoriaParamsTeste(),
			erroEsperado: mysql.ErrCategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				categoriaRepo, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria original no repositório")
				verificarNulo(t, categoria)
				verificarNaoAlterado(t, novoCategoriaTeste(), categoriaRepo, compararCategorias)
				verificarDatas(t, categoriaRepo.CriadoEm(), categoriaRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro de validação ao atualizar categoria",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				return repo
			},
			id:           "categoria-123",
			params:       novoAtualizarCategoriaParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				categoriaRepo, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria original no repositório")
				verificarNulo(t, categoria)
				verificarNaoAlterado(t, novoCategoriaTeste(), categoriaRepo, compararCategorias)
				verificarDatas(t, categoriaRepo.CriadoEm(), categoriaRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaService(nil, repo)
			categoriaAtualizada, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaAtualizada)
			}
		})
	}
}

// Test_CategoriaService_BuscarPorID testa o método BuscarPorID do CategoriaService.
func Test_CategoriaService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() ctg.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			nome: "buscar categoria por id com sucesso",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				return repo
			},
			id:           categoriaTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoNulo(t, categoria)
				verificarIDs(t, categoriaTesteID, categoria.ID())
			},
		},
		{
			nome: "erro ao buscar categoria por id no repositório",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           categoriaTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria no repositório")
				verificarNulo(t, categoria)
			},
		},
		{
			nome: "categoria não encontrada por id",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteID] = novaCategoria
				repo.erroAoBuscar = mysql.ErrCategoriaNaoEncontrada
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrCategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteID]
				verificarOk(t, ok, "deveria existir categoria original no repositório")
				verificarNulo(t, categoria)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaService(nil, repo)
			categoriaEncontrada, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaEncontrada)
			}
		})
	}
}

// Test_CategoriaService_BuscarPorNome testa o método BuscarPorNome do CategoriaService.
func Test_CategoriaService_BuscarPorNome(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() ctg.Repository
		nomeCategoria      string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			nome: "buscar categoria por nome com sucesso",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteNome] = novaCategoria
				return repo
			},
			nomeCategoria: categoriaTesteNome,
			erroEsperado:  nil,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				verificarNaoNulo(t, categoria)
				verificarParamBusca(t, categoriaTesteNome, categoria.Nome())
			},
		},
		{
			nome: "erro ao buscar categoria por nome no repositório",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteNome] = novaCategoria
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			nomeCategoria: categoriaTesteNome,
			erroEsperado:  errFakeRepo,

			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteNome]
				verificarOk(t, ok, "deveria existir categoria no repositório")
				verificarNulo(t, categoria)
			},
		},
		{
			nome: "categoria não encontrada por nome",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				novaCategoria := novoCategoriaTeste()
				repo.categorias[categoriaTesteNome] = novaCategoria
				repo.erroAoBuscar = mysql.ErrCategoriaNaoEncontrada
				return repo
			},
			nomeCategoria: "nome-inexistente",
			erroEsperado:  mysql.ErrCategoriaNaoEncontrada,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaRepository).categorias[categoriaTesteNome]
				verificarOk(t, ok, "deveria existir categoria original no repositório")
				verificarNulo(t, categoria)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaService(nil, repo)
			categoriaEncontrada, erroRecebido := service.BuscarPorNome(ctx, tt.nomeCategoria)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaEncontrada)
			}
		})
	}
}

// Test_CategoriaService_Listar testa o método Listar do CategoriaService.
func Test_CategoriaService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() ctg.Repository
		filtro             ctg.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo ctg.Repository, categorias []ctg.Categoria, total int, filtroRetornado ctg.Filtro)
	}{
		{
			nome: "listar 10 categorias paginadas com sucesso - pagina 1 limite 5",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				popularCategoriasTeste(repo, 10)
				return repo
			},
			filtro:       ctg.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categorias []ctg.Categoria, total int, filtroRetornado ctg.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeCategoriaRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.categorias))
				verificarLimite(t, categorias, 5)
			},
		},
		{
			nome: "erro ao listar categorias paginadas no repositório",
			prepararRepo: func() ctg.Repository {
				repo := novoFakeCategoriaRepository()
				popularCategoriasTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       ctg.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo ctg.Repository, categorias []ctg.Categoria, total int, filtroRetornado ctg.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, categorias, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaService(nil, repo)
			categorias, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categorias, total, filtroRetornado)
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

// compararCategorias verifica se duas categorias são iguais comparando seus campos.
func compararCategorias(t *testing.T, antes, depois *ctg.Categoria) {
	t.Helper()

	if antes == nil || depois == nil {
		t.Fatalf("categorias para comparação não podem ser nil: antes=%v depois=%v", antes, depois)
	}

	if antes.ID() != depois.ID() ||
		antes.Nome() != depois.Nome() ||
		antes.Status() != depois.Status() {

		t.Fatalf(
			"categorias não deveriam ser diferentes: \nantes=%+v\ndepois=%+v",
			antes, depois,
		)
	}
}

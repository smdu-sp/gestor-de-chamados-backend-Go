package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que fakeCategoriaPermissaoRepository implementa cpm.Repository
var _ cpm.Repository = (*fakeCategoriaPermissaoRepository)(nil)

// fakeCategoriaPermissaoRepository é um repositório falso para testes de categoria de permissão.
type fakeCategoriaPermissaoRepository struct {
	categoriasPermissoes map[string]*cpm.CategoriaPermissao
	erroAoCriar          error
	erroAoAtualizar      error
	erroAoBuscar         error
	erroAoDeletar        error
	erroAoListar         error
}

// novoFakeCategoriaPermissaoRepository cria uma nova instância do repositório fake.
func novoFakeCategoriaPermissaoRepository() *fakeCategoriaPermissaoRepository {
	return &fakeCategoriaPermissaoRepository{
		categoriasPermissoes: make(map[string]*cpm.CategoriaPermissao),
	}
}

// Criar adiciona uma nova categoria ao repositório fake.
func (f *fakeCategoriaPermissaoRepository) Criar(ctx context.Context, c cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	categoriaPermissao := c
	f.categoriasPermissoes[chaveCompostaCategoriaPermissao] = &categoriaPermissao
	return &categoriaPermissao, nil
}

// BuscarPorID busca uma categoria pelo ID composto no repositório fake.
func (f *fakeCategoriaPermissaoRepository) BuscarPorID(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	c, existe := f.categoriasPermissoes[chaveCompostaCategoriaPermissao]
	if !existe {
		return nil, mysql.ErrCategoriaPermissaoNaoEncontrada
	}
	return c, nil
}

// Atualizar atualiza uma categoria existente no repositório fake.
func (f *fakeCategoriaPermissaoRepository) Atualizar(ctx context.Context, categoriaID, usuarioID string, c cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	categoriaPermissao := c
	f.categoriasPermissoes[chaveCompostaCategoriaPermissao] = &categoriaPermissao
	return &categoriaPermissao, nil
}

// Listar retorna uma lista de categorias do repositório fake.
func (f *fakeCategoriaPermissaoRepository) Listar(ctx context.Context, filtro cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para uma slice
	categoriasPermissoes := make([]cpm.CategoriaPermissao, 0, len(f.categoriasPermissoes))
	for _, cpm := range f.categoriasPermissoes {
		categoriasPermissoes = append(categoriasPermissoes, *cpm)
	}

	// Ordernar ORDER BY criado_em DESC
	sort.Slice(categoriasPermissoes, func(i, j int) bool {
		return categoriasPermissoes[i].CriadoEm().After(categoriasPermissoes[j].CriadoEm())
	})

	// Contar total antes da paginação
	total := len(categoriasPermissoes)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio > total {
			return []cpm.CategoriaPermissao{}, total, nil
		}

		if fim > total {
			fim = total
		}

		categoriasPermissoes = categoriasPermissoes[inicio:fim]
	}

	return categoriasPermissoes, total, nil
}

// Deletar remove uma categoria do repositório fake.
func (f *fakeCategoriaPermissaoRepository) Deletar(ctx context.Context, categoriaID, usuarioID string) error {
	if f.erroAoDeletar != nil {
		return f.erroAoDeletar
	}
	delete(f.categoriasPermissoes, chaveCompostaCategoriaPermissao)
	return nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const chaveCompostaCategoriaPermissao = categoriaTesteID + ":" + usuarioTesteID

// novoCategoriaPermissaoTest cria uma nova categoria de permissão para testes.
func novoCategoriaPermissaoTeste() *cpm.CategoriaPermissao {
	categoriaPermissao, _ := cpm.Novo(categoriaTesteID, usuarioTesteID, usr.PermADM)
	return categoriaPermissao
}

// popularCategoriaPermissaoTeste adiciona múltiplas categorias de permissão ao repositório para testes.
func popularCategoriaPermissaoTeste(repo *fakeCategoriaPermissaoRepository, quantidade int) {
	// Gerar categorias/permissões com IDs únicos para evitar colisões
	// levando em conta que o id é composto por categoriaID e usuarioID, podemos variar ambos para garantir unicidade
	for i := 0; i <= quantidade; i++ {
		categoriaID := fmt.Sprintf("categoria-%02d", i)
		usuarioID := fmt.Sprintf("usuario-%02d", i)
		chaveComposta := categoriaID + ":" + usuarioID
		categoriaPermissao, _ := cpm.Novo(categoriaID, usuarioID, usr.PermADM)
		repo.categoriasPermissoes[chaveComposta] = categoriaPermissao
	}
}

// novoCriarCategoriaPermissaoParamsTeste cria parâmetros para criar uma categoria de permissão para testes.
func novoCriarCategoriaPermissaoParamsTeste() cpm.CriarParams {
	return cpm.CriarParams{
		CategoriaID: categoriaTesteID,
		UsuarioID:   usuarioTesteID,
		Permissao:   usr.PermADM,
	}
}

// novoCriarCategoriaPermissaoParamsInvalidosTeste cria parâmetros inválidos para criar uma categoria de permissão para testes.
func novoCriarCategoriaPermissaoParamsInvalidosTeste() cpm.CriarParams {
	return cpm.CriarParams{
		CategoriaID: "",
		UsuarioID:   "",
		Permissao:   "",
	}
}

// novoAtualizarCategoriaPermissaoParamsTeste cria parâmetros para atualizar uma categoria de permissão para testes.
func novoAtualizarCategoriaPermissaoParamsTeste() cpm.AtualizarParams {
	return cpm.AtualizarParams{
		Permissao: usr.PermADM,
	}
}

// novoAtualizarCategoriaPermissaoParamsInvalidosTeste cria parâmetros inválidos para atualizar uma categoria de permissão para testes.
func novoAtualizarCategoriaPermissaoParamsInvalidosTeste() cpm.AtualizarParams {
	return cpm.AtualizarParams{
		Permissao: "",
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_CategoriaPermissaoService_Criar testa o método Criar do serviço de categoria de permissão.
func Test_CategoriaPermissaoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() cpm.Repository
		params             cpm.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao)
	}{
		{
			nome: "criar categoriaPermissão com sucesso",
			prepararRepo: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			params:       novoCriarCategoriaPermissaoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				verificarNaoNulo(t, categoriaPermissao)
				verificarPersistido(t, len(repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes))

				if categoriaPermissao.CategoriaID() != novoCriarCategoriaPermissaoParamsTeste().CategoriaID {
					t.Errorf("CategoriaID esperado '%s', recebido '%s'",
						novoCriarCategoriaPermissaoParamsTeste().CategoriaID, categoriaPermissao.CategoriaID())
				}

				if categoriaPermissao.UsuarioID() != novoCriarCategoriaPermissaoParamsTeste().UsuarioID {
					t.Errorf("UsuarioID esperado '%s', recebido '%s'",
						novoCriarCategoriaPermissaoParamsTeste().UsuarioID, categoriaPermissao.UsuarioID())
				}

				if categoriaPermissao.Permissao() != novoCriarCategoriaPermissaoParamsTeste().Permissao.String() {
					t.Errorf("Permissao esperada '%s', recebida '%s'",
						novoCriarCategoriaPermissaoParamsTeste().Permissao, categoriaPermissao.Permissao())
				}

				verificarDatas(t, categoriaPermissao.CriadoEm(), categoriaPermissao.AtualizadoEm())
			},
		},
		{
			nome: "erro ao salvar categoriaPermissão no repositório",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarCategoriaPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes))
				verificarNulo(t, categoriaPermissao)
			},
		},
		{
			nome: "erro ao criar categoriaPermissão com parâmetros inválidos",
			prepararRepo: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			params:       novoCriarCategoriaPermissaoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes))
				verificarNulo(t, categoriaPermissao)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaPermissaoService(repo)
			categoriaPermissaoCriada, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaPermissaoCriada)
			}
		})
	}
}

// Test_CategoriaPermissaoService_BuscarPorID testa o método BuscarPorID do serviço de categoria de permissão.
func Test_CategoriaPermissaoService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		categoriaID        string
		usuarioID          string
		prepararRepo       func() cpm.Repository
		erroEsperado       error
		verificarResultado func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao)
	}{
		{
			nome:        "buscar categoriaPermissão com sucesso",
			categoriaID: categoriaTesteID,
			usuarioID:   usuarioTesteID,
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				return repo
			},
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				verificarNaoNulo(t, categoriaPermissao)
				verificarIDs(t, categoriaTesteID, categoriaPermissao.CategoriaID())
				verificarIDs(t, usuarioTesteID, categoriaPermissao.UsuarioID())
			},
		},
		{
			nome:        "erro ao buscar categoriaPermissão por ids no repositório",
			categoriaID: "id-inexistente",
			usuarioID:   "id-inexistente",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
			},
		},
		{
			nome:        "categoriaPermissão não encontrada por ids",
			categoriaID: categoriaTesteID,
			usuarioID:   usuarioTesteID,
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				repo.erroAoBuscar = mysql.ErrCategoriaPermissaoNaoEncontrada
				return repo
			},
			erroEsperado: mysql.ErrCategoriaPermissaoNaoEncontrada,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaPermissaoService(repo)
			categoriaPermissaoEncontrada, erroRecebido := service.BuscarPorID(ctx, tt.categoriaID, tt.usuarioID)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaPermissaoEncontrada)
			}
		})
	}
}

// TestCategoriaPermissaoService_Atualizar testa o método Atualizar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() cpm.Repository
		categoriaID        string
		usuarioID          string
		params             cpm.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao)
	}{
		{
			nome: "atualizar categoriaPermissão com sucesso",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				return repo
			},
			categoriaID:  categoriaTesteID,
			usuarioID:    usuarioTesteID,
			params:       novoAtualizarCategoriaPermissaoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				verificarNaoNulo(t, categoriaPermissao)

				if categoriaPermissao.Permissao() != novoAtualizarCategoriaPermissaoParamsTeste().Permissao.String() {
					t.Errorf("Permissao esperada '%s', recebida '%s'",
						novoAtualizarCategoriaPermissaoParamsTeste().Permissao, categoriaPermissao.Permissao())
				}

				verificarIDs(t, categoriaTesteID, categoriaPermissao.CategoriaID())
				verificarIDs(t, usuarioTesteID, categoriaPermissao.UsuarioID())
				verificarDatas(t, categoriaPermissao.CriadoEm(), categoriaPermissao.AtualizadoEm())

			},
		},
		{
			nome: "erro ao buscar categoriaPermissão para atualizar",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissa := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissa
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			categoriaID:  categoriaTesteID,
			usuarioID:    usuarioTesteID,
			params:       novoAtualizarCategoriaPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				categoriaPermissaoRepo, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
				verificarNaoAlterado(t, novoCategoriaPermissaoTeste(), categoriaPermissaoRepo, compararCategoriasPermissao)
				verificarDatas(t, categoriaPermissaoRepo.CriadoEm(), categoriaPermissaoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir categoriaPermissão atualizada no repositório",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			categoriaID:  categoriaTesteID,
			usuarioID:    usuarioTesteID,
			params:       novoAtualizarCategoriaPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				categoriaPermissaoRepo, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
				verificarNaoAlterado(t, novoCategoriaPermissaoTeste(), categoriaPermissaoRepo, compararCategoriasPermissao)
				verificarDatas(t, categoriaPermissaoRepo.CriadoEm(), categoriaPermissaoRepo.AtualizadoEm())
			},
		},
		{
			nome: "categoriaPermissão não encontrada para atualizar",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				repo.erroAoBuscar = mysql.ErrCategoriaPermissaoNaoEncontrada
				return repo
			},
			categoriaID:  "id-inexistente",
			usuarioID:    "id-inexistente",
			params:       novoAtualizarCategoriaPermissaoParamsTeste(),
			erroEsperado: mysql.ErrCategoriaPermissaoNaoEncontrada,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				categoriaPermissaoRepo, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
				verificarNaoAlterado(t, novoCategoriaPermissaoTeste(), categoriaPermissaoRepo, compararCategoriasPermissao)
				verificarDatas(t, categoriaPermissaoRepo.CriadoEm(), categoriaPermissaoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro de validação ao atualizar categoriaPermissão",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				return repo
			},
			categoriaID:  categoriaTesteID,
			usuarioID:    usuarioTesteID,
			params:       novoAtualizarCategoriaPermissaoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriaPermissao *cpm.CategoriaPermissao) {
				t.Helper()

				categoriaPermissaoRepo, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
				verificarNulo(t, categoriaPermissao)
				verificarNaoAlterado(t, novoCategoriaPermissaoTeste(), categoriaPermissaoRepo, compararCategoriasPermissao)
				verificarDatas(t, categoriaPermissaoRepo.CriadoEm(), categoriaPermissaoRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaPermissaoService(repo)
			categoriaPermissaoAtualizada, erroRecebido := service.Atualizar(ctx, categoriaTesteID, usuarioTesteID, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriaPermissaoAtualizada)
			}
		})
	}
}

// TestCategoriaPermissaoService_Deletar testa o método Deletar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Deletar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() cpm.Repository
		categoriaID        string
		usuarioID          string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo cpm.Repository)
	}{
		{
			nome:        "deletar categoriaPermissão com sucesso",
			categoriaID: categoriaTesteID,
			usuarioID:   usuarioTesteID,
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				return repo
			},
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo cpm.Repository) {
				t.Helper()

				_, existe := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				if existe {
					t.Errorf("não deveria existir categoriaPermissão no repositório após deleção")
				}
			},
		},
		{
			nome:        "categoriaPermissão não encontrada para deletar",
			categoriaID: "id-inexistente",
			usuarioID:   "id-inexistente",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				novaCategoriaPermissao := novoCategoriaPermissaoTeste()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novaCategoriaPermissao
				repo.erroAoBuscar = mysql.ErrCategoriaPermissaoNaoEncontrada
				return repo
			},
			erroEsperado: mysql.ErrCategoriaPermissaoNaoEncontrada,
			verificarResultado: func(t *testing.T, repo cpm.Repository) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
			},
		},
		{
			nome:        "erro ao deletar categoriaPermissão do repositório",
			categoriaID: categoriaTesteID,
			usuarioID:   usuarioTesteID,
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()
				repo.erroAoDeletar = errFakeRepo
				return repo
			},
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository) {
				t.Helper()

				_, ok := repo.(*fakeCategoriaPermissaoRepository).categoriasPermissoes[chaveCompostaCategoriaPermissao]
				verificarOk(t, ok, "deveria existir categoriaPermissão no repositório")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaPermissaoService(repo)
			erroRecebido := service.Deletar(ctx, tt.categoriaID, tt.usuarioID)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo)
			}
		})
	}
}

// Test_CategoriaPermissaoService_Listar testa o método Listar do serviço de categoria de permissão.
func Test_CategoriaPermissaoService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() cpm.Repository
		filtro             cpm.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo cpm.Repository, categoriasPermissoes []cpm.CategoriaPermissao, total int, filtroRetornado cpm.Filtro)
	}{
		{
			nome: "listar 10 usuarios paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				popularCategoriaPermissaoTeste(repo, 10)
				return repo
			},
			filtro:       cpm.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriasPermissoes []cpm.CategoriaPermissao, total int, filtroRetornado cpm.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeCategoriaPermissaoRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.categoriasPermissoes))
				verificarLimite(t, categoriasPermissoes, 5)
			},
		},
		{
			nome: "erro ao listar categoriasPermissoes paginadas no repositório",
			prepararRepo: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				popularCategoriaPermissaoTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       cpm.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo cpm.Repository, categoriasPermissoes []cpm.CategoriaPermissao, total int, filtroRetornado cpm.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, categoriasPermissoes, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoCategoriaPermissaoService(repo)
			categoriasPermissoes, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, categoriasPermissoes, total, filtroRetornado)
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

func compararCategoriasPermissao(t *testing.T, antes, depois *cpm.CategoriaPermissao) {
	t.Helper()

	if antes == nil || depois == nil {
		t.Fatalf("categoriaPermissao para comparação não pode ser nil: antes=%v depois=%v", antes, depois)
	}

	if antes.CategoriaID() != depois.CategoriaID() ||
		antes.UsuarioID() != depois.UsuarioID() ||
		antes.Permissao() != depois.Permissao() {

		t.Errorf(
			"categoriaPermissao não deveriam ser diferentes: \nantes=%+v\ndepois=%+v",
			antes, depois,
		)
	}
}

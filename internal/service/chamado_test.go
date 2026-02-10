package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que fakeChamadoRepository implementa chm.Repository
var _ chm.Repository = (*fakeChamadoRepository)(nil)

// fakeChamadoRepository é uma implementação fake do repositório de chamados para testes.
type fakeChamadoRepository struct {
	chamados        map[string]*chm.Chamado
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
}

// novoFakeChamadoRepository cria uma nova instância do repositório fake.
func novoFakeChamadoRepository() *fakeChamadoRepository {
	return &fakeChamadoRepository{
		chamados: make(map[string]*chm.Chamado),
	}
}

// Criar adiciona um novo chamado ao repositório fake.
func (f *fakeChamadoRepository) Criar(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	chamado := c
	f.chamados[c.ID()] = &chamado
	return &chamado, nil
}

// Atualizar atualiza um chamado existente no repositório fake.
func (f *fakeChamadoRepository) Atualizar(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	chamado := c
	f.chamados[c.ID()] = &chamado
	return &chamado, nil
}

// BuscarPorID recupera um chamado pelo ID do repositório fake.
func (f *fakeChamadoRepository) BuscarPorID(ctx context.Context, id string) (*chm.Chamado, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	c, existe := f.chamados[id]
	if !existe {
		return nil, mysql.ErrChamadoNaoEncontrado
	}

	return c, nil
}

// Listar lista chamados do repositório fake com paginação.
func (f *fakeChamadoRepository) Listar(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para uma slice
	chamados := make([]chm.Chamado, 0, len(f.chamados))
	for _, c := range f.chamados {
		chamados = append(chamados, *c)
	}

	// Ordenação ORDER BY nome ASC
	sort.Slice(chamados, func(i, j int) bool {
		return chamados[i].Titulo() < chamados[j].Titulo()
	})

	// Contar o total antes da paginação
	total := len(chamados)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio >= total {
			return []chm.Chamado{}, total, nil
		}

		if fim > total {
			fim = total
		}

		chamados = chamados[inicio:fim]
	}

	return chamados, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const (
	chamadoTesteID = "chamado-123"
)

// novoChamadoTeste cria um novo chamado para testes.
func novoChamadoTeste() *chm.Chamado {
	c, _ := chm.Novo(
		chamadoTesteID,
		categoriaTesteID,
		subcategoriaTesteID,
		usuarioTesteID,
		"Título do Chamado",
		"Descrição detalhada do chamado para fins de teste.",
	)
	return c
}

// popularChamadosTeste popula o repositório fake com chamados de teste.
func popularChamadosTeste(repo *fakeChamadoRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("chamado-%02d", i)

		c, _ := chm.Novo(
			id,
			categoriaTesteID,
			subcategoriaTesteID,
			usuarioTesteID,
			fmt.Sprintf("Título do Chamado %02d", i),
			fmt.Sprintf("Descrição detalhada do chamado %02d para fins de teste.", i),
		)
		repo.chamados[id] = c
	}
}

// novoCriarChamadoParamsTeste cria parâmetros para criar um chamado para testes.
func novoCriarChamadoParamsTeste() chm.CriarParams {
	return chm.CriarParams{
		Titulo:         "Título do Chamado",
		Descricao:      "Descrição detalhada do chamado para fins de teste.",
		CategoriaID:    categoriaTesteID,
		SubcategoriaID: subcategoriaTesteID,
		CriadorID:      usuarioTesteID,
	}
}

// novoCriarChamadoParamsInvalidosTeste cria parâmetros inválidos para criar um chamado para testes.
func novoCriarChamadoParamsInvalidosTeste() chm.CriarParams {
	return chm.CriarParams{
		Titulo:         "",
		Descricao:      "",
		CategoriaID:    categoriaTesteID,
		SubcategoriaID: subcategoriaTesteID,
		CriadorID:      usuarioTesteID,
	}
}

// novoAtualizarChamadoParamsTeste cria parâmetros para atualizar um chamado para testes.
func novoAtualizarChamadoParamsTeste() chm.AtualizarParams {
	return chm.AtualizarParams{
		Titulo:         ptr("Título Atualizado do Chamado"),
		Descricao:      ptr("Descrição atualizada do chamado para fins de teste."),
		Arquivado:      ptr(false),
		CategoriaID:    ptr(categoriaTesteID),
		SubcategoriaID: ptr(subcategoriaTesteID),
	}
}

// novoAtualizarChamadoParamsInvalidosTeste cria parâmetros inválidos para atualizar um chamado para testes.
func novoAtualizarChamadoParamsInvalidosTeste() chm.AtualizarParams {
	return chm.AtualizarParams{
		Titulo:         ptr(""),
		Descricao:      ptr(""),
		Arquivado:      ptr(false),
		CategoriaID:    ptr(categoriaTesteID),
		SubcategoriaID: ptr(subcategoriaTesteID),
	}
}

// novoAtualizarStatusChamadoParamsTeste cria parâmetros para atualizar o status de um chamado para testes.
func novoAtualizarStatusChamadoParamsTeste() chm.AtualizarStatusParams {
	return chm.AtualizarStatusParams{Status: chm.StatusAberto}
}

// novoAtualizarStatusChamadoParamsInvalidosTeste cria parâmetros inválidos para atualizar o status de um chamado para testes.
func novoAtualizarStatusChamadoParamsInvalidosTeste() chm.AtualizarStatusParams {
	return chm.AtualizarStatusParams{Status: "status-invalido"}
}

// novoAtualizarSolucaoChamadoParamsTeste cria parâmetros para atualizar a solução de um chamado para testes.
func novoAtualizarSolucaoChamadoParamsTeste() chm.AtualizarSolucaoParams {
	return chm.AtualizarSolucaoParams{Solucao: "Solução detalhada para o chamado para fins de teste."}
}

// novoAtualizarSolucaoChamadoParamsInvalidosTeste cria parâmetros inválidos para atualizar a solução de um chamado para testes.
func novoAtualizarSolucaoChamadoParamsInvalidosTeste() chm.AtualizarSolucaoParams {
	return chm.AtualizarSolucaoParams{Solucao: ""}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestChamadoService_Criar testa o método Criar do serviço de chamado.
func TestChamadoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() chm.Repository
		params             chm.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "criar chamado com sucesso",
			geradorID: &fakeGeradorID{
				id: chamadoTesteID,
			},
			prepararRepo: func() chm.Repository {
				return novoFakeChamadoRepository()
			},
			params:       novoCriarChamadoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				if chamado.Titulo() != novoCriarChamadoParamsTeste().Titulo {
					t.Errorf("Título esperado '%s', recebido '%s'",
						novoCriarChamadoParamsTeste().Titulo, chamado.Titulo())
				}

				if chamado.Descricao() != novoCriarChamadoParamsTeste().Descricao {
					t.Errorf("Descrição esperada '%s', recebida '%s'",
						novoCriarChamadoParamsTeste().Descricao, chamado.Descricao())
				}

				if chamado.Status() != chm.StatusAberto {
					t.Errorf("Status esperado '%s', recebido '%s'",
						chm.StatusAberto, chamado.Status())
				}

				if chamado.Arquivado() != false {
					t.Errorf("Arquivado esperado 'false', recebido '%t'", chamado.Arquivado())
				}

				if chamado.Solucao() != nil {
					t.Errorf("Solução esperada 'nil', recebida '%v'", chamado.Solucao())
				}

				if chamado.SolucionadoEm() != nil {
					t.Errorf("SolucionadoEm esperado 'nil', recebido '%v'", chamado.SolucionadoEm())
				}

				if chamado.FechadoEm() != nil {
					t.Errorf("FechadoEm esperado 'nil', recebido '%v'", chamado.FechadoEm())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarIDs(t, categoriaTesteID, chamado.CategoriaID())
				verificarIDs(t, subcategoriaTesteID, chamado.SubcategoriaID())
				verificarIDs(t, usuarioTesteID, chamado.CriadorID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
				verificarNaoNulo(t, chamado, "esperava chamado criado, recebeu nil")
				verificarPersistido(t, len(repo.(*fakeChamadoRepository).chamados))
			},
		},
		{
			nome: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() chm.Repository {
				return novoFakeChamadoRepository()
			},
			params:       novoCriarChamadoParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeChamadoRepository).chamados))
				verificarNulo(t, chamado, "não deveria retornar chamado quando gerador de ID falha")
			},
		},
		{
			nome: "erro ao salvar chamado no repositório",
			geradorID: &fakeGeradorID{
				id: chamadoTesteID,
			},
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeChamadoRepository).chamados))
				verificarNulo(t, chamado, "não deveria retornar chamado quando repositório falha ao criar")
			},
		},
		{
			nome: "erro de validação ao criar chamado",
			geradorID: &fakeGeradorID{
				id: chamadoTesteID,
			},
			prepararRepo: func() chm.Repository {
				return novoFakeChamadoRepository()
			},
			params:       novoCriarChamadoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeChamadoRepository).chamados))
				verificarNulo(t, chamado, "não deveria retornar chamado quando parâmetros são inválidos")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			categoriaRepo := novoFakeCategoriaRepository()
			categoriaRepo.categorias[categoriaTesteID] = novoCategoriaTeste()
			subcategoriaRepo := novoFakeSubcategoriaRepository()
			subcategoriaRepo.subcategorias[subcategoriaTesteID] = novoSubcategoriaTeste()
			service := NovoChamadoService(tt.geradorID, repo, categoriaRepo, subcategoriaRepo, nil, nil, nil)
			chamadoCriado, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoCriado)
			}
		})
	}
}

// TestChamadoService_Atualizar testa a função Atualizar do ChamadoService.
func TestChamadoService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		params             chm.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "Atualizar chamado com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarChamadoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				if novoAtualizarChamadoParamsTeste().Titulo != nil {
					if chamado.Titulo() != *novoAtualizarChamadoParamsTeste().Titulo {
						t.Errorf("Título esperado '%s', recebido '%s'",
							*novoAtualizarChamadoParamsTeste().Titulo, chamado.Titulo())
					}
				}

				if novoAtualizarChamadoParamsTeste().Descricao != nil {
					if chamado.Descricao() != *novoAtualizarChamadoParamsTeste().Descricao {
						t.Errorf("Descrição esperada '%s', recebida '%s'",
							*novoAtualizarChamadoParamsTeste().Descricao, chamado.Descricao())
					}
				}

				if chamado.Arquivado() != *novoAtualizarChamadoParamsTeste().Arquivado {
					t.Errorf("Arquivado esperado '%t', recebido '%t'",
						*novoAtualizarChamadoParamsTeste().Arquivado, chamado.Arquivado())
				}

				if chamado.CategoriaID() != *novoAtualizarChamadoParamsTeste().CategoriaID {
					t.Errorf("CategoriaID esperado '%s', recebido '%s'",
						*novoAtualizarChamadoParamsTeste().CategoriaID, chamado.CategoriaID())
				}

				if chamado.SubcategoriaID() != *novoAtualizarChamadoParamsTeste().SubcategoriaID {
					t.Errorf("SubcategoriaID esperado '%s', recebido '%s'",
						*novoAtualizarChamadoParamsTeste().SubcategoriaID, chamado.SubcategoriaID())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
				verificarNulo(t, chamado, "esperava chamado atualizado, recebeu nil")
			},
		},
		{
			nome: "Erro ao buscar chamado no repositório para atualizar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
				verificarNulo(t, c, "não esperava chamado retornado quando há erro ao buscar")
			},
		},
		{
			nome: "erro ao persistir chamado atualizado no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
				verificarNulo(t, chamado, "não esperava chamado retornado quando há erro ao atualizar")
			},
		},
		{
			nome: "chamado não encontrado para atualizar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarChamadoParamsTeste(),
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
				verificarNulo(t, chamado, "não esperava chamado retornado quando chamado não é encontrado para atualizar")
			},
		},
		{
			nome: "Erro de validação ao atualizar chamado",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarChamadoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
				verificarNulo(t, chamado, "não esperava chamado retornado quando parâmetros são inválidos")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			categoriaRepo := novoFakeCategoriaRepository()
			categoriaRepo.categorias[categoriaTesteID] = novoCategoriaTeste()
			subcategoriaRepo := novoFakeSubcategoriaRepository()
			subcategoriaRepo.subcategorias[subcategoriaTesteID] = novoSubcategoriaTeste()
			service := NovoChamadoService(nil, repo, categoriaRepo, subcategoriaRepo, nil, nil, nil)
			chamadoAtualizado, erroRecebido := service.Atualizar(ctx, "chamado-123", tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoAtualizado)
			}
		})
	}
}

// TestChamadoService_BuscarPorID2 testa a função BuscarPorID do ChamadoService.
func TestChamadoService_BuscarPorID2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		repoSetup  func() chm.Repository
		id         string
		wantErr    bool
		assertFunc func(t *testing.T, c *chm.Chamado)
	}{
		{
			name: "Busca bem-sucedida de chamado por ID",
			repoSetup: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				repo.chamados["chamado-123"] = novoChamadoTeste()
				return repo
			},
			id:      "chamado-123",
			wantErr: false,
			assertFunc: func(t *testing.T, c *chm.Chamado) {
				t.Helper()

				if c.ID() != "chamado-123" {
					t.Errorf("ID do chamado incorreto. obtido: %s, esperado: %s", c.ID(), "chamado-123")
				}
			},
		},
		{
			name: "Erro ao buscar chamado inexistente",
			repoSetup: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				return repo
			},
			id:      "chamado-inexistente",
			wantErr: true,
			assertFunc: func(t *testing.T, c *chm.Chamado) {
				t.Helper()
				// Nenhum chamado deve ser retornado
				if c != nil {
					t.Errorf("Chamado retornado apesar de não existir. ID: %s", c.ID())
				}
			},
		},
		{
			name: "Erro ao buscar chamado devido a falha no repositório",
			repoSetup: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				repo.erroAoBuscar = errors.New("falha ao buscar chamado")
				return repo
			},
			id:      "chamado-123",
			wantErr: true,
			assertFunc: func(t *testing.T, c *chm.Chamado) {
				t.Helper()
				// Nenhum chamado deve ser retornado
				if c != nil {
					t.Errorf("Chamado retornado apesar do erro no repositório. ID: %s", c.ID())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)

			c, err := service.BuscarPorID(ctx, tt.id)

			assertError(t, err, tt.wantErr)

			if tt.assertFunc != nil {
				tt.assertFunc(t, c)
			}
		})
	}
}

// TestChamadoService_Listar2 testa a função Listar do ChamadoService.
func TestChamadoService_Listar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		repoSetup  func() chm.Repository
		filtro     chm.Filtro
		wantTotal  int
		wantErr    bool
		assertFunc func(t *testing.T, chamados []chm.Chamado, filtroRetornado chm.Filtro)
	}{
		{
			name: "Listagem bem-sucedida de chamados",
			repoSetup: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				repo.chamados["chamado-1"] = novoChamadoTeste()
				repo.chamados["chamado-2"] = novoChamadoTeste()
				return repo
			},
			filtro:    chm.Filtro{},
			wantTotal: 2,
			wantErr:   false,
			assertFunc: func(t *testing.T, chamados []chm.Chamado, filtroRetornado chm.Filtro) {
				t.Helper()

				if len(chamados) != 2 {
					t.Errorf("Número incorreto de chamados retornados. obtido: %d, esperado: %d", len(chamados), 2)
				}

				// filtro deve estar normalizado
				if filtroRetornado.Pagina() <= 0 {
					t.Errorf("Filtro retornado com página inválida: %d", filtroRetornado.Pagina())
				}
				if filtroRetornado.Limite() <= 0 {
					t.Errorf("Filtro retornado com limite inválido: %d", filtroRetornado.Limite())
				}
			},
		},
		{
			name: "Erro ao listar chamados devido a falha no repositório",
			repoSetup: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				repo.erroAoListar = errors.New("falha ao listar chamados")
				return repo
			},
			filtro:    chm.Filtro{},
			wantTotal: 0,
			wantErr:   true,
			assertFunc: func(t *testing.T, chamados []chm.Chamado, filtroRetornado chm.Filtro) {
				t.Helper()
				// Nenhum chamado deve ser retornado
				if len(chamados) != 0 {
					t.Errorf("Chamados retornados apesar do erro no repositório. Total de chamados: %d", len(chamados))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)

			chamados, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("Total de chamados incorreto. obtido: %d, esperado: %d", total, tt.wantTotal)
				}

				if tt.assertFunc != nil {
					tt.assertFunc(t, chamados, filtroRetornado)
				}
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

func compararChamados(t *testing.T, antes, depois *chm.Chamado) {
	t.Helper()

	if antes == nil || depois == nil {
		t.Fatalf(
			"Chamados para comparação não podem ser nil. antes: %v, depois: %v",
			antes, depois,
		)
	}

	solucaoDiferente :=
		(antes.Solucao() == nil) != (depois.Solucao() == nil) ||
			(antes.Solucao() != nil && depois.Solucao() != nil &&
				*antes.Solucao() != *depois.Solucao())

	solucionadoEmDiferente :=
		(antes.SolucionadoEm() == nil) != (depois.SolucionadoEm() == nil) ||
			(antes.SolucionadoEm() != nil && depois.SolucionadoEm() != nil &&
				!antes.SolucionadoEm().Equal(*depois.SolucionadoEm()))

	if antes.ID() != depois.ID() ||
		antes.CategoriaID() != depois.CategoriaID() ||
		antes.SubcategoriaID() != depois.SubcategoriaID() ||
		antes.CriadorID() != depois.CriadorID() ||
		antes.Titulo() != depois.Titulo() ||
		antes.Descricao() != depois.Descricao() ||
		antes.Status() != depois.Status() ||
		antes.Arquivado() != depois.Arquivado() ||
		solucaoDiferente ||
		solucionadoEmDiferente {

		t.Fatalf(
			"chamados não deveriam ser diferentes:\nantes=%+v\ndepois=%+v",
			antes, depois,
		)
	}
}

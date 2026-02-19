package service

import (
	"context"
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

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestChamadoService_Criar testa o método Criar do serviço de chamado.
func Test_ChamadoService_Criar(t *testing.T) {
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

				verificarNaoNulo(t, chamado)
				verificarPersistido(t, len(repo.(*fakeChamadoRepository).chamados))

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
				verificarNulo(t, chamado)
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
				verificarNulo(t, chamado)
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
				verificarNulo(t, chamado)
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
func Test_ChamadoService_Atualizar(t *testing.T) {
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
			nome: "atualizar chamado com sucesso",
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

				verificarNaoNulo(t, chamado)

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
			},
		},
		{
			nome: "erro ao buscar chamado no repositório para atualizar",
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
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
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
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
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
			id:           "id-inexistente",
			params:       novoAtualizarChamadoParamsTeste(),
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())

			},
		},
		{
			nome: "erro de validação ao atualizar chamado",
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
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())

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
			atendimentoRepo := novoFakeAtendimentoRepository()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimentoTeste()
			service := NovoChamadoService(nil, repo, categoriaRepo, subcategoriaRepo, nil, atendimentoRepo, nil)
			chamadoAtualizado, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoAtualizado)
			}
		})
	}
}

// TestChamadoService_BuscarPorID testa a função BuscarPorID do ChamadoService.
func Test_ChamadoService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "buscar chamado por id com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoNulo(t, chamado)
				verificarIDs(t, chamadoTesteID, chamado.ID())
			},
		},
		{
			nome: "erro ao buscar chamado por id no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           "chamado-inexistente",
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				_, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
			},
		},
		{
			nome: "chamado não encontrado  por id",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				_, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)

			chamadoEncontrado, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoEncontrado)
			}
		})
	}
}

// Test_ChamadoService_Arquivar testa a função Arquivar do ChamadoService.
func Test_ChamadoService_Arquivar(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "arquivar chamado com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				novoChamado.Desarquivar()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoNulo(t, chamado)

				if chamado.Arquivado() != true {
					t.Errorf("Arquivado esperado 'true', recebido '%t'", chamado.Arquivado())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar chamado no repositório para arquivar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir chamado arquivado no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "chamado não encontrado para arquivar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)
			chamadoArquivado, erroRecebido := service.Arquivar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoArquivado)
			}
		})
	}
}

// Test_ChamadoService_Desarquivar testa a função Desarquivar do ChamadoService.
func Test_ChamadoService_Desarquivar(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "desarquivar chamado com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				novoChamado.Arquivar()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				verificarNaoNulo(t, chamado)

				if chamado.Arquivado() != false {
					t.Errorf("Arquivado esperado 'false', recebido '%t'", chamado.Arquivado())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar chamado no repositório para desarquivar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir chamado desarquivado no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				novoChamado.Arquivar()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "chamado não encontrado para desarquivar",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)
			chamadoDesarquivado, erroRecebido := service.Desarquivar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoDesarquivado)
			}
		})
	}
}

// Test_ChamadoService_AtualizarStatus testa a função AtualizarStatus do ChamadoService.
func Test_ChamadoService_AtualizarStatus(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		params             chm.AtualizarStatusParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "atualizar status do chamado com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarStatusChamadoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				if chamado.Status() != novoAtualizarStatusChamadoParamsTeste().Status {
					t.Errorf("Status esperado '%s', recebido '%s'",
						novoAtualizarStatusChamadoParamsTeste().Status, chamado.Status())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
				verificarNaoNulo(t, chamado)
			},
		},
		{
			nome: "erro ao buscar chamado no repositório para atualizar status",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarStatusChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir chamado com status atualizado no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarStatusChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "chamado não encontrado para atualizar status",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarStatusChamadoParamsTeste(),
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
			{
			nome: "erro de validação ao atualizar status do chamado",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarStatusChamadoParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
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
			atendimentoRepo := novoFakeAtendimentoRepository()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimentoTeste()
			categoriaPermissaoRepo := novoFakeCategoriaPermissaoRepository()
			categoriaPermissaoRepo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()
			usuarioRepo := novoFakeUsuarioRepository()
			usuarioRepo.usuarios[usuarioTesteID] = novoUsuarioTeste()
			service := NovoChamadoService(nil, repo, categoriaRepo, subcategoriaRepo, categoriaPermissaoRepo, atendimentoRepo, nil)
			chamadoAtualizado, erroRecebido := service.AtualizarStatus(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoAtualizado)
			}
		})
	}
}

// Test_ChamadoService_AtualizarSolucao testa a função AtualizarSolucao do ChamadoService.
func Test_ChamadoService_AtualizarSolucao(t *testing.T) {
	t.Parallel()

	ctx := contextoTeste()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		id                 string
		params             chm.AtualizarSolucaoParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamado *chm.Chamado)
	}{
		{
			nome: "atualizar solução do chamado com sucesso",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarSolucaoChamadoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				if *chamado.Solucao() != novoAtualizarSolucaoChamadoParamsTeste().Solucao {
					t.Errorf("Solução esperada '%s', recebida '%v'",
						novoAtualizarSolucaoChamadoParamsTeste().Solucao, chamado.Solucao())
				}

				verificarIDs(t, chamadoTesteID, chamado.ID())
				verificarDatas(t, chamado.CriadoEm(), chamado.AtualizadoEm())
				verificarNaoNulo(t, chamado)
			},
		},
		{
			nome: "erro ao buscar chamado no repositório para atualizar solução",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarSolucaoChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir chamado com solução atualizada no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           chamadoTesteID,
			params:       novoAtualizarSolucaoChamadoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
			},
		},
		{
			nome: "chamado não encontrado para atualizar solução",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				novoChamado := novoChamadoTeste()
				repo.chamados[novoChamado.ID()] = novoChamado
				repo.erroAoBuscar = mysql.ErrChamadoNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarSolucaoChamadoParamsTeste(),
			erroEsperado: mysql.ErrChamadoNaoEncontrado,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamado *chm.Chamado) {
				t.Helper()

				chamadoRepo, ok := repo.(*fakeChamadoRepository).chamados[chamadoTesteID]
				verificarOk(t, ok, "deveria existir chamado no repositório")
				verificarNulo(t, chamado)
				verificarNaoAlterado(t, novoChamadoTeste(), chamadoRepo, compararChamados)
				verificarDatas(t, chamadoRepo.CriadoEm(), chamadoRepo.AtualizadoEm())
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
			atendimentoRepo := novoFakeAtendimentoRepository()
			atendimentoRepo.atendimentos[atendimentoTesteID] = novoAtendimentoTeste()
			categoriaPermissaoRepo := novoFakeCategoriaPermissaoRepository()
			categoriaPermissaoRepo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()
			usuarioRepo := novoFakeUsuarioRepository()
			usuarioRepo.usuarios[usuarioTesteID] = novoUsuarioTeste()
			service := NovoChamadoService(nil, repo, categoriaRepo, subcategoriaRepo, categoriaPermissaoRepo, atendimentoRepo, nil)
			chamadoAtualizado, erroRecebido := service.AtualizarSolucao(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamadoAtualizado)
			}
		})
	}
}

// Test_ChamadoService_Listar testa a função Listar do ChamadoService.
func Test_ChamadoService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() chm.Repository
		filtro             chm.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo chm.Repository, chamados []chm.Chamado, total int, filtroRetornado chm.Filtro)
	}{
		{
			nome: "listar 10 chamados paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				popularChamadosTeste(repo, 10)
				return repo
			},
			filtro:       chm.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamados []chm.Chamado, total int, filtroRetornado chm.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeChamadoRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.chamados))
				verificarLimite(t, chamados, 5)
			},
		},
		{
			nome: "erro ao listar chamados paginados no repositório",
			prepararRepo: func() chm.Repository {
				repo := novoFakeChamadoRepository()
				popularChamadosTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       chm.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil, nil, nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo chm.Repository, chamados []chm.Chamado, total int, filtroRetornado chm.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, chamados, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoChamadoService(nil, repo, nil, nil, nil, nil, nil)
			chamados, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, chamados, total, filtroRetornado)
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

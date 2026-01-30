package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeChamadoRepository2 é uma implementação fake do repositório de chamados para testes.
type fakeChamadoRepository2 struct {
	chamados        map[string]*chm.Chamado
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
	totalListar     int
}

// NovoFakeChamadoRepository2 cria uma nova instância do repositório fake.
func NovoFakeChamadoRepository2() *fakeChamadoRepository2 {
	return &fakeChamadoRepository2{
		chamados: make(map[string]*chm.Chamado),
	}
}

// Criar adiciona um novo chamado ao repositório fake.
func (f *fakeChamadoRepository2) Criar(ctx context.Context, u chm.Chamado) (*chm.Chamado, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	f.chamados[u.ID()] = &u
	return &u, nil
}

// Atualizar atualiza um chamado existente no repositório fake.
func (f *fakeChamadoRepository2) Atualizar(ctx context.Context, id string, u chm.Chamado) (*chm.Chamado, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	f.chamados[u.ID()] = &u
	return &u, nil
}

// BuscarPorID recupera um chamado pelo ID do repositório fake.
func (f *fakeChamadoRepository2) BuscarPorID(ctx context.Context, id string) (*chm.Chamado, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	c, existe := f.chamados[id]
	if !existe {
		return nil, errors.New("chamado não encontrado")
	}

	return c, nil
}

// Listar lista chamados do repositório fake com paginação.
func (f *fakeChamadoRepository2) Listar(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	chamados := make([]chm.Chamado, 0, len(f.chamados))
	for _, c := range f.chamados {
		chamados = append(chamados, *c)
	}

	total := f.totalListar
	if total == 0 {
		total = len(chamados)
	}

	return chamados, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

// novoChamadoTeste cria um novo chamado para testes.
func novoChamadoTeste(id string) *chm.Chamado {
	c, _ := chm.Novo(
		id,
		"categoria-123",
		"subcategoria-123",
		"criador-123",
		"Título do Chamado",
		"Descrição detalhada do chamado para fins de teste.",
	)
	return c
}

// novoCriarChamadoParamsTeste cria parâmetros para criar um chamado para testes.
func novoCriarChamadoParamsTeste() chm.CriarParams {
	return chm.CriarParams{
		Titulo:         "Título do Chamado",
		Descricao:      "Descrição detalhada do chamado para fins de teste.",
		CategoriaID:    "categoria-123",
		SubcategoriaID: "subcategoria-123",
		CriadorID:      "criador-123",
	}
}

// novoAtualizarChamadoParamsTeste cria parâmetros para atualizar um chamado para testes.
func novoAtualizarChamadoParamsTeste() chm.AtualizarParams {
	return chm.AtualizarParams{
		Titulo:         ptr("Título Atualizado do Chamado"),
		Descricao:      ptr("Descrição atualizada do chamado para fins de teste."),
		Arquivado:      ptr(false),
		CategoriaID:    ptr("categoria-456"),
		SubcategoriaID: ptr("subcategoria-456"),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

func TestChamadoService_Criar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		geradorID  dmn.GeradorID
		repoSetup  func() chm.Repository
		params     chm.CriarParams
		wantErr    bool
		assertFunc func(t *testing.T, repo chm.Repository, c *chm.Chamado)
	}{
		{
			name: "Criação bem-sucedida de chamado",
			geradorID: &fakeGeradorID{
				id: "chamado-123",
			},
			repoSetup: func() chm.Repository {
				return NovoFakeChamadoRepository2()
			},
			params:  novoCriarChamadoParamsTeste(),
			wantErr: false,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()

				if c.ID() != "chamado-123" {
					t.Errorf("ID do chamado incorreto. obtido: %s, esperado: %s", c.ID(), "chamado-123")
				}

				fake := repo.(*fakeChamadoRepository2)
				if len(fake.chamados) != 1 {
					t.Errorf("Número incorreto de chamados no repositório. obtido: %d, esperado: %d", len(fake.chamados), 1)
				}
			},
		},
		{
			name: "Erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errors.New("falha ao gerar ID"),
			},
			repoSetup: func() chm.Repository {
				return NovoFakeChamadoRepository2()
			},
			params:  novoCriarChamadoParamsTeste(),
			wantErr: true,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()

				// repositório não deve ter nenhum chamado criado
				fake := repo.(*fakeChamadoRepository2)
				if len(fake.chamados) != 0 {
					t.Errorf("Chamado foi criado no repositório apesar do erro ao gerar ID. Total de chamados: %d", len(fake.chamados))
				}
			},
		},
		{
			name: "Erro ao salvar chamado no repositório",
			geradorID: &fakeGeradorID{
				id: "chamado-123",
			},
			repoSetup: func() chm.Repository {
				repo := NovoFakeChamadoRepository2()
				repo.erroAoCriar = errors.New("falha ao salvar chamado")
				return repo
			},
			params:  novoCriarChamadoParamsTeste(),
			wantErr: true,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()

				// repositório não deve ter nenhum chamado criado
				fake := repo.(*fakeChamadoRepository2)
				if len(fake.chamados) != 0 {
					t.Errorf("Chamado foi criado no repositório apesar do erro ao salvar. Total de chamados: %d", len(fake.chamados))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			categoriaRepo := novoFakeCategoriaRepository2()
			categoriaRepo.categorias["categoria-123"] = novoCategoriaTeste("categoria-123", "Categoria Teste")
			subcategoriaRepo := novoFakeSubcategoriaRepository2()
			subcategoriaRepo.subcategorias["subcategoria-123"] = novoSubcategoriaTeste("subcategoria-123", "Subcategoria Teste", "categoria-123")

			service := NovoChamadoService(tt.geradorID, repo, categoriaRepo, subcategoriaRepo, nil, nil, nil)

			c, err := service.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFunc != nil {
				tt.assertFunc(t, repo, c)
			}
		})
	}
}

// TestChamadoService_Atualizar2 testa a função Atualizar do ChamadoService.
func TestChamadoService_Atualizar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "criador-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name       string
		repoSetup  func() chm.Repository
		params     chm.AtualizarParams
		wantErr    bool
		assertFunc func(t *testing.T, repo chm.Repository, c *chm.Chamado)
	}{
		{
			name: "Atualização bem-sucedida de chamado",
			repoSetup: func() chm.Repository {
				repo := NovoFakeChamadoRepository2()
				repo.chamados["chamado-123"] = novoChamadoTeste("chamado-123")
				return repo
			},
			params:  novoAtualizarChamadoParamsTeste(),
			wantErr: false,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()

				if c.Titulo() != "Título Atualizado do Chamado" {
					t.Errorf("Título do chamado incorreto. obtido: %s, esperado: %s", c.Titulo(), "Título Atualizado do Chamado")
				}

				fake := repo.(*fakeChamadoRepository2)
				armazenado, existe := fake.chamados["chamado-123"]
				if !existe {
					t.Errorf("Chamado não encontrado no repositório após atualização.")
					return
				}
				if armazenado.Titulo() != "Título Atualizado do Chamado" {
					t.Errorf("Título do chamado no repositório incorreto. obtido: %s, esperado: %s", armazenado.Titulo(), "Título Atualizado do Chamado")
				}
			},
		},
		{
			name: "Erro ao buscar chamado existente",
			repoSetup: func() chm.Repository {
				repo := NovoFakeChamadoRepository2()
				repo.erroAoBuscar = errors.New("falha ao buscar chamado")
				return repo
			},
			params:  novoAtualizarChamadoParamsTeste(),
			wantErr: true,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()
				// Nenhum chamado deve ser atualizado no repositório
				fake := repo.(*fakeChamadoRepository2)
				if len(fake.chamados) != 0 {
					t.Errorf("Chamado foi atualizado no repositório apesar do erro ao buscar. Total de chamados: %d", len(fake.chamados))
				}
			},
		},
		{
			name: "Erro ao atualizar chamado no repositório",
			repoSetup: func() chm.Repository {
				repo := NovoFakeChamadoRepository2()
				repo.chamados["chamado-123"] = novoChamadoTeste("chamado-123")
				repo.erroAoAtualizar = errors.New("falha ao atualizar chamado")
				return repo
			},
			params:  novoAtualizarChamadoParamsTeste(),
			wantErr: true,
			assertFunc: func(t *testing.T, repo chm.Repository, c *chm.Chamado) {
				t.Helper()
				// Chamado deve permanecer inalterado no repositório
				fake := repo.(*fakeChamadoRepository2)
				armazenado, existe := fake.chamados["chamado-123"]
				if !existe {
					t.Errorf("Chamado não encontrado no repositório após tentativa de atualização.")
					return
				}
				if armazenado.Titulo() != "Título do Chamado" {
					t.Errorf("Título do chamado no repositório foi alterado apesar do erro na atualização. obtido: %s, esperado: %s", armazenado.Titulo(), "Título do Chamado")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			categoriaRepo := novoFakeCategoriaRepository2()
			categoriaRepo.categorias["categoria-456"] = novoCategoriaTeste("categoria-456", "Categoria Atualizada")
			subcategoriaRepo := novoFakeSubcategoriaRepository2()
			subcategoriaRepo.subcategorias["subcategoria-456"] = novoSubcategoriaTeste("subcategoria-456", "Subcategoria Atualizada", "categoria-456")

			service := NovoChamadoService(nil, repo, categoriaRepo, subcategoriaRepo, nil, nil, nil)

			c, err := service.Atualizar(ctx, "chamado-123", tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFunc != nil {
				tt.assertFunc(t, repo, c)
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
				repo := NovoFakeChamadoRepository2()
				repo.chamados["chamado-123"] = novoChamadoTeste("chamado-123")
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
				repo := NovoFakeChamadoRepository2()
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
				repo := NovoFakeChamadoRepository2()
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
				repo := NovoFakeChamadoRepository2()
				repo.chamados["chamado-1"] = novoChamadoTeste("chamado-1")
				repo.chamados["chamado-2"] = novoChamadoTeste("chamado-2")
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
				repo := NovoFakeChamadoRepository2()
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

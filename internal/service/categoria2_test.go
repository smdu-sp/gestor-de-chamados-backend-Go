package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeCategoriaRepository2 é uma implementação fake do repositório de categorias para testes.
type fakeCategoriaRepository2 struct {
	categorias      map[string]*ctg.Categoria
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
	totalListar     int
}

// novoFakeCategoriaRepository2 cria uma nova instância do repositório fake.
func novoFakeCategoriaRepository2() *fakeCategoriaRepository2 {
	return &fakeCategoriaRepository2{
		categorias: make(map[string]*ctg.Categoria),
	}
}

// Criar adiciona uma nova categoria ao repositório fake.
func (f *fakeCategoriaRepository2) Criar(ctx context.Context, u ctg.Categoria) (*ctg.Categoria, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	f.categorias[u.ID()] = &u
	return &u, nil
}

// Atualizar atualiza uma categoria existente no repositório fake.
func (f *fakeCategoriaRepository2) Atualizar(ctx context.Context, id string, u ctg.Categoria) (*ctg.Categoria, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	f.categorias[id] = &u
	return &u, nil
}

// BuscarPorID recupera uma categoria pelo ID do repositório fake.
func (f *fakeCategoriaRepository2) BuscarPorID(ctx context.Context, id string) (*ctg.Categoria, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	c, existe := f.categorias[id]
	if !existe {
		return nil, errors.New("categoria não encontrada")
	}

	return c, nil
}

// BuscarPorNome recupera uma categoria pelo nome do repositório fake.
func (f *fakeCategoriaRepository2) BuscarPorNome(ctx context.Context, nome string) (*ctg.Categoria, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	
	c, existe := f.categorias[nome]
	if !existe {
		return nil, errors.New("categoria não encontrada")
	}

	return c, nil
}

// Listar lista categorias do repositório fake com paginação.
func (f *fakeCategoriaRepository2) Listar(ctx context.Context, filtro ctg.Filtro) ([]ctg.Categoria, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	categorias := make([]ctg.Categoria, 0, len(f.categorias))
	for _, c := range f.categorias {
		categorias = append(categorias, *c)
	}

	total := f.totalListar
	if total == 0 {
		total = len(categorias)
	}

	return categorias, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

// novoCategoriaTeste cria uma nova categoria para testes.
func novoCategoriaTeste(id, nome string) *ctg.Categoria {
	c, _ := ctg.Novo(id, nome)
	return c
}

// novoCriarCategoriaParamsTeste cria parâmetros de criação de categoria para testes.
func novoCriarCategoriaParamsTeste() ctg.CriarParams {
	return ctg.CriarParams{
		Nome: "Categoria de Teste",
	}
}

// novoAtualizarCategoriaParamsTeste cria parâmetros de atualização de categoria para testes.
func novoAtualizarCategoriaParamsTeste() ctg.AtualizarParams {
	return ctg.AtualizarParams{
		Nome:   ptr("Categoria Atualizada"),
		Status: ptr(true),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestCategoriaService_Criar2 testa o método Criar do CategoriaService.
func TestCategoriaService_Criar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		geradorID dmn.GeradorID
		repoSetup   func() ctg.Repository
		params      ctg.CriarParams
		wantErr     bool
		assertFn 	func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			name: "criar categoria com sucesso",
			geradorID: &fakeGeradorID{
				id: "categoria-123",
			},
			repoSetup: func() ctg.Repository {
				return novoFakeCategoriaRepository2()
			},
			params:  novoCriarCategoriaParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				if c.ID() != "categoria-123" {
					t.Errorf("ID esperado 'categoria-123', recebeu '%s'", c.ID())
				}
				if c.Nome() != "Categoria de Teste" {
					t.Errorf("Nome esperado 'Categoria de Teste', recebeu '%s'", c.Nome())
				}

				// Verificar se a categoria foi salva no repositório
				fake := repo.(*fakeCategoriaRepository2)
				if len(fake.categorias) != 1 {
					t.Errorf("Esperava 1 categoria no repositório, recebeu %d", len(fake.categorias))
				}
			},
		},
		{
			name: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errors.New("falha ao gerar ID"),
			},
			repoSetup: func() ctg.Repository {
				return novoFakeCategoriaRepository2()
			},
			params:  novoCriarCategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				fake := repo.(*fakeCategoriaRepository2)
				if len(fake.categorias) != 0 {
					t.Errorf("Esperava 0 categorias no repositório, recebeu %d", len(fake.categorias))
				}
			},
		},
		{
			name: "erro ao salvar categoria no repositório",
			geradorID: &fakeGeradorID{
				id: "categoria-123",
			},
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.erroAoCriar = errors.New("falha ao salvar categoria")
				return repo
			},
			params:  novoCriarCategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				fake := repo.(*fakeCategoriaRepository2)
				if len(fake.categorias) != 0 {
					t.Errorf("Esperava 0 categorias no repositório, recebeu %d", len(fake.categorias))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaService(tt.geradorID, repo)

			c, err := service.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, c)
			}
		})
	}
}

// TestCategoriaService_Atualizar2 testa o método Atualizar do CategoriaService.
func TestCategoriaService_Atualizar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		repoSetup   func() ctg.Repository
		params      ctg.AtualizarParams
		wantErr     bool
		assertFn    func(t *testing.T, repo ctg.Repository, categoria *ctg.Categoria)
	}{
		{
			name: "atualizar categoria com sucesso",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.categorias["categoria-123"] = novoCategoriaTeste("categoria-123", "Categoria Antiga")
				return repo
			},
			params:  novoAtualizarCategoriaParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				if c.Nome() != "Categoria Atualizada" {
					t.Errorf("Nome esperado 'Categoria Atualizada', recebeu '%s'", c.Nome())
				}

				// Verificar se a categoria foi atualizada no repositório
				fake := repo.(*fakeCategoriaRepository2)
				atualizada, existe := fake.categorias["categoria-123"]
				if !existe {
					t.Errorf("Categoria não encontrada no repositório")
				}
				if atualizada.Nome() != "Categoria Atualizada" {
					t.Errorf("Nome no repositório esperado 'Categoria Atualizada', recebeu '%s'", atualizada.Nome())
				}
			},
		},
		{
			name: "erro ao buscar categoria inexistente",
			repoSetup: func() ctg.Repository {
				return novoFakeCategoriaRepository2()
			},
			params:  novoAtualizarCategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				// Verificar se nenhuma categoria foi adicionada no repositório
				fake := repo.(*fakeCategoriaRepository2)
				if len(fake.categorias) != 0 {
					t.Errorf("Esperava 0 categorias no repositório, recebeu %d", len(fake.categorias))
				}
			},
		},
		{
			name: "erro ao atualizar categoria no repositório",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.categorias["categoria-123"] = novoCategoriaTeste("categoria-123", "Categoria Antiga")
				repo.erroAoAtualizar = errors.New("falha ao atualizar categoria")
				return repo
			},
			params:  novoAtualizarCategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				// Verificar se a categoria não foi alterada no repositório
				fake := repo.(*fakeCategoriaRepository2)
				atual, existe := fake.categorias["categoria-123"]
				if !existe {
					t.Errorf("Categoria não encontrada no repositório")
				}
				if atual.Nome() != "Categoria Antiga" {
					t.Errorf("Nome no repositório esperado 'Categoria Antiga', recebeu '%s'", atual.Nome())
				}
			},
		},
		{
			name: "erro ao buscar categoria devido a erro no repositório",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.erroAoBuscar = errors.New("falha ao buscar categoria")
				return repo
			},
			params:  novoAtualizarCategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo ctg.Repository, c *ctg.Categoria) {
				t.Helper()

				// Verificar se nenhuma categoria foi adicionada no repositório
				fake := repo.(*fakeCategoriaRepository2)
				if len(fake.categorias) != 0 {
					t.Errorf("Esperava 0 categorias no repositório, recebeu %d", len(fake.categorias))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaService(nil, repo)

			c, err := service.Atualizar(ctx, "categoria-123", tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, c)
			}
		})
	}
}

// TestCategoriaService_BuscarPorID2 testa o método BuscarPorID do CategoriaService.
func TestCategoriaService_BuscarPorID2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() ctg.Repository
		id        string
		wantErr   bool
		assertFn  func(t *testing.T, categoria *ctg.Categoria)
	}{
		{
			name: "buscar categoria por ID com sucesso",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.categorias["categoria-123"] = novoCategoriaTeste("categoria-123", "Categoria de Teste")
				return repo
			},
			id:      "categoria-123",
			wantErr: false,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c.ID() != "categoria-123" {
					t.Errorf("ID esperado 'categoria-123', recebeu '%s'", c.ID())
				}
				if c.Nome() != "Categoria de Teste" {
					t.Errorf("Nome esperado 'Categoria de Teste', recebeu '%s'", c.Nome())
				}
			},
		},
		{
			name: "erro ao buscar categoria inexistente por ID",
			repoSetup: func() ctg.Repository {
				return novoFakeCategoriaRepository2()
			},
			id:      "categoria-inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c != nil {
					t.Errorf("Esperava categoria nula, recebeu '%v'", c)
				}
			},
		},
		{
			name: "erro ao buscar categoria por ID devido a erro no repositório",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.erroAoBuscar = errors.New("falha ao buscar categoria")
				return repo
			},
			id:      "categoria-123",
			wantErr: true,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c != nil {
					t.Errorf("Esperava categoria nula, recebeu '%v'", c)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaService(nil, repo)

			c, err := service.BuscarPorID(ctx, tt.id)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, c)
			}
		})
	}
}

// TestCategoriaService_BuscarPorNome2 testa o método BuscarPorNome do CategoriaService.
func TestCategoriaService_BuscarPorNome2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() ctg.Repository
		nome      string
		wantErr   bool
		assertFn  func(t *testing.T, categoria *ctg.Categoria)
	}{
		{
			name: "buscar categoria por nome com sucesso",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.categorias["Categoria de Teste"] = novoCategoriaTeste("categoria-123", "Categoria de Teste")
				return repo
			},
			nome:    "Categoria de Teste",
			wantErr: false,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c.ID() != "categoria-123" {
					t.Errorf("ID esperado 'categoria-123', recebeu '%s'", c.ID())
				}
				if c.Nome() != "Categoria de Teste" {
					t.Errorf("Nome esperado 'Categoria de Teste', recebeu '%s'", c.Nome())
				}
			},
		},
		{
			name: "erro ao buscar categoria inexistente por nome",
			repoSetup: func() ctg.Repository {
				return novoFakeCategoriaRepository2()
			},
			nome:    "Categoria Inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c != nil {
					t.Errorf("Esperava categoria nula, recebeu '%v'", c)
				}
			},
		},
		{
			name: "erro ao buscar categoria por nome devido a erro no repositório",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.erroAoBuscar = errors.New("falha ao buscar categoria")
				return repo
			},
			nome:    "Categoria de Teste",
			wantErr: true,
			assertFn: func(t *testing.T, c *ctg.Categoria) {
				t.Helper()

				if c != nil {
					t.Errorf("Esperava categoria nula, recebeu '%v'", c)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaService(nil, repo)

			c, err := service.BuscarPorNome(ctx, tt.nome)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, c)
			}
		})
	}
}

// TestCategoriaService_Listar2 testa o método Listar do CategoriaService.
func TestCategoriaService_Listar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		repoSetup   func() ctg.Repository
		filtro      ctg.Filtro
		wantTotal  int
		wantErr     bool
		assertFn    func(t *testing.T, categorias []ctg.Categoria, filtroRetornado ctg.Filtro)
	}{
		{
			name: "listar categorias com sucesso",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.categorias["categoria-1"] = novoCategoriaTeste("categoria-1", "Categoria 1")
				repo.categorias["categoria-2"] = novoCategoriaTeste("categoria-2", "Categoria 2")
				return repo
			},
			filtro: ctg.Filtro{},
			wantTotal: 2,
			wantErr:   false,
			assertFn: func(t *testing.T, categorias []ctg.Categoria, filtroRetornado ctg.Filtro) {
				t.Helper()
				if len(categorias) != 2 {
					t.Errorf("Esperava 2 categorias, recebeu %d", len(categorias))
				}

				// Filtro deve estar normalizado
				if filtroRetornado.Pagina() <= 0 {
					t.Error("esperava filtro.Pagina normalizado")
				}
				if filtroRetornado.Limite() <= 0 {
					t.Error("esperava filtro.Limite normalizado")
				}
			},
		},
		{
			name: "erro ao listar categorias devido a erro no repositório",
			repoSetup: func() ctg.Repository {
				repo := novoFakeCategoriaRepository2()
				repo.erroAoListar = errors.New("falha ao listar categorias")
				return repo
			},
			filtro: ctg.Filtro{},
			wantTotal: 0,
			wantErr:   true,
			assertFn: func(t *testing.T, categorias []ctg.Categoria, filtroRetornado ctg.Filtro) {
				t.Helper()
				if categorias != nil {
					t.Errorf("Esperava categorias nulo, recebeu %v", categorias)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaService(nil, repo)

			categorias, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("Total esperado %d, recebeu %d", tt.wantTotal, total)
				}

				if tt.assertFn != nil {
					tt.assertFn(t, categorias, filtroRetornado)
				}
			}
		})
	}
}
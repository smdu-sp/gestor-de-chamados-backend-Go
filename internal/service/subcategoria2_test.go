package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeSubcategoriaRepository2 é um repositório fake para testes, com estado interno.
type fakeSubcategoriaRepository2 struct {
	subcategorias   map[string]*subc.Subcategoria
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
	totalListar     int
}

// novoFakeSubcategoriaRepository2 cria uma nova instância do repositório fake de subcategorias.
func novoFakeSubcategoriaRepository2() *fakeSubcategoriaRepository2 {
	return &fakeSubcategoriaRepository2{
		subcategorias: make(map[string]*subc.Subcategoria),
	}
}

// Implementação dos métodos da interface sub.Repository
func (r *fakeSubcategoriaRepository2) Criar(ctx context.Context, s subc.Subcategoria) (*subc.Subcategoria, error) {
	if r.erroAoCriar != nil {
		return nil, r.erroAoCriar
	}
	r.subcategorias[s.ID()] = &s
	return &s, nil
}

// Atualizar atualiza uma subcategoria existente no repositório fake.
func (r *fakeSubcategoriaRepository2) Atualizar(ctx context.Context, id string, s subc.Subcategoria) (*subc.Subcategoria, error) {
	if r.erroAoCriar != nil {
		return nil, r.erroAoAtualizar
	}
	r.subcategorias[id] = &s
	return &s, nil
}

// BuscarPorID busca uma subcategoria por ID no repositório fake.
func (r *fakeSubcategoriaRepository2) BuscarPorID(ctx context.Context, id string) (*subc.Subcategoria, error) {
	if r.erroAoBuscar != nil {
		return nil, r.erroAoBuscar
	}
	s, existe := r.subcategorias[id]
	if !existe {
		return nil, errors.New("subcategoria não encontrada")
	}
	return s, nil
}

// BuscarPorID busca uma subcategoria por ID no repositório fake.
func (r *fakeSubcategoriaRepository2) BuscarPorNome(ctx context.Context, nome string) (*subc.Subcategoria, error) {
	if r.erroAoBuscar != nil {
		return nil, r.erroAoBuscar
	}
	for _, s := range r.subcategorias {
		if s.Nome() == nome {
			return s, nil
		}
	}
	return nil, errors.New("subcategoria não encontrada")
}

// Listar lista subcategorias no repositório fake com base no filtro fornecido.
func (r *fakeSubcategoriaRepository2) Listar(ctx context.Context, filtro subc.Filtro) ([]subc.Subcategoria, int, error) {
	if r.erroAoListar != nil {
		return nil, 0, r.erroAoListar
	}

	subcategorias := make([]subc.Subcategoria, 0, len(r.subcategorias))
	for _, s := range r.subcategorias {
		subcategorias = append(subcategorias, *s)
	}

	total := r.totalListar
	if total == 0 {
		total = len(subcategorias)
	}

	return subcategorias, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

// novoSubcategoriaTeste cria uma subcategoria de teste com os valores fornecidos.
func novoSubcategoriaTeste(id, nome string, categoriaID string) *subc.Subcategoria {
	s, _ := subc.Novo(
		id,
		nome,
		categoriaID,
	)
	return s
}

// novoCriarSubcategoriaParamsTeste cria parâmetros de criação de subcategoria para testes.
func novoCriarSubcategoriaParamsTeste() subc.CriarParams {
	return subc.CriarParams{
		Nome:        "Suporte Técnico",
		CategoriaID: "cat-123",
	}
}

// novoAtualizarSubcategoriaParamsTeste cria parâmetros de atualização de subcategoria para testes.
func novoAtualizarSubcategoriaParamsTeste() subc.AtualizarParams {
	return subc.AtualizarParams{
		Nome: ptr("Suporte Avançado"),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestServicoSubcategoria_Criar testa o método Criar do serviço de subcategorias.
func TestServicoSubcategoria_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repoSetup func() subc.Repository
		params    subc.CriarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo subc.Repository, s *subc.Subcategoria)
	}{
		{
			name: "criar subcategoria com sucesso",
			geradorID: &fakeGeradorID{
				id: "subcat-123",
			},
			repoSetup: func() subc.Repository {
				return novoFakeSubcategoriaRepository2()
			},
			params:  novoCriarSubcategoriaParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				if s.ID() != "subcat-123" {
					t.Errorf("ID esperado 'subcat-123', got '%s'", s.ID())
				}
				if s.Nome() != "Suporte Técnico" {
					t.Errorf("Nome esperado 'Suporte Técnico', got '%s'", s.Nome())
				}

				// Verifica se a subcategoria foi salva no repositório
				fake := repo.(*fakeSubcategoriaRepository2)
				if len(fake.subcategorias) != 1 {
					t.Errorf("Esperado 1 subcategoria no repositório, tem %d", len(fake.subcategorias))
				}
			},
		},
		{
			name: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errors.New("falha ao gerar ID"),
			},
			repoSetup: func() subc.Repository {
				return novoFakeSubcategoriaRepository2()
			},
			params:  novoCriarSubcategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				// Repositório não deve ter sido chamado
				fake := repo.(*fakeSubcategoriaRepository2)
				if len(fake.subcategorias) != 0 {
					t.Errorf("Esperado 0 subcategorias no repositório, tem %d", len(fake.subcategorias))
				}
			},
		},
		{
			name: "erro ao salvar no repositorio",
			geradorID: &fakeGeradorID{
				id: "subcat-123",
			},
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.erroAoCriar = errors.New("erro no banco")
				return repo
			},
			params:  novoCriarSubcategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				// Repositório não deve ter subcategorias salvas
				fake := repo.(*fakeSubcategoriaRepository2)
				if len(fake.subcategorias) != 0 {
					t.Errorf("Esperado 0 subcategorias no repositório, tem %d", len(fake.subcategorias))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			servico := NovoSubcategoriaService(tt.geradorID, repo)

			s, err := servico.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, s)
			}
		})
	}
}

// TestServicoSubcategoria_Atualizar2 testa o método Atualizar do serviço de subcategorias.
func TestServicoSubcategoria_Atualizar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() subc.Repository
		params    subc.AtualizarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo subc.Repository, s *subc.Subcategoria)
	}{
		{
			name: "atualizar subcategoria com sucesso",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.subcategorias["subcat-123"] = novoSubcategoriaTeste("subcat-123", "Suporte Antigo", "cat-456")
				return repo
			},
			params:  novoAtualizarSubcategoriaParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				if s.Nome() != "Suporte Avançado" {
					t.Errorf("Nome esperado 'Suporte Avançado', got '%s'", s.Nome())
				}

				// Repositório deve ter a subcategoria atualizada
				fake := repo.(*fakeSubcategoriaRepository2)
				updated, exists := fake.subcategorias["subcat-123"]
				if !exists {
					t.Errorf("Subcategoria com ID 'subcat-123' não encontrada no repositório")
				}
				if updated.Nome() != "Suporte Avançado" {
					t.Errorf("Nome esperado 'Suporte Avançado', recebido '%s'", updated.Nome())
				}
			},
		},
		{
			name: "erro ao buscar subcategoria",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			params:  subc.AtualizarParams{},
			wantErr: true,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				// Repositório não deve ter subcategorias salvas
				fake := repo.(*fakeSubcategoriaRepository2)
				if len(fake.subcategorias) != 0 {
					t.Errorf("Esperado 0 subcategorias no repositório, tem %d", len(fake.subcategorias))
				}
			},
		},
		{
			name: "erro ao salvar atualização",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.subcategorias["subcat-123"] = novoSubcategoriaTeste("subcat-123", "Suporte Antigo", "cat-456")
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			params:  novoAtualizarSubcategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				// Repositório deve manter a subcategoria original sem alterações
				fake := repo.(*fakeSubcategoriaRepository2)
				original, exists := fake.subcategorias["subcat-123"]
				if !exists {
					t.Errorf("Subcategoria com ID 'subcat-123' não encontrada no repositório")
				}
				if original.Nome() != "Suporte Antigo" {
					t.Errorf("Nome esperado 'Suporte Antigo', recebido '%s'", original.Nome())
				}
			},
		},
		{
			name: "subcategoria não encontrada para atualização",
			repoSetup: func() subc.Repository {
				return novoFakeSubcategoriaRepository2() // vazio
			},
			params:  novoAtualizarSubcategoriaParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo subc.Repository, s *subc.Subcategoria) {
				t.Helper()

				// Repositório não deve ter subcategorias salvas
				fake := repo.(*fakeSubcategoriaRepository2)
				if len(fake.subcategorias) != 0 {
					t.Errorf("Esperado 0 subcategorias no repositório, tem %d", len(fake.subcategorias))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			servico := NovoSubcategoriaService(nil, repo)

			s, err := servico.Atualizar(ctx, "subcat-123", tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, s)
			}
		})
	}
}

// TestServicoSubcategoria_BuscarPorID2 testa o método BuscarPorID do serviço de subcategorias.
func TestServicoSubcategoria_BuscarPorID2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() subc.Repository
		id        string
		wantErr   bool
		assertFn  func(t *testing.T, s *subc.Subcategoria)
	}{
		{
			name: "buscar subcategoria com sucesso",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.subcategorias["subcat-123"] = novoSubcategoriaTeste("subcat-123", "Suporte Técnico", "cat-456")
				return repo
			},
			id:      "subcat-123",
			wantErr: false,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s.ID() != "subcat-123" {
					t.Errorf("ID esperado 'subcat-123', recebeu '%s'", s.ID())
				}
				if s.Nome() != "Suporte Técnico" {
					t.Errorf("Nome esperado 'Suporte Técnico', recebeu '%s'", s.Nome())
				}
			},
		},
		{
			name: "erro ao buscar subcategoria",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			id:      "subcat-123",
			wantErr: true,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s != nil {
					t.Errorf("Esperado subcategoria nil, recebeu '%v'", s)
				}
			},
		},
		{
			name: "subcategoria não encontrada",
			repoSetup: func() subc.Repository {
				return novoFakeSubcategoriaRepository2()
			},
			id:      "id-inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s != nil {
					t.Errorf("Esperado subcategoria nil, recebeu '%v'", s)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			servico := NovoSubcategoriaService(nil, repo)

			s, err := servico.BuscarPorID(ctx, tt.id)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, s)
			}
		})
	}
}

// TestServicoSubcategoria_BuscarPorNome2 testa o método BuscarPorNome do serviço de subcategorias.
func TestServicoSubcategoria_BuscarPorNome2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() subc.Repository
		nome      string
		wantErr   bool
		assertFn  func(t *testing.T, s *subc.Subcategoria)
	}{
		{
			name: "buscar subcategoria com sucesso",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.subcategorias["subcat-123"] = novoSubcategoriaTeste("subcat-123", "Suporte Técnico", "cat-456")
				return repo
			},
			nome:    "Suporte Técnico",
			wantErr: false,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s.ID() != "subcat-123" {
					t.Errorf("ID esperado 'subcat-123', recebeu '%s'", s.ID())
				}
				if s.Nome() != "Suporte Técnico" {
					t.Errorf("Nome esperado 'Suporte Técnico', recebeu '%s'", s.Nome())
				}
			},
		},
		{
			name: "erro ao buscar subcategoria",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			nome:    "Suporte Técnico",
			wantErr: true,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s != nil {
					t.Errorf("Esperado subcategoria nil, recebeu '%v'", s)
				}
			},
		},
		{
			name: "subcategoria não encontrada por nome",
			repoSetup: func() subc.Repository {
				return novoFakeSubcategoriaRepository2()
			},
			nome:    "Nome Inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, s *subc.Subcategoria) {
				t.Helper()

				if s != nil {
					t.Errorf("Esperado subcategoria nil, recebeu '%v'", s)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			servico := NovoSubcategoriaService(nil, repo)

			s, err := servico.BuscarPorNome(ctx, tt.nome)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, s)
			}
		})
	}
}

// TestServicoSubcategoria_Listar2 testa o método Listar do serviço de subcategorias.
func TestServicoSubcategoria_Listar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() subc.Repository
		filtro    subc.Filtro
		wantTotal int
		wantErr   bool
		assertFn  func(t *testing.T, subs []subc.Subcategoria, filtroRetornado subc.Filtro)
	}{
		{
			name:   "listar subcategorias com sucesso e normalizar filtro",
			filtro: subc.Filtro{},
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.subcategorias["subcat-1"] = novoSubcategoriaTeste("subcat-1", "Suporte Técnico", "cat-456")
				return repo
			},
			wantTotal: 1,
			wantErr:   false,
			assertFn: func(t *testing.T, subcategorias []subc.Subcategoria, filtroRetornado subc.Filtro) {
				t.Helper()

				if len(subcategorias) != 1 {
					t.Errorf("Esperado 1 subcategoria, recebeu %d", len(subcategorias))
				}

				// filtro deve estar normalizado
				if filtroRetornado.Pagina() != 1 {
					t.Errorf("Página esperada 1, recebeu %d", filtroRetornado.Pagina())
				}

				if filtroRetornado.Limite() != 10 {
					t.Errorf("Limite esperado 10, recebeu %d", filtroRetornado.Limite())
				}
			},
		},
		{
			name: "erro ao listar subcategorias",
			repoSetup: func() subc.Repository {
				repo := novoFakeSubcategoriaRepository2()
				repo.erroAoListar = errors.New("erro no banco")
				return repo
			},
			filtro:  subc.Filtro{},
			wantErr: true,
			assertFn: func(t *testing.T, subcategorias []subc.Subcategoria, filtroRetornado subc.Filtro) {
				t.Helper()

				if len(subcategorias) != 0 {
					t.Errorf("Esperado 0 subcategorias, recebeu %d", len(subcategorias))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			servico := NovoSubcategoriaService(nil, repo)

			subcategorias, total, filtroRetornado, err := servico.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("Total esperado %d, recebeu %d", tt.wantTotal, total)
				}
			}

			if tt.assertFn != nil {
				tt.assertFn(t, subcategorias, filtroRetornado)
			}
		})
	}
}

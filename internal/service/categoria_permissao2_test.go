package service

import (
	"context"
	"errors"
	"testing"

	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

type fakeCategoriaPermissaoRepository struct {
	categoriasPermissoes map[string]*cpm.CategoriaPermissao
	erroAoCriar          error
	erroAoAtualizar      error
	erroAoBuscar         error
	erroAoListar         error
	erroAoDeletar        error
	totalListar          int
}

// novoFakeCategoriaPermissaoRepository cria uma nova instância do repositório fake.
func novoFakeCategoriaPermissaoRepository() *fakeCategoriaPermissaoRepository {
	return &fakeCategoriaPermissaoRepository{
		categoriasPermissoes: make(map[string]*cpm.CategoriaPermissao),
	}
}

// chaveComposta retorna uma chave única baseada em categoriaID e usuarioID.
func chaveComposta(categoriaID, usuarioID string) string {
	return categoriaID + ":" + usuarioID
}

// Criar adiciona uma nova categoria ao repositório fake.
func (f *fakeCategoriaPermissaoRepository) Criar(ctx context.Context, u cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	chave := chaveComposta(u.CategoriaID(), u.UsuarioID())
	f.categoriasPermissoes[chave] = &u
	return &u, nil
}

// BuscarPorID busca uma categoria pelo ID composto no repositório fake.
func (f *fakeCategoriaPermissaoRepository) BuscarPorID(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	chave := chaveComposta(categoriaID, usuarioID)
	cpm, existe := f.categoriasPermissoes[chave]
	if !existe {
		return nil, errors.New("categoria/permissão não encontrada")
	}
	return cpm, nil
}

// Atualizar atualiza uma categoria existente no repositório fake.
func (f *fakeCategoriaPermissaoRepository) Atualizar(ctx context.Context, categoriaID, usuarioID string, c cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	chave := chaveComposta(categoriaID, usuarioID)
	f.categoriasPermissoes[chave] = &c
	return &c, nil
}

// Listar retorna uma lista de categorias do repositório fake.
func (f *fakeCategoriaPermissaoRepository) Listar(ctx context.Context, filtro cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	categoriasPermissoes := make([]cpm.CategoriaPermissao, 0, len(f.categoriasPermissoes))
	for _, cpm := range f.categoriasPermissoes {
		categoriasPermissoes = append(categoriasPermissoes, *cpm)
	}

	total := f.totalListar
	if total == 0 {
		total = len(categoriasPermissoes)
	}

	return categoriasPermissoes, total, nil
}

// Deletar remove uma categoria do repositório fake.
func (f *fakeCategoriaPermissaoRepository) Deletar(ctx context.Context, categoriaID, usuarioID string) error {
	if f.erroAoDeletar != nil {
		return f.erroAoDeletar
	}
	chave := chaveComposta(categoriaID, usuarioID)
	delete(f.categoriasPermissoes, chave)
	return nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

// novoCategoriaPermissaoTest cria uma nova categoria de permissão para testes.
func novoCategoriaPermissaoTeste() *cpm.CategoriaPermissao {
	categoriaPermissao, _ := cpm.Novo(categoriaTesteID, usuarioTesteID, usr.PermADM)
	return categoriaPermissao
}

// novoCriarCategoriaPermissaoParamsTeste cria parâmetros para criar uma categoria de permissão para testes.
func novoCriarCategoriaPermissaoParamsTeste() cpm.CriarParams {
	return cpm.CriarParams{
		CategoriaID: categoriaTesteID,
		UsuarioID:   usuarioTesteID,
		Permissao:   usr.PermADM,
	}
}

// novoAtualizarCategoriaPermissaoParamsTeste cria parâmetros para atualizar uma categoria de permissão para testes.
func novoAtualizarCategoriaPermissaoParamsTeste() cpm.AtualizarParams {
	return cpm.AtualizarParams{
		Permissao: usr.PermADM,
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestCategoriaPermissaoService_Criar testa o método Criar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() cpm.Repository
		params    cpm.CriarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao)
	}{
		{
			name: "criar categoria/permissão com sucesso",
			repoSetup: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			params:  novoCriarCategoriaPermissaoParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao) {
				t.Helper()

				if cpm.CategoriaID() != "ctg-123" {
					t.Errorf("esperado CategoriaID 'ctg-123', obtido '%s'", cpm.CategoriaID())
				}
				if cpm.UsuarioID() != "usr-123" {
					t.Errorf("esperado UsuarioID 'usr-123', obtido '%s'", cpm.UsuarioID())
				}
				if cpm.Permissao() != usr.PermADM.String() {
					t.Errorf("esperado Permissao 'PermADM', obtido '%s'", cpm.Permissao())
				}

				// Verificar se foi salvo no repositório
				fake := repo.(*fakeCategoriaPermissaoRepository)
				chave := chaveComposta("ctg-123", "usr-123")
				salvo, existe := fake.categoriasPermissoes[chave]
				if !existe {
					t.Errorf("categoria/permissão não encontrada no repositório")
				}
				if salvo.CategoriaID() != "ctg-123" || salvo.UsuarioID() != "usr-123" || salvo.Permissao() != usr.PermADM.String() {
					t.Errorf("dados salvos no repositório não correspondem aos esperados")
				}
			},
		},
		{
			name: "erro ao criar categoria/permissão devido a erro no repositório",
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.erroAoCriar = errors.New("erro ao criar no repositório")
				return repo
			},
			params:  novoCriarCategoriaPermissaoParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao) {
				t.Helper()
				// repositório deve estar vazio
				fake := repo.(*fakeCategoriaPermissaoRepository)
				if len(fake.categoriasPermissoes) != 0 {
					t.Errorf("esperado repositório vazio, mas contém dados")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaPermissaoService(repo)

			cpm, err := service.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, cpm)
			}
		})
	}
}

// TestCategoriaPermissaoService_BuscarPorID testa o método BuscarPorID do serviço de categoria de permissão.
func TestCategoriaPermissaoService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		ctgID     string
		usrID     string
		repoSetup func() cpm.Repository
		wantErr   bool
		assertFn  func(t *testing.T, cpm *cpm.CategoriaPermissao)
	}{
		{
			name:  "buscar categoria/permissão com sucesso",
			ctgID: categoriaTesteID,
			usrID: usuarioTesteID,
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				return repo
			},
			wantErr: false,
			assertFn: func(t *testing.T, cpm *cpm.CategoriaPermissao) {
				t.Helper()

				if cpm.CategoriaID() != categoriaTesteID {
					t.Errorf("esperado CategoriaID '%s', obtido '%s'", categoriaTesteID, cpm.CategoriaID())
				}
				if cpm.UsuarioID() != usuarioTesteID {
					t.Errorf("esperado UsuarioID '%s', obtido '%s'", usuarioTesteID, cpm.UsuarioID())
				}
				if cpm.Permissao() != usr.PermADM.String() {
					t.Errorf("esperado Permissao 'PermADM', obtido '%s'", cpm.Permissao())
				}
			},
		},
		{
			name:  "erro ao buscar categoria/permissão não existente",
			ctgID: "id-inexistente",
			usrID: "id-inexistente",
			repoSetup: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			wantErr: true,
			assertFn: func(t *testing.T, cpm *cpm.CategoriaPermissao) {
				t.Helper()
				if cpm != nil {
					t.Errorf("esperado categoria/permissão nula, mas obteve valor")
				}
			},
		},
		{
			name:  "erro ao buscar categoria/permissão devido a erro no repositório",
			ctgID: categoriaTesteID,
			usrID: usuarioTesteID,
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.erroAoBuscar = errors.New("erro ao buscar no repositório")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, cpm *cpm.CategoriaPermissao) {
				t.Helper()
				if cpm != nil {
					t.Errorf("esperado categoria/permissão nula, mas obteve valor")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaPermissaoService(repo)

			cpm, err := service.BuscarPorID(ctx, tt.ctgID, tt.usrID)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, cpm)
			}
		})
	}
}

// TestCategoriaPermissaoService_Atualizar testa o método Atualizar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() cpm.Repository
		params    cpm.AtualizarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao)
	}{
		{
			name: "atualizar categoria/permissão com sucesso",
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				return repo
			},
			params:  novoAtualizarCategoriaPermissaoParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao) {
				t.Helper()

				if cpm.Permissao() != usr.PermADM.String() {
					t.Errorf("esperado Permissao 'PermADM', obtido '%s'", cpm.Permissao())
				}

				// Verificar se foi atualizado no repositório
				fake := repo.(*fakeCategoriaPermissaoRepository)
				chave := chaveComposta(categoriaTesteID, usuarioTesteID)
				atualizado, existe := fake.categoriasPermissoes[chave]
				if !existe {
					t.Errorf("categoria/permissão não encontrada no repositório")
					return
				}
				if atualizado.Permissao() != usr.PermADM.String() {
					t.Errorf("dados atualizados no repositório não correspondem aos esperados")
				}
			},
		},
		{
			name: "erro ao atualizar categoria/permissão não existente",
			repoSetup: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			params:  novoAtualizarCategoriaPermissaoParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao) {
				t.Helper()
				if cpm != nil {
					t.Errorf("esperado categoria/permissão nula, mas obteve valor")
				}
			},
		},
		{
			name: "erro ao atualizar categoria/permissão devido a erro no repositório",
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				repo.erroAoAtualizar = errors.New("erro ao atualizar no repositório")
				return repo
			},
			params:  novoAtualizarCategoriaPermissaoParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo cpm.Repository, cpm *cpm.CategoriaPermissao) {
				t.Helper()
				if cpm != nil {
					t.Errorf("esperado categoria/permissão nula, mas obteve valor")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaPermissaoService(repo)

			cpm, err := service.Atualizar(ctx, categoriaTesteID, usuarioTesteID, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, cpm)
			}
		})
	}
}

// TestCategoriaPermissaoService_Deletar testa o método Deletar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Deletar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		ctgID     string
		usrID     string
		repoSetup func() cpm.Repository
		wantErr   bool
		assertFn  func(t *testing.T, repo cpm.Repository)
	}{
		{
			name:  "deletar categoria/permissão com sucesso",
			ctgID: categoriaTesteID,
			usrID: usuarioTesteID,
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				return repo
			},
			wantErr: false,
			assertFn: func(t *testing.T, repo cpm.Repository) {
				t.Helper()

				// Verificar se foi removido do repositório
				fake := repo.(*fakeCategoriaPermissaoRepository)
				chave := chaveComposta(categoriaTesteID, usuarioTesteID)
				_, existe := fake.categoriasPermissoes[chave]
				if existe {
					t.Errorf("categoria/permissão ainda existe no repositório após deleção")
				}
			},
		},
		{
			name:  "erro ao deletar categoria/permissão não existente",
			ctgID: "id-inexistente",
			usrID: "id-inexistente",
			repoSetup: func() cpm.Repository {
				return novoFakeCategoriaPermissaoRepository()
			},
			wantErr: true,
			assertFn: func(t *testing.T, repo cpm.Repository) {
				t.Helper()
				// repositório deve estar vazio
				fake := repo.(*fakeCategoriaPermissaoRepository)
				if len(fake.categoriasPermissoes) != 0 {
					t.Errorf("esperado repositório vazio, mas contém dados")
				}
			},
		},
		{
			name:  "erro ao deletar categoria/permissão devido a erro no repositório",
			ctgID: categoriaTesteID,
			usrID: usuarioTesteID,
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				repo.erroAoDeletar = errors.New("erro ao deletar no repositório")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, repo cpm.Repository) {
				t.Helper()
				// categoria/permissão deve continuar existindo no repositório
				fake := repo.(*fakeCategoriaPermissaoRepository)
				chave := chaveComposta(categoriaTesteID, usuarioTesteID)
				_, existe := fake.categoriasPermissoes[chave]
				if !existe {
					t.Errorf("categoria/permissão não encontrada no repositório após falha na deleção")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaPermissaoService(repo)

			err := service.Deletar(ctx, tt.ctgID, tt.usrID)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo)
			}
		})
	}
}

// TestCategoriaPermissaoService_Listar testa o método Listar do serviço de categoria de permissão.
func TestCategoriaPermissaoService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() cpm.Repository
		filtro    cpm.Filtro
		wantTotal int
		wantErr   bool
		assertFn  func(t *testing.T, categoriasPermissoes []cpm.CategoriaPermissao, filtro cpm.Filtro)
	}{
		{
			name: "listar categorias/permissões com sucesso",
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.categoriasPermissoes[chaveComposta(categoriaTesteID, usuarioTesteID)] = novoCategoriaPermissaoTeste()
				return repo
			},
			filtro:    cpm.Filtro{},
			wantTotal: 2,
			wantErr:   false,
			assertFn: func(t *testing.T, categoriasPermissoes []cpm.CategoriaPermissao, filtro cpm.Filtro) {
				t.Helper()

				if len(categoriasPermissoes) != 2 {
					t.Errorf("esperado 2 categorias/permissões, obtido %d", len(categoriasPermissoes))
				}

				// filtro deve estar normalizado
				if filtro.Limite() <= 0 {
					t.Errorf("esperado filtro com Limite positivo, obtido %d", filtro.Limite())
				}
				if filtro.Pagina() < 0 {
					t.Errorf("esperado filtro com Página não negativa, obtido %d", filtro.Pagina())
				}
			},
		},
		{
			name: "erro ao listar categorias/permissões devido a erro no repositório",
			repoSetup: func() cpm.Repository {
				repo := novoFakeCategoriaPermissaoRepository()
				repo.erroAoListar = errors.New("erro ao listar no repositório")
				return repo
			},
			filtro:    cpm.Filtro{},
			wantTotal: 0,
			wantErr:   true,
			assertFn: func(t *testing.T, categoriasPermissoes []cpm.CategoriaPermissao, filtro cpm.Filtro) {
				t.Helper()
				if categoriasPermissoes != nil {
					t.Errorf("esperado categorias/permissões nulas, mas obteve valor")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoCategoriaPermissaoService(repo)

			categoriasPermissoes, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("esperado total %d, obtido %d", tt.wantTotal, total)
				}
			}

			if tt.assertFn != nil {
				tt.assertFn(t, categoriasPermissoes, filtroRetornado)
			}
		})
	}
}

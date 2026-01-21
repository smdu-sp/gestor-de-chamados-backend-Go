package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeUsuarioRepository é uma implementação falsa de usr.Repository para testes.
type fakeUsuarioRepository struct {
	criarFn          func(ctx context.Context, u usr.Usuario) (*usr.Usuario, error)
	atualizarFn      func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error)
	buscarPorIDFn    func(ctx context.Context, id string) (*usr.Usuario, error)
	buscarPorLoginFn func(ctx context.Context, login string) (*usr.Usuario, error)
	listarFn         func(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, error)
}

// Criar chama a função criarFn configurada no fakeUsuarioRepository.
func (f *fakeUsuarioRepository) Criar(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
	return f.criarFn(ctx, u)
}

// Atualizar chama a função atualizarFn configurada no fakeUsuarioRepository.
func (f *fakeUsuarioRepository) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	return f.atualizarFn(ctx, id, u)
}

// BuscarPorID chama a função buscarPorIDFn configurada no fakeUsuarioRepository.
func (f *fakeUsuarioRepository) BuscarPorID(ctx context.Context, id string) (*usr.Usuario, error) {
	return f.buscarPorIDFn(ctx, id)
}

// BuscarPorLogin chama a função buscarPorLoginFn configurada no fakeUsuarioRepository.
func (f *fakeUsuarioRepository) BuscarPorLogin(ctx context.Context, login string) (*usr.Usuario, error) {
	return f.buscarPorLoginFn(ctx, login)
}

// Listar chama a função listarFn configurada no fakeUsuarioRepository.
func (f *fakeUsuarioRepository) Listar(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestUsuarioService_Criar testa o método Criar do UsuarioService.
func TestUsuarioService_Criar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repo      usr.Repository
		wantErr   bool
	}{
		{
			name: "criar usuario com sucesso",
			geradorID: &fakeGeradorID{
				id: "user-123",
			},
			repo: &fakeUsuarioRepository{
				criarFn: func(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
					if u.ID() != "user-123" {
						t.Fatalf("id inesperado: %s", u.ID())
					}
					if u.Login() != "rogerio" {
						t.Fatalf("login inesperado: %s", u.Login())
					}
					return &u, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errors.New("erro ao gerar id"),
			},
			repo: &fakeUsuarioRepository{
				criarFn: func(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
					t.Fatalf("Criar não deveria ser chamado quando gerador de ID falha")
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			geradorID: &fakeGeradorID{
				id: "user-123",
			},
			repo: &fakeUsuarioRepository{
				criarFn: func(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(tt.geradorID, tt.repo)

			params := usr.CriarParams{
				Nome:      "Rogério",
				Login:     "rogerio",
				Email:     usr.NovoEmail("rogerio@email.com"),
				Permissao: usr.PermADM,
			}

			// Act
			_, err := service.Criar(ctx, params)

			// Assert
			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestUsuarioService_Atualizar testa o método Atualizar do UsuarioService.
func TestUsuarioService_Atualizar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    usr.Repository
		params  usr.AtualizarParams
		wantErr bool
	}{
		{
			name: "atualizar usuario com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					u, _ := usr.Novo(
						"user-123",
						"Rogério",
						"rogerio",
						usr.NovoEmail("rogerio@email.com"),
						usr.PermADM,
						nil,
					)
					return u, nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					// Assert dentro do fake (idiomático)
					if u.Nome() != "Novo Nome" {
						t.Fatalf("esperava nome 'Novo Nome', recebeu '%s'", u.Nome())
					}
					return &u, nil
				},
			},
			params: usr.AtualizarParams{
				Nome: ptr("Novo Nome"),
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					t.Fatalf("Atualizar não deveria ser chamado")
					return nil, nil
				},
			},
			params:  usr.AtualizarParams{},
			wantErr: true,
		},
		{
			name: "erro ao salvar atualizacao",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					u, _ := usr.Novo(
						"user-123",
						"Rogério",
						"rogerio",
						usr.NovoEmail("rogerio@email.com"),
						usr.PermADM,
						nil,
					)
					return u, nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro ao atualizar")
				},
			},
			params:  usr.AtualizarParams{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			_, err := service.Atualizar(ctx, "user-123", tt.params)

			// Assert
			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestUsuarioService_BuscarPorID testa o método BuscarPorID do UsuarioService.
func TestUsuarioService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	usuarioOK := &usr.Usuario{}

	tests := []struct {
		name    string
		repo    usr.Repository
		want    *usr.Usuario
		wantErr bool
	}{
		{
			name: "buscar usuario com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					if id != "user-123" {
						t.Fatalf("esperava id 'user-123', recebeu '%s'", id)
					}
					return usuarioOK, nil
				},
			},
			want:    usuarioOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "usuario não encontrado",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("usuario não encontrado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			got, err := service.BuscarPorID(ctx, "user-123")

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if got.ID() != tt.want.ID() {
				t.Fatalf("esperava id '%s', recebeu '%s'", tt.want.ID(), got.ID())
			}
		})
	}
}

// TestUsuarioService_BuscarPorLogin testa o método BuscarPorLogin do UsuarioService.
func TestUsuarioService_BuscarPorLogin(t *testing.T) {
	ctx := context.Background()

	usuarioOK := &usr.Usuario{}

	tests := []struct {
		name    string
		repo    usr.Repository
		want    *usr.Usuario
		wantErr bool
	}{
		{
			name: "buscar usuario por login com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorLoginFn: func(ctx context.Context, login string) (*usr.Usuario, error) {
					if login != "rogerio" {
						t.Fatalf("esperava login 'rogerio', recebeu '%s'", login)
					}
					return usuarioOK, nil
				},
			},
			want:    usuarioOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario por login",
			repo: &fakeUsuarioRepository{
				buscarPorLoginFn: func(ctx context.Context, login string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			got, err := service.BuscarPorLogin(ctx, "rogerio")

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if got.Login() != tt.want.Login() {
				t.Fatalf("esperava login '%s', recebeu '%s'", tt.want.Login(), got.Login())
			}
		})
	}
}

// TestUsuarioService_AtualizarPermissao testa o método AtualizarPermissao do UsuarioService.
func TestUsuarioService_AtualizarPermissao(t *testing.T) {
	ctx := context.Background()

	usuarioBase := func() *usr.Usuario {
		u, _ := usr.Novo(
			"user-123",
			"Rogério",
			"rogerio",
			usr.NovoEmail("rogerio@email.com"),
			usr.PermADM,
			nil,
		)
		return u
	}

	tests := []struct {
		name    string
		repo    usr.Repository
		params  usr.AtualizarPermissaoParams
		wantErr bool
	}{
		{
			name: "atualizar permissao com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioBase(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return &u, nil
				},
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: true,
		},
		{
			name: "erro de validacao ao atualizar permissao",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioBase(), nil
				},
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: "INVALIDA", // força erro de domínio
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioBase(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			u, err := service.AtualizarPermissao(ctx, "user-123", tt.params)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if u.Permissao() != tt.params.Permissao {
				t.Fatalf(
					"esperava permissao '%s', recebeu '%s'",
					tt.params.Permissao,
					u.Permissao(),
				)
			}
		})
	}
}

// TestUsuarioService_AtualizarUltimoLogin testa o método AtualizarUltimoLogin do UsuarioService.
func TestUsuarioService_AtualizarUltimoLogin(t *testing.T) {
	ctx := context.Background()

	usuarioBase := func() *usr.Usuario {
		u, _ := usr.Novo(
			"user-123",
			"Rogério",
			"rogerio",
			usr.NovoEmail("rogerio@email.com"),
			usr.PermADM,
			nil,
		)
		return u
	}

	tests := []struct {
		name    string
		repo    usr.Repository
		wantErr bool
	}{
		{
			name: "atualizar ultimo login com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioBase(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return &u, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioBase(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			u, err := service.AtualizarUltimoLogin(ctx, "user-123")

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if u == nil {
				t.Fatalf("esperava usuario, recebeu nil")
			}
		})
	}
}

// TestUsuarioService_Desativar testa o método Desativar do UsuarioService.
func TestUsuarioService_Desativar(t *testing.T) {
	ctx := context.Background()

	usuarioAtivo := func() *usr.Usuario {
		u, _ := usr.Novo(
			"user-123",
			"Rogério",
			"rogerio",
			usr.NovoEmail("rogerio@email.com"),
			usr.PermADM,
			nil,
		)
		return u
	}

	tests := []struct {
		name    string
		repo    usr.Repository
		wantErr bool
	}{
		{
			name: "desativar usuario com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioAtivo(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return &u, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioAtivo(), nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			u, err := service.Desativar(ctx, "user-123")

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if u.Status() != false {
				t.Fatalf("esperava usuario desativado")
			}
		})
	}
}

// TestUsuarioService_Ativar testa o método Ativar do UsuarioService.
func TestUsuarioService_Ativar(t *testing.T) {
	ctx := context.Background()

	usuarioDesativado, _ := usr.Novo(
		"user-123",
		"Rogério",
		"rogerio",
		usr.NovoEmail("rogerio@email.com"),
		usr.PermADM,
		nil,
	)

	tests := []struct {
		name    string
		repo    usr.Repository
		wantErr bool
	}{
		{
			name: "ativar usuario com sucesso",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioDesativado, nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return &u, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioDesativado, nil
				},
				atualizarFn: func(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			u, err := service.Ativar(ctx, "user-123")

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if u.Status() != true {
				t.Fatalf("esperava usuario ativado")
			}
		})
	}
}

// TestUsuarioService_Listar testa o método Listar do UsuarioService.
func TestUsuarioService_Listar(t *testing.T) {
	ctx := context.Background()

	u, _ := usr.Novo(
		"user-1",
		"Rogério",
		"rogerio",
		usr.NovoEmail("rogerio@email.com"),
		usr.PermADM,
		nil,
	)
	usuarios := []usr.Usuario{*u}

	tests := []struct {
		name      string
		repo      usr.Repository
		filtro    usr.Filtro
		wantTotal int
		wantErr   bool
	}{
		{
			name:   "listar usuarios com sucesso e normalizar filtro",
			filtro: usr.Filtro{},
			repo: &fakeUsuarioRepository{
				listarFn: func(ctx context.Context, f usr.Filtro) ([]usr.Usuario, int, error) {
					// Assert indireto: o filtro chegou normalizado no repo
					if f.Pagina() <= 0 {
						t.Fatalf("esperava pagina normalizada, recebeu %d", f.Pagina())
					}
					if f.Limite() <= 0 {
						t.Fatalf("esperava limite normalizado, recebeu %d", f.Limite())
					}
					return usuarios, len(usuarios), nil
				},
			},
			wantTotal: len(usuarios),
			wantErr:   false,
		},
		{
			name:   "erro ao listar usuarios",
			filtro: usr.Filtro{},
			repo: &fakeUsuarioRepository{
				listarFn: func(ctx context.Context, f usr.Filtro) ([]usr.Usuario, int, error) {
					return nil, 0, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			us, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if total != tt.wantTotal {
				t.Fatalf("esperava total %d, recebeu %d", tt.wantTotal, total)
			}

			if len(us) != tt.wantTotal {
				t.Fatalf("esperava %d usuarios, recebeu %d", tt.wantTotal, len(us))
			}

			// Assert importante: filtro retornado é o normalizado
			if filtroRetornado.Pagina() <= 0 {
				t.Fatalf("esperava filtro pagina normalizada")
			}
			if filtroRetornado.Limite() <= 0 {
				t.Fatalf("esperava filtro limite normalizado")
			}
		})
	}
}

// TestUsuarioService_VerificarPermissao testa o método VerificarPermissao do UsuarioService.
func TestUsuarioService_VerificarPermissao(t *testing.T) {
	ctx := context.Background()

	usuarioADM := func() *usr.Usuario {
		u, _ := usr.Novo(
			"user-123",
			"Rogério",
			"rogerio",
			usr.NovoEmail("rogerio@email.com"),
			usr.PermADM,
			nil,
		)
		return u
	}

	tests := []struct {
		name       string
		repo       usr.Repository
		permissoes []usr.Permissao
		want       bool
		wantErr    bool
	}{
		{
			name: "usuario possui permissao",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioADM(), nil
				},
			},
			permissoes: []usr.Permissao{usr.PermUSR, usr.PermADM},
			want:       true,
			wantErr:    false,
		},
		{
			name: "usuario nao possui permissao",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usuarioADM(), nil
				},
			},
			permissoes: []usr.Permissao{usr.PermUSR},
			want:       false,
			wantErr:    false,
		},
		{
			name: "erro ao buscar usuario",
			repo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return nil, errors.New("erro no banco")
				},
			},
			permissoes: []usr.Permissao{usr.PermADM},
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoUsuarioService(nil, tt.repo)

			// Act
			ok, err := service.VerificarPermissao(ctx, "user-123", tt.permissoes)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if ok != tt.want {
				t.Fatalf("esperava %v, recebeu %v", tt.want, ok)
			}
		})
	}
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// ptr é um helper que retorna um ponteiro para o valor fornecido.
func ptr[T any](v T) *T {
	return &v
}

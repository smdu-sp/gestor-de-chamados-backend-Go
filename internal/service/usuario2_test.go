package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeUsuarioRepository é uma implementação falsa de usr.Repository para testes.
// Usa estado em vez de funções injetadas, tornando os testes mais legíveis.
type fakeUsuarioRepository2 struct {
	usuarios        map[string]*usr.Usuario
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
	totalListar     int
}

// novoFakeUsuarioRepository cria uma nova instância de fakeUsuarioRepository.
func novoFakeUsuarioRepository2() *fakeUsuarioRepository2 {
	return &fakeUsuarioRepository2{
		usuarios: make(map[string]*usr.Usuario),
	}
}

// Criar salva um novo usuário no repositório falso.
func (f *fakeUsuarioRepository2) Criar(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	f.usuarios[u.ID()] = &u
	return &u, nil
}

// Atualizar atualiza um usuário existente no repositório falso.
func (f *fakeUsuarioRepository2) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoAtualizar != nil { // reusa erro genérico
		return nil, f.erroAoAtualizar
	}
	f.usuarios[id] = &u
	return &u, nil
}

// BuscarPorID busca um usuário por ID no repositório falso.
func (f *fakeUsuarioRepository2) BuscarPorID(ctx context.Context, id string) (*usr.Usuario, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	u, existe := f.usuarios[id]
	if !existe {
		return nil, errors.New("usuário não encontrado")
	}
	return u, nil
}

// BuscarPorLogin busca um usuário por login no repositório falso.
func (f *fakeUsuarioRepository2) BuscarPorLogin(ctx context.Context, login string) (*usr.Usuario, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	for _, u := range f.usuarios {
		if u.Login() == login {
			return u, nil
		}
	}
	return nil, errors.New("usuário não encontrado")
}

// Listar lista usuários com base no filtro no repositório falso.
func (f *fakeUsuarioRepository2) Listar(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	usuarios := make([]usr.Usuario, 0, len(f.usuarios))
	for _, u := range f.usuarios {
		usuarios = append(usuarios, *u)
	}

	total := f.totalListar
	if total == 0 {
		total = len(usuarios)
	}

	return usuarios, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

// novoUsuarioTeste cria um usuário de teste com os dados fornecidos.
func novoUsuarioTeste(id, nome, login string) *usr.Usuario {
	u, _ := usr.Novo(
		id,
		nome,
		login,
		usr.NovoEmail(login+"@email.com"),
		usr.PermADM,
		nil,
	)
	return u
}

// novoCriarUsuarioParamsTeste cria parâmetros de criação de usuário para testes.
func novoCriarUsuarioParamsTeste() usr.CriarParams {
	return usr.CriarParams{
		Nome:      "Rogério",
		Login:     "rogerio",
		Email:     usr.NovoEmail("rogerio@email.com"),
		Permissao: usr.PermADM,
		Avatar:    nil,
	}
}

// novoAtualizarUsuarioParamsTeste cria parâmetros de atualização de usuário para testes.
func novoAtualizarUsuarioParamsTeste() usr.AtualizarParams {
	return usr.AtualizarParams{
		Nome:      ptr("Rogério Atualizado"),
		Login:     ptr("rogerioatualizado"),
		Email:     ptr(usr.NovoEmail("rogerioatualizado@email.com")),
		Permissao: ptr(usr.PermADM),
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

func TestUsuarioService_Criar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repoSetup func() usr.Repository
		params    usr.CriarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo usr.Repository, u *usr.Usuario)
	}{
		{
			name: "criar usuario com sucesso",
			geradorID: &fakeGeradorID{
				id: "user-123",
			},
			repoSetup: func() usr.Repository {
				return novoFakeUsuarioRepository2()
			},
			params:  novoCriarUsuarioParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				if u.ID() != "user-123" {
					t.Errorf("esperava ID 'user-123', recebeu '%s'", u.ID())
				}
				if u.Login() != "rogerio" {
					t.Errorf("esperava login 'rogerio', recebeu '%s'", u.Login())
				}

				// Verifica se foi salvo no repositório
				fake := repo.(*fakeUsuarioRepository2)
				if len(fake.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fake.usuarios))
				}
			},
		},
		{
			name: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errors.New("erro ao gerar id"),
			},
			repoSetup: func() usr.Repository {
				return novoFakeUsuarioRepository2()
			},
			params:  novoCriarUsuarioParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				// Repositório não deve ter usuários salvos
				fake := repo.(*fakeUsuarioRepository2)
				if len(fake.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando gerador de ID falha")
				}
			},
		},
		{
			name: "erro ao salvar no repositorio",
			geradorID: &fakeGeradorID{
				id: "user-123",
			},
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoCriar = errors.New("erro no banco")
				return repo
			},
			params:  novoCriarUsuarioParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				// Repositório não deve ter usuários salvos
				fake := repo.(*fakeUsuarioRepository2)
				if len(fake.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando salvar falha")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(tt.geradorID, repo)

			u, err := service.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, u)
			}
		})
	}
}

// TestUsuarioService_Atualizar2 testa o método Atualizar do serviço de usuário.
func TestUsuarioService_Atualizar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		params    usr.AtualizarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo usr.Repository, u *usr.Usuario)
	}{
		{
			name: "atualizar usuario com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			params:  novoAtualizarUsuarioParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()
				if u.Nome() != "Rogério Atualizado" {
					t.Errorf("esperava nome 'Rogério Atualizado', recebeu '%s'", u.Nome())
				}

				// Repositório deve ter o usuário atualizado
				fake := repo.(*fakeUsuarioRepository2)
				updated, exists := fake.usuarios["user-123"]
				if !exists {
					t.Error("usuário atualizado não encontrado no repositório")
					return
				}
				if updated.Nome() != "Rogério Atualizado" {
					t.Errorf("Nome esperado 'Rogério Atualizado', recebido '%s'", updated.Nome())
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			params:  novoAtualizarUsuarioParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				// Repositório não deve ter usuários salvos
				fake := repo.(*fakeUsuarioRepository2)
				if len(fake.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando buscar falha")
				}
			},
		},
		{
			name: "erro ao salvar atualizacao",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				repo.erroAoAtualizar = errors.New("erro ao atualizar")
				return repo
			},
			params:  novoAtualizarUsuarioParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				// Repositório deve ter o usuário original sem alterações
				fake := repo.(*fakeUsuarioRepository2)
				original, exists := fake.usuarios["user-123"]
				if !exists {
					t.Error("usuário atualizado não encontrado no repositório")
					return
				}
				if original.Nome() != "Rogério" {
					t.Errorf("Nome esperado 'Rogério', recebido '%s'", original.Nome())
				}
			},
		},
		{
			name: "usuario nao encontrado para atualizar",
			repoSetup: func() usr.Repository {
				return novoFakeUsuarioRepository2() // vazio
			},
			params:  novoAtualizarUsuarioParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo usr.Repository, u *usr.Usuario) {
				t.Helper()

				// Repositório não deve ter usuários salvos
				fake := repo.(*fakeUsuarioRepository2)
				if len(fake.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando usuário não encontrado")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.Atualizar(ctx, "user-123", tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, u)
			}
		})
	}
}

func TestUsuarioService_BuscarPorID2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		id        string
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "buscar usuario com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			id:      "user-123",
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u.ID() != "user-123" {
					t.Errorf("esperava ID 'user-123', recebeu '%s'", u.ID())
				}
				if u.Login() != "rogerio" {
					t.Errorf("esperava login 'rogerio', recebeu '%s'", u.Login())
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			id:      "user-123",
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "usuario não encontrado",
			repoSetup: func() usr.Repository {
				return novoFakeUsuarioRepository2() // vazio
			},
			id:      "user-inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.BuscarPorID(ctx, tt.id)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_BuscarPorLogin2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		login     string
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "buscar usuario por login com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			login:   "rogerio",
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u.Login() != "rogerio" {
					t.Errorf("esperava login 'rogerio', recebeu '%s'", u.Login())
				}
			},
		},
		{
			name: "erro ao buscar usuario por login",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			login:   "rogerio",
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "usuario por login nao encontrado",
			repoSetup: func() usr.Repository {
				return novoFakeUsuarioRepository2() // vazio
			},
			login:   "login-inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.BuscarPorLogin(ctx, tt.login)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_AtualizarPermissao2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		params    usr.AtualizarPermissaoParams
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "atualizar permissao com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u.Permissao() != usr.PermUSR {
					t.Errorf("esperava permissão USR, recebeu '%s'", u.Permissao())
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "erro de validacao ao atualizar permissao",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: "INVALIDA",
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "erro ao salvar no repositorio",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			params: usr.AtualizarPermissaoParams{
				Permissao: usr.PermUSR,
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.AtualizarPermissao(ctx, "user-123", tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_AtualizarUltimoLogin2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "atualizar ultimo login com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u == nil {
					t.Fatal("esperava usuário, recebeu nil")
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.AtualizarUltimoLogin(ctx, "user-123")

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_Desativar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "desativar usuario com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u.Status() != false {
					t.Error("esperava usuário desativado")
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "erro ao salvar no repositorio",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.Desativar(ctx, "user-123")

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_Ativar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		wantErr   bool
		assertFn  func(t *testing.T, u *usr.Usuario)
	}{
		{
			name: "ativar usuario com sucesso",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			wantErr: false,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()
				if u.Status() != true {
					t.Error("esperava usuário ativado")
				}
			},
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
		{
			name: "erro ao salvar no repositorio",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				repo.erroAoAtualizar = errors.New("erro no banco")
				return repo
			},
			wantErr: true,
			assertFn: func(t *testing.T, u *usr.Usuario) {
				t.Helper()

				if u != nil {
					t.Errorf("esperava nil, recebeu '%v'", u)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			u, err := service.Ativar(ctx, "user-123")

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, u)
			}
		})
	}
}

func TestUsuarioService_Listar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() usr.Repository
		filtro    usr.Filtro
		wantTotal int
		wantErr   bool
		assertFn  func(t *testing.T, usuarios []usr.Usuario, filtroRetornado usr.Filtro)
	}{
		{
			name:   "listar usuarios com sucesso e normalizar filtro",
			filtro: usr.Filtro{},
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-1"] = novoUsuarioTeste("user-1", "Rogério", "rogerio")
				return repo
			},
			wantTotal: 1,
			wantErr:   false,
			assertFn: func(t *testing.T, usuarios []usr.Usuario, filtroRetornado usr.Filtro) {
				t.Helper()

				if len(usuarios) != 1 {
					t.Errorf("esperava 1 usuário, recebeu %d", len(usuarios))
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
			name:   "erro ao listar usuarios",
			filtro: usr.Filtro{},
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoListar = errors.New("erro no banco")
				return repo
			},
			wantTotal: 0,
			wantErr: true,
			assertFn: func(t *testing.T, usuarios []usr.Usuario, filtroRetornado usr.Filtro) {
				t.Helper()
				if len(usuarios) != 0 {
					t.Errorf("esperava 0 usuários, recebeu %d", len(usuarios))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			usuarios, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("esperava total %d, recebeu %d", tt.wantTotal, total)
				}

				if tt.assertFn != nil {
					tt.assertFn(t, usuarios, filtroRetornado)
				}
			}
		})
	}
}

func TestUsuarioService_VerificarPermissao2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		repoSetup  func() usr.Repository
		permissoes []usr.Permissao
		want       bool
		wantErr    bool
	}{
		{
			name: "usuario possui permissao",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			permissoes: []usr.Permissao{usr.PermUSR, usr.PermADM},
			want:       true,
			wantErr:    false,
		},
		{
			name: "usuario nao possui permissao",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.usuarios["user-123"] = novoUsuarioTeste("user-123", "Rogério", "rogerio")
				return repo
			},
			permissoes: []usr.Permissao{usr.PermUSR},
			want:       false,
			wantErr:    false,
		},
		{
			name: "erro ao buscar usuario",
			repoSetup: func() usr.Repository {
				repo := novoFakeUsuarioRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			permissoes: []usr.Permissao{usr.PermADM},
			want:       false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoUsuarioService(nil, repo)

			ok, err := service.VerificarPermissao(ctx, "user-123", tt.permissoes)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr && ok != tt.want {
				t.Errorf("esperava %v, recebeu %v", tt.want, ok)
			}
		})
	}
}

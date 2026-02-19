package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// asserção de interface para garantir que fakeUsuarioRepository implementa usr.Repository
var _ usr.Repository = (*fakeUsuarioRepository)(nil)

// fakeUsuarioRepository é uma implementação falsa do repositório de usuários para testes.
type fakeUsuarioRepository struct {
	usuarios        map[string]*usr.Usuario
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
}

// novoFakeUsuarioRepository cria uma nova instância de fakeUsuarioRepository.
func novoFakeUsuarioRepository() *fakeUsuarioRepository {
	return &fakeUsuarioRepository{
		usuarios: make(map[string]*usr.Usuario),
	}
}

// Criar salva um novo usuário no repositório falso.
func (f *fakeUsuarioRepository) Criar(ctx context.Context, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	usuario := u
	f.usuarios[u.ID()] = &usuario
	return &usuario, nil
}

// Atualizar atualiza um usuário existente no repositório falso.
func (f *fakeUsuarioRepository) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	usuario := u
	f.usuarios[id] = &usuario
	return &usuario, nil
}

// BuscarPorID busca um usuário por ID no repositório falso.
func (f *fakeUsuarioRepository) BuscarPorID(ctx context.Context, id string) (*usr.Usuario, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	u, existe := f.usuarios[id]
	if !existe {
		return nil, mysql.ErrUsuarioNaoEncontrado
	}
	return u, nil
}

// BuscarPorLogin busca um usuário por login no repositório falso.
func (f *fakeUsuarioRepository) BuscarPorLogin(ctx context.Context, login string) (*usr.Usuario, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	for _, u := range f.usuarios {
		if u.Login() == login {
			return u, nil
		}
	}
	return nil, mysql.ErrUsuarioNaoEncontrado
}

// Listar lista usuários com base no filtro no repositório falso.
func (f *fakeUsuarioRepository) Listar(ctx context.Context, filtro usr.Filtro) ([]usr.Usuario, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para uma slice
	usuarios := make([]usr.Usuario, 0, len(f.usuarios))
	for _, u := range f.usuarios {
		usuarios = append(usuarios, *u)
	}

	// Ordenação ORDER BY nome ASC
	sort.Slice(usuarios, func(i, j int) bool {
		return usuarios[i].Nome() < usuarios[j].Nome()
	})

	// Contar o total antes da paginação
	total := len(usuarios)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio >= total {
			return []usr.Usuario{}, total, nil
		}

		if fim > total {
			fim = total
		}

		usuarios = usuarios[inicio:fim]
	}

	return usuarios, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const (
	usuarioTesteID    = "usuario-123"
	usuarioTesteLogin = "rogerio"
)

// novoUsuarioTeste cria um usuário de teste com valores fixos.
func novoUsuarioTeste() *usr.Usuario {
	u, _ := usr.Novo(
		usuarioTesteID,
		"Rogério",
		usuarioTesteLogin,
		usr.NovoEmail("rogerio@email.com"),
		usr.PermADM,
		nil,
	)
	return u
}

// popularUsuariosTeste adiciona múltiplos usuários ao repositório falso para testes.
func popularUsuariosTeste(repo *fakeUsuarioRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("usuario-%02d", i)

		u, _ := usr.Novo(
			id,
			fmt.Sprintf("Usuario %02d", i),
			fmt.Sprintf("login%02d", i),
			usr.NovoEmail(fmt.Sprintf("user%02d@email.com", i)),
			usr.PermUSR,
			nil,
		)

		repo.usuarios[id] = u
	}
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

// novoCriarUsuarioParamsInvalidosTeste cria parâmetros inválidos de criação de usuário para testes.
func novoCriarUsuarioParamsInvalidosTeste() usr.CriarParams {
	return usr.CriarParams{
		Nome:      "",
		Login:     "",
		Email:     usr.NovoEmail(""),
		Permissao: "",
	}
}

// novoAtualizarUsuarioParamsTeste cria parâmetros de atualização de usuário para testes.
func novoAtualizarUsuarioParamsTeste() usr.AtualizarParams {
	return usr.AtualizarParams{
		Nome:      ptr("Rogério Atualizado"),
		Login:     ptr("rogerioatualizado"),
		Email:     ptr(usr.NovoEmail("rogerioatualizado@email.com")),
		Permissao: ptr(usr.PermADM),
		Status:    ptr(true),
		Avatar:    nil,
	}
}

// novoAtualizarUsuarioParamsInvalidosTeste cria parâmetros inválidos de atualização de usuário para testes.
func novoAtualizarUsuarioParamsInvalidosTeste() usr.AtualizarParams {
	return usr.AtualizarParams{
		Nome:  ptr(""),
		Login: ptr(""),
		Email: ptr(usr.NovoEmail("")),
	}
}

// novoAtualizarPermissaoUsuarioParamsTeste cria parâmetros de atualização de permissão de usuário para testes.
func novoAtualizarPermissaoParamsTeste() usr.AtualizarPermissaoParams {
	return usr.AtualizarPermissaoParams{Permissao: usr.PermUSR}
}

// novoAtualizarPermissaoParamsInvalidosTeste cria parâmetros inválidos de atualização de permissão de usuário para testes.
func novoAtualizarPermissaoParamsInvalidosTeste() usr.AtualizarPermissaoParams {
	return usr.AtualizarPermissaoParams{Permissao: ""}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_UsuarioService_Criar testa o método Criar do serviço de usuário.
func Test_UsuarioService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() usr.Repository
		params             usr.CriarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "criar usuario com sucesso",
			geradorID: &fakeGeradorID{
				id: usuarioTesteID,
			},
			prepararRepo: func() usr.Repository {
				return novoFakeUsuarioRepository()
			},
			params:       novoCriarUsuarioParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)
				verificarPersistido(t, len(repo.(*fakeUsuarioRepository).usuarios))

				if usuario.Nome() != novoCriarUsuarioParamsTeste().Nome {
					t.Errorf("Nome esperado '%s', recebido '%s'",
						novoCriarUsuarioParamsTeste().Nome, usuario.Nome())
				}

				if usuario.Login() != novoCriarUsuarioParamsTeste().Login {
					t.Errorf("Login esperado '%s', recebido '%s'",
						novoCriarUsuarioParamsTeste().Login, usuario.Login())
				}

				if usuario.Email().String() != novoCriarUsuarioParamsTeste().Email.String() {
					t.Errorf("Email esperado '%s', recebido '%s'",
						novoCriarUsuarioParamsTeste().Email.String(), usuario.Email().String())
				}

				if usuario.Permissao() != novoCriarUsuarioParamsTeste().Permissao {
					t.Errorf("Permissão esperada '%s', recebida '%s'",
						novoCriarUsuarioParamsTeste().Permissao, usuario.Permissao())
				}

				if usuario.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", usuario.Avatar())
				}

				if usuario.ID() != usuarioTesteID {
					t.Errorf("ID esperado '%s', recebido '%s'", usuarioTesteID, usuario.ID())
				}
				if usuario.Status() != true {
					t.Errorf("Status esperado 'true', recebido '%v'", usuario.Status())
				}

				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errFakeGeradorID,
			},
			prepararRepo: func() usr.Repository {
				return novoFakeUsuarioRepository()
			},
			params:       novoCriarUsuarioParamsTeste(),
			erroEsperado: errFakeGeradorID,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeUsuarioRepository).usuarios))
				verificarNulo(t, usuario)
			},
		},
		{
			nome: "erro ao salvar usuário no repositorio",
			geradorID: &fakeGeradorID{
				id: usuarioTesteID,
			},
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			params:       novoCriarUsuarioParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeUsuarioRepository).usuarios))
				verificarNulo(t, usuario)
			},
		},
		{
			nome: "erro de validacao ao criar usuario",
			geradorID: &fakeGeradorID{
				id: usuarioTesteID,
			},
			prepararRepo: func() usr.Repository {
				return novoFakeUsuarioRepository()
			},
			params:       novoCriarUsuarioParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeUsuarioRepository).usuarios))
				verificarNulo(t, usuario)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(tt.geradorID, repo)
			usuarioCriado, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuarioCriado)
			}
		})
	}
}

// Test_UsuarioService_Atualizar testa o método Atualizar do serviço de usuário.
func Test_UsuarioService_Atualizar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		params             usr.AtualizarParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "atualizar usuario com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)

				if novoAtualizarUsuarioParamsTeste().Nome != nil {
					if usuario.Nome() != *novoAtualizarUsuarioParamsTeste().Nome {
						t.Errorf("Nome esperado '%s', recebido '%s'",
							*novoAtualizarUsuarioParamsTeste().Nome, usuario.Nome())
					}
				}

				if novoAtualizarUsuarioParamsTeste().Login != nil {
					if usuario.Login() != *novoAtualizarUsuarioParamsTeste().Login {
						t.Errorf("Login esperado '%s', recebido '%s'",
							*novoAtualizarUsuarioParamsTeste().Login, usuario.Login())
					}
				}

				if novoAtualizarUsuarioParamsTeste().Email != nil {
					if usuario.Email().String() != novoAtualizarUsuarioParamsTeste().Email.String() {
						t.Errorf("Email esperado '%s', recebido '%s'", (*novoAtualizarUsuarioParamsTeste().Email).String(), usuario.Email().String())
					}
				}

				if novoAtualizarUsuarioParamsTeste().Permissao != nil {
					if usuario.Permissao() != *novoAtualizarUsuarioParamsTeste().Permissao {
						t.Errorf("Permissão esperada '%s', recebida '%s'", *novoAtualizarUsuarioParamsTeste().Permissao, usuario.Permissao())
					}
				}

				if novoAtualizarUsuarioParamsTeste().Status != nil {
					if usuario.Status() != *novoAtualizarUsuarioParamsTeste().Status {
						t.Errorf("Status esperado '%v', recebido '%v'", *novoAtualizarUsuarioParamsTeste().Status, usuario.Status())
					}
				}

				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())

			},
		},
		{
			nome: "erro ao persistir usuario atualizado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "usuario nao encontrado para atualizar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário original no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro de validacao ao atualizar usuario",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário original no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuarioAtualizado, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuarioAtualizado)
			}
		})
	}
}

// Test_UsuarioService_BuscarPorID testa o método BuscarPorID do serviço de usuário.
func Test_UsuarioService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "buscar usuario por id com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)
				verificarIDs(t, usuarioTesteID, usuario.ID())
			},
		},
		{
			nome: "erro ao buscar usuario por id no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				_, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
			},
		},
		{
			nome: "usuario não encontrado por id",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				_, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuarioEncontrado, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuarioEncontrado)
			}
		})
	}
}

// Test_UsuarioService_BuscarPorLogin testa o método BuscarPorLogin do serviço de usuário.
func Test_UsuarioService_BuscarPorLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		login              string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "buscar usuario por login com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			login:        usuarioTesteLogin,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)
				verificarParamBusca(t, usuarioTesteLogin, usuario.Login())

			},
		},
		{
			nome: "erro ao buscar usuario por login no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			login:        usuarioTesteLogin,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				_, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
			},
		},
		{
			nome: "usuario não encontrado por login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			login:        "login-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				_, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuarioEncontrado, erroRecebido := service.BuscarPorLogin(ctx, tt.login)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuarioEncontrado)
			}
		})
	}
}

// Test_UsuarioService_AtualizarPermissao testa o método AtualizarPermissao do serviço de usuário.
func Test_UsuarioService_AtualizarPermissao(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		params             usr.AtualizarPermissaoParams
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "atualizar permissao com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)

				if usuario.Permissao() != novoAtualizarPermissaoParamsTeste().Permissao {
					t.Errorf("Permissão esperada '%s', recebida '%s'", novoAtualizarPermissaoParamsTeste().Permissao, usuario.Permissao())
				}

				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir usuario com permissao atualizada no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "usuario nao encontrado para atualizar permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro de validacao ao atualizar permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsInvalidosTeste(),
			erroEsperado: dmn.NovoErrosValidacao(),
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário original no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, err := service.AtualizarPermissao(ctx, tt.id, tt.params)

			verificarErro(t, err, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
			}
		})
	}
}

// Test_UsuarioService_AtualizarUltimoLogin testa o método AtualizarUltimoLogin do serviço de usuário.
func Test_UsuarioService_AtualizarUltimoLogin(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "atualizar ultimo login com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)
				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar ultimo login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir usuario com ultimo login atualizado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "usuario nao encontrado para atualizar ultimo login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, erroRecebido := service.AtualizarUltimoLogin(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
			}
		})
	}
}

// Test_UsuarioService_Desativar testa o método Desativar do serviço de usuário.
func Test_UsuarioService_Desativar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "desativar usuario com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Ativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)

				if usuario.Status() {
					t.Errorf("esperava usuário desativado, recebeu Status=%t", usuario.Status())
				}

				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para desativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Ativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNulo(t, usuario)

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir usuario desativado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Ativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "usuario nao encontrado para desativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, erroRecebido := service.Desativar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
			}
		})
	}
}

// Test_UsuarioService_Ativar testa o método Ativar do serviço de usuário.
func Test_UsuarioService_Ativar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuario *usr.Usuario)
	}{
		{
			nome: "ativar usuario com sucesso",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Desativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				verificarNaoNulo(t, usuario)

				if !usuario.Status() {
					t.Errorf("esperava usuário ativado, recebeu Status=%t", usuario.Status())
				}

				verificarIDs(t, usuarioTesteID, usuario.ID())
				verificarDatas(t, usuario.CriadoEm(), usuario.AtualizadoEm())
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para ativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Desativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())
			},
		},
		{
			nome: "erro ao persistir usuario ativado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Desativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())

			},
		},
		{
			nome: "usuario nao encontrado para ativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Ativar()
				repo.usuarios[usuarioTesteID] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				verificarOk(t, ok, "deveria existir usuário no repositório")
				verificarNulo(t, usuario)
				verificarNaoAlterado(t, novoUsuarioTeste(), usuarioRepo, compararUsuarios)
				verificarDatas(t, usuarioRepo.CriadoEm(), usuarioRepo.AtualizadoEm())

			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, erroRecebido := service.Ativar(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
			}
		})
	}
}

// Test_UsuarioService_Listar testa o método Listar do serviço de usuário.
func Test_UsuarioService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		filtro             usr.Filtro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro)
	}{
		{
			nome: "listar 10 usuarios paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				popularUsuariosTeste(repo, 10)
				return repo
			},
			filtro:       usr.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.usuarios))
				verificarLimite(t, usuarios, 5)
			},
		},
		{
			nome: "erro ao listar usuarios paginados no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				popularUsuariosTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       usr.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, usuarios, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuarios, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuarios, total, filtroRetornado)
			}
		})
	}
}

// Test_UsuarioService_VerificarPermissao testa o método VerificarPermissao do serviço de usuário.
func Test_UsuarioService_VerificarPermissao(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() usr.Repository
		id                 string
		permissoes         []usr.Permissao
		erroEsperado       error
		verificarResultado func(t *testing.T, repo usr.Repository, ok bool)
	}{
		{
			nome: "usuario possui permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			permissoes:   []usr.Permissao{usr.PermUSR, usr.PermADM},
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, ok bool) {
				t.Helper()

				if !ok {
					t.Error("esperava que o retorno ok fosse true, recebeu false")
				}
			},
		},
		{
			nome: "usuario nao possui permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			permissoes:   []usr.Permissao{usr.PermUSR},
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, ok bool) {
				t.Helper()

				if ok {
					t.Error("esperava que o retorno ok fosse false, recebeu true")
				}
			},
		},
		{
			nome: "erro ao buscar usuario",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			permissoes:   []usr.Permissao{usr.PermADM},
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, ok bool) {
				t.Helper()

				if ok {
					t.Error("não esperava valor ok quando há erro ao buscar usuário")
				}
			},
		},
		{
			nome: "usuario nao encontrado",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			permissoes:   []usr.Permissao{usr.PermADM},
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, ok bool) {
				t.Helper()

				if ok {
					t.Error("não esperava valor ok quando o usuário não é encontrado")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			ok, err := service.VerificarPermissao(ctx, tt.id, tt.permissoes)

			verificarErro(t, err, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, ok)
			}
		})
	}
}

// =====================================================================================================================
// Funções auxiliares para os testes
// =====================================================================================================================

// compararUsuarios verifica se dois usuários são iguais comparando seus campos.
func compararUsuarios(t *testing.T, antes, depois *usr.Usuario) {
	t.Helper()

	if antes == nil || depois == nil {
		t.Fatalf("usuários para comparação não podem ser nil: antes=%v depois=%v", antes, depois)
	}

	if antes.ID() != depois.ID() ||
		antes.Nome() != depois.Nome() ||
		antes.Login() != depois.Login() ||
		antes.Email().String() != depois.Email().String() ||
		antes.Permissao() != depois.Permissao() ||
		antes.Status() != depois.Status() {

		t.Fatalf(
			"usuários não deveriam ser diferentes:\nantes=%+v\ndepois=%+v",
			antes, depois,
		)
	}
}

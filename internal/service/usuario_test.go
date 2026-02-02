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

// --- Fake Repository ---------------------------------------------------

// fakeUsuarioRepository é uma implementação falsa de usr.Repository para testes.
// Usa estado em vez de funções injetadas, tornando os testes mais legíveis.
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
	usuario := u // Criar uma cópia para evitar efeitos colaterais
	f.usuarios[u.ID()] = &usuario
	return &usuario, nil
}

// Atualizar atualiza um usuário existente no repositório falso.
func (f *fakeUsuarioRepository) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	usuario := u // Criar uma cópia para evitar efeitos colaterais
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
// MOCKS
// =====================================================================================================================

const (
	usuarioTesteID    = "usuario-123"       // ID fixo para testes
	usuarioTesteNome  = "Rogério"           // Nome fixo para testes
	usuarioTesteLogin = "rogerio"           // Login fixo para testes
	usuarioTesteEmail = "rogerio@email.com" // Email fixo para testes
)

// novoUsuarioTeste cria um usuário de teste com os dados fornecidos.
func novoUsuarioTeste() *usr.Usuario {
	u, _ := usr.Novo(
		"usuario-123",
		"Rogério",
		"rogerio",
		usr.NovoEmail("rogerio@email.com"),
		usr.PermADM,
		nil,
	)
	return u
}

// popularUsuarios adiciona múltiplos usuários ao repositório falso para testes.
func popularUsuarios(repo *fakeUsuarioRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("user-%02d", i)

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

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário criado, recebeu nil")
				}

				// verifica se foi salvo no repositório
				_, ok := repo.(*fakeUsuarioRepository).usuarios[usuario.ID()]
				if !ok {
					t.Fatalf("usuário com ID '%s' não foi salvo no repositório", usuario.ID())
					return
				}

				// verificar se os parametros foram salvos corretamente
				if usuario.Nome() != novoCriarUsuarioParamsTeste().Nome {
					t.Errorf("Nome esperado '%s', recebido '%s'", novoCriarUsuarioParamsTeste().Nome, usuario.Nome())
				}
				if usuario.Login() != novoCriarUsuarioParamsTeste().Login {
					t.Errorf("Login esperado '%s', recebido '%s'", novoCriarUsuarioParamsTeste().Login, usuario.Login())
				}
				if usuario.Email().String() != novoCriarUsuarioParamsTeste().Email.String() {
					t.Errorf("Email esperado '%s', recebido '%s'", novoCriarUsuarioParamsTeste().Email.String(), usuario.Email().String())
				}
				if usuario.Permissao() != novoCriarUsuarioParamsTeste().Permissao {
					t.Errorf("Permissão esperada '%s', recebida '%s'", novoCriarUsuarioParamsTeste().Permissao, usuario.Permissao())
				}
				if usuario.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", usuario.Avatar())
				}

				// verificar se campos automáticos foram definidos corretamente
				if usuario.ID() != usuarioTesteID {
					t.Errorf("ID esperado '%s', recebido '%s'", usuarioTesteID, usuario.ID())
				}
				if usuario.Status() != true {
					t.Errorf("Status esperado 'true', recebido '%v'", usuario.Status())
				}
				if usuario.UltimoLogin().IsZero() {
					t.Errorf("ÚltimoLogin(%v) não deveria ser zero", usuario.UltimoLogin())
				}
				if usuario.CriadoEm().IsZero() {
					t.Errorf("CriadoEm(%v) não deveria ser zero", usuario.CriadoEm())
				}

				// verificar invariantes do usuário
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.AtualizadoEm(), usuario.CriadoEm())
				}
				if usuario.UltimoLogin().Before(usuario.CriadoEm()) {
					t.Errorf("ÚltimoLogin(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.UltimoLogin(), usuario.CriadoEm())
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando gerador de ID falha")
				}
				if usuario != nil {
					t.Error("não deveria retornar usuário quando gerador de ID falha")
				}
			},
		},
		{
			nome: "erro ao salvar no repositorio",
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando salvar falha")
				}
				if usuario != nil {
					t.Error("não deveria retornar usuário quando salvar falha")
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 0 {
					t.Error("repositório não deveria ter usuários quando validação falha")
				}
				if usuario != nil {
					t.Error("não deveria retornar usuário quando validação falha")
				}
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário atualizado, recebeu nil")
				}

				// verificar se os parametros foram atualizados corretamente
				if novoAtualizarUsuarioParamsTeste().Nome != nil {
					if usuario.Nome() != *novoAtualizarUsuarioParamsTeste().Nome {
						t.Errorf("Nome esperado '%s', recebido '%s'", *novoAtualizarUsuarioParamsTeste().Nome, usuario.Nome())
					}
				}
				if novoAtualizarUsuarioParamsTeste().Login != nil {
					if usuario.Login() != *novoAtualizarUsuarioParamsTeste().Login {
						t.Errorf("Login esperado '%s', recebido '%s'", *novoAtualizarUsuarioParamsTeste().Login, usuario.Login())
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

				// verificar se campos automáticos foram atualizados corretamente
				if usuario.AtualizadoEm().IsZero() {
					t.Errorf("AtualizadoEm(%v) não deveria ser zero", usuario.AtualizadoEm())
				}
				if usuario.CriadoEm().IsZero() {
					t.Errorf("CriadoEm(%v) não deveria ser zero", usuario.CriadoEm())
				}

				// verificar invariantes do usuário
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.AtualizadoEm(), usuario.CriadoEm())
				}
				if usuario.UltimoLogin().Before(usuario.CriadoEm()) {
					t.Errorf("ÚltimoLogin(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.UltimoLogin(), usuario.CriadoEm())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}

				if usuarioRepo.Nome() != novoUsuarioTeste().Nome() {
					t.Errorf("Nome do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Nome(), usuarioRepo.Nome())
				}
				if usuarioRepo.Login() != novoUsuarioTeste().Login() {
					t.Errorf("Login do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Login(), usuarioRepo.Login())
				}
				if usuarioRepo.Email().String() != novoUsuarioTeste().Email().String() {
					t.Errorf("Email do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Email().String(), usuarioRepo.Email().String())
				}
				if usuarioRepo.Permissao() != novoUsuarioTeste().Permissao() {
					t.Errorf("Permissão do usuário original esperada '%s', recebida '%s'", novoUsuarioTeste().Permissao(), usuarioRepo.Permissao())
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}
				if usuarioRepo.CriadoEm() != novoUsuarioTeste().CriadoEm() {
					t.Errorf("CriadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().CriadoEm(), usuarioRepo.CriadoEm())
				}
				if usuarioRepo.AtualizadoEm() != novoUsuarioTeste().AtualizadoEm() {
					t.Errorf("AtualizadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().AtualizadoEm(), usuarioRepo.AtualizadoEm())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro ao persistir usuario atualizado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}

				if usuarioRepo.Nome() != novoUsuarioTeste().Nome() {
					t.Errorf("Nome do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Nome(), usuarioRepo.Nome())
				}
				if usuarioRepo.Login() != novoUsuarioTeste().Login() {
					t.Errorf("Login do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Login(), usuarioRepo.Login())
				}
				if usuarioRepo.Email().String() != novoUsuarioTeste().Email().String() {
					t.Errorf("Email do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Email().String(), usuarioRepo.Email().String())
				}
				if usuarioRepo.Permissao() != novoUsuarioTeste().Permissao() {
					t.Errorf("Permissão do usuário original esperada '%s', recebida '%s'", novoUsuarioTeste().Permissao(), usuarioRepo.Permissao())
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}
				if usuarioRepo.CriadoEm() != novoUsuarioTeste().CriadoEm() {
					t.Errorf("CriadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().CriadoEm(), usuarioRepo.CriadoEm())
				}
				if usuarioRepo.AtualizadoEm() != novoUsuarioTeste().AtualizadoEm() {
					t.Errorf("AtualizadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().AtualizadoEm(), usuarioRepo.AtualizadoEm())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao atualizar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario nao encontrado para atualizar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarUsuarioParamsTeste(),
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro de validacao ao atualizar usuario",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarUsuarioParamsInvalidosTeste(),
			erroEsperado: &dmn.ErrosValidacao{},
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}

				if usuarioRepo.Nome() != novoUsuarioTeste().Nome() {
					t.Errorf("Nome do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Nome(), usuarioRepo.Nome())
				}
				if usuarioRepo.Login() != novoUsuarioTeste().Login() {
					t.Errorf("Login do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Login(), usuarioRepo.Login())
				}
				if usuarioRepo.Email().String() != novoUsuarioTeste().Email().String() {
					t.Errorf("Email do usuário original esperado '%s', recebido '%s'", novoUsuarioTeste().Email().String(), usuarioRepo.Email().String())
				}
				if usuarioRepo.Permissao() != novoUsuarioTeste().Permissao() {
					t.Errorf("Permissão do usuário original esperada '%s', recebida '%s'", novoUsuarioTeste().Permissao(), usuarioRepo.Permissao())
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}
				if usuarioRepo.CriadoEm() != novoUsuarioTeste().CriadoEm() {
					t.Errorf("CriadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().CriadoEm(), usuarioRepo.CriadoEm())
				}
				if usuarioRepo.AtualizadoEm() != novoUsuarioTeste().AtualizadoEm() {
					t.Errorf("AtualizadoEm do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().AtualizadoEm(), usuarioRepo.AtualizadoEm())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro de validação: '%v'", usuario)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, erroRecebido := service.Atualizar(ctx, tt.id, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário encontrado, recebeu nil")
				}

				// verifica se o usuário retornado é o esperado
				if usuario.ID() != usuarioTesteID {
					t.Errorf("ID esperado '%s', recebido '%s'", usuarioTesteID, usuario.ID())
				}
			},
		},
		{
			nome: "erro ao buscar usuario por id no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario não encontrado por id",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, err := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, err, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			login:        usuarioTesteLogin,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verfica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário encontrado, recebeu nil")
				}

				// verifica se o usuário retornado é o esperado
				if usuario.Login() != usuarioTesteLogin {
					t.Errorf("Login esperado '%s', recebido '%s'", usuarioTesteLogin, usuario.Login())
				}
			},
		},
		{
			nome: "erro ao buscar usuario por login no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			login:        usuarioTesteLogin,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario não encontrado por login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			login:        "login-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com login diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.Login() == "login-inexistente" {
					t.Error("o login buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoUsuarioService(nil, repo)
			usuario, err := service.BuscarPorLogin(ctx, tt.login)

			verificarErro(t, err, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário atualizado, recebeu nil")
				}

				// verificar se a permissão foi atualizada corretamente
				if usuario.Permissao() != novoAtualizarPermissaoParamsTeste().Permissao {
					t.Errorf("Permissão esperada '%s', recebida '%s'", novoAtualizarPermissaoParamsTeste().Permissao, usuario.Permissao())
				}

				// verificar se campos automáticos foram atualizados corretamente
				if usuario.AtualizadoEm().IsZero() {
					t.Errorf("AtualizadoEm(%v) não deveria ser zero", usuario.AtualizadoEm())
				}
				if usuario.CriadoEm().IsZero() {
					t.Errorf("CriadoEm(%v) não deveria ser zero", usuario.CriadoEm())
				}

				// verificar invariantes do usuário
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.AtualizadoEm(), usuario.CriadoEm())
				}
				if usuario.UltimoLogin().Before(usuario.CriadoEm()) {
					t.Errorf("ÚltimoLogin(%v) não deveria ser anterior a CriadoEm(%v)",
						usuario.UltimoLogin(), usuario.CriadoEm())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro ao persistir usuario com permissao atualizada no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Permissao() != novoUsuarioTeste().Permissao() {
					t.Errorf("Permissão do usuário original esperada '%s', recebida '%s'", novoUsuarioTeste().Permissao(), usuarioRepo.Permissao())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao atualizar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario nao encontrado para atualizar permissao",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			params:       novoAtualizarPermissaoParamsTeste(),
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
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

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Permissao() != novoUsuarioTeste().Permissao() {
					t.Errorf("Permissão do usuário original esperada '%s', recebida '%s'", novoUsuarioTeste().Permissao(), usuarioRepo.Permissao())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro de validação: '%v'", usuario)
				}
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário atualizado, recebeu nil")
				}

				// verificar se o ultimo login foi atualizado corretamente
				if usuario.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if usuario.UltimoLogin().Before(usuario.CriadoEm()) {
					t.Errorf("UltimoLogin(%v) não pode ser anterior a CriadoEm(%v)", usuario.UltimoLogin(), usuario.CriadoEm())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para atualizar ultimo login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.UltimoLogin() != novoUsuarioTeste().UltimoLogin() {
					t.Errorf("UltimoLogin do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().UltimoLogin(), usuarioRepo.UltimoLogin())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro ao persistir usuario com ultimo login atualizado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}

				if usuarioRepo.UltimoLogin() != novoUsuarioTeste().UltimoLogin() {
					t.Errorf("UltimoLogin do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().UltimoLogin(), usuarioRepo.UltimoLogin())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao atualizar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario nao encontrado para atualizar ultimo login",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
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
				novoUsuario.Ativar() // garante que o usuário está ativo antes de desativar
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário atualizado, recebeu nil")
				}

				// verificar se o status foi atualizado corretamente
				if usuario.Status() {
					t.Errorf("esperava usuário desativado, recebeu Status=%v", usuario.Status())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para desativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Ativar() // garante que o usuário está ativo antes de desativar
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro ao persistir usuario desativado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao atualizar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario nao encontrado para desativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
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
				novoUsuario.Desativar() // garante que o usuário está desativado antes de ativar
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário foi retornado
				if usuario == nil {
					t.Fatal("esperava usuário atualizado, recebeu nil")
				}

				// verificar se o status foi atualizado corretamente
				if !usuario.Status() {
					t.Errorf("esperava usuário ativado, recebeu Status=%v", usuario.Status())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para ativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Desativar() // garante que o usuário está desativado antes de ativar
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
			},
		},
		{
			nome: "erro ao persistir usuario ativado no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				novoUsuario.Desativar() // garante que o usuário está desativado antes de ativar
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se usuário original permanece no repositório sem atualização
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuarioRepo.Status() != novoUsuarioTeste().Status() {
					t.Errorf("Status do usuário original esperado '%v', recebido '%v'", novoUsuarioTeste().Status(), usuarioRepo.Status())
				}

				// verifica se não foi retornado usuário atualizado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao atualizar: '%v'", usuario)
				}
			},
		},
		{
			nome: "usuario nao encontrado para ativar",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoBuscar = mysql.ErrUsuarioNaoEncontrado
				return repo
			},
			id:           "id-inexistente",
			erroEsperado: mysql.ErrUsuarioNaoEncontrado,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, ok := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !ok {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado usuário encontrado
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao buscar: '%v'", usuario)
				}
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
				popularUsuarios(repo, 10)
				return repo
			},
			filtro:       usr.NovoFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 10 {
					t.Errorf("esperava 10 usuários no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if filtroRetornado.Pagina() != 1 || filtroRetornado.Limite() != 5 {
					t.Errorf("esperava filtro retornado com página 1 e limite 5, recebeu página %d e limite %d", filtroRetornado.Pagina(), filtroRetornado.Limite())
				}
				if total != 10 {
					t.Errorf("esperava total 10 usuários retornados, recebeu %d", total)
				}
				if len(usuarios) != 5 {
					t.Errorf("esperava 5 usuários na página, recebeu %d", len(usuarios))
				}
			},
		},
		{
			nome: "listar 10 usuarios paginados com sucesso - pagina 3 limite 4",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				popularUsuarios(repo, 10)
				return repo
			},
			filtro:       usr.NovoFiltro(dmn.NovoPaginacao(3, 4), nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 10 {
					t.Errorf("esperava 10 usuários no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if filtroRetornado.Pagina() != 3 || filtroRetornado.Limite() != 4 {
					t.Errorf("esperava filtro retornado com página 3 e limite 4, recebeu página %d e limite %d", filtroRetornado.Pagina(), filtroRetornado.Limite())
				}
				if total != 10 {
					t.Errorf("esperava total 10 usuários retornados, recebeu %d", total)
				}
				if len(usuarios) != 2 {
					t.Errorf("esperava 2 usuários na página, recebeu %d", len(usuarios))
				}
			},
		},
		{
			nome: "erro ao listar usuarios paginados no repositório",
			prepararRepo: func() usr.Repository {
				repo := novoFakeUsuarioRepository()
				novoUsuario := novoUsuarioTeste()
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       usr.NovoFiltro(dmn.NovoPaginacao(1, 10), nil, nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuarios []usr.Usuario, total int, filtroRetornado usr.Filtro) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if len(usuarios) != 0 {
					t.Errorf("não esperava usuários retornados quando há erro ao listar: '%v'", usuarios)
				}
				if total != 0 {
					t.Errorf("não esperava total de usuários quando há erro ao listar: %d", total)
				}
				if filtroRetornado.Pagina() != 0 || filtroRetornado.Limite() != 0 {
					t.Errorf("deve retornar filtro vazio quando há erro ao listar")
				}
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

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if !ok {
					t.Error("esperava que o usuário possuísse a permissão")
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

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if ok {
					t.Error("não esperava que o usuário possuísse a permissão")
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

				// verifica se havia 1 usuário no repositório
				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
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

				// verifica se havia 1 usuário no repositório com id diferente
				usuarioRepo, existe := repo.(*fakeUsuarioRepository).usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário no repositório")
					return
				}
				if usuarioRepo.ID() == "id-inexistente" {
					t.Error("o id buscado não deve existir no repositório")
				}

				// verifica se não foi retornado ok
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

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
	f.usuarios[u.ID()] = &u
	return &u, nil
}

// Atualizar atualiza um usuário existente no repositório falso.
func (f *fakeUsuarioRepository) Atualizar(ctx context.Context, id string, u usr.Usuario) (*usr.Usuario, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	f.usuarios[id] = &u
	return &u, nil
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
func popularUsuarios(repo *fakeUsuarioRepository2, quantidade int) {
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
	return usr.AtualizarPermissaoParams{
		Permissao: usr.PermUSR,
	}
}

// novoAtualizarPermissaoParamsInvalidosTeste cria parâmetros inválidos de atualização de permissão de usuário para testes.
func novoAtualizarPermissaoParamsInvalidosTeste() usr.AtualizarPermissaoParams {
	return usr.AtualizarPermissaoParams{
		Permissao: "",
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestUsuarioService_Criar testa o método Criar do serviço de usuário.
func TestUsuarioService_Criar(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario.ID() != usuarioTesteID {
					t.Errorf("esperava ID '%s', recebeu '%s'", usuarioTesteID, usuario.ID())
				}
				if usuario.Nome() != "Rogério" {
					t.Errorf("esperava nome 'Rogério', recebeu '%s'", usuario.Nome())
				}
				if usuario.Login() != "rogerio" {
					t.Errorf("esperava login 'rogerio', recebeu '%s'", usuario.Login())
				}
				if usuario.Email() != "rogerio@email.com" {
					t.Errorf("esperava email 'rogerio@email.com', recebeu '%s'", usuario.Email())
				}
				if usuario.Permissao() != usr.PermADM {
					t.Errorf("esperava permissão 'ADM', recebeu '%s'", usuario.Permissao())
				}
				if !usuario.Status() {
					t.Error("esperava usuário ativo ao criar")
				}
				if usuario.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if usuario.CriadoEm().IsZero() {
					t.Error("esperava campo CriadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf(
						"AtualizadoEm não pode ser anterior a CriadoEm (CriadoEm=%v, AtualizadoEm=%v)",
						usuario.CriadoEm(),
						usuario.AtualizadoEm(),
					)
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
			usuario, erroRecebido := service.Criar(ctx, tt.params)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, usuario)
			}
		})
	}
}

// TestUsuarioService_Atualizar2 testa o método Atualizar do serviço de usuário.
func TestUsuarioService_Atualizar(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				atual, existe := fakeRepo.usuarios[usuario.ID()]
				if !existe {
					t.Error("deveia existir usuário atualizado no repositório")
					return
				}
				if atual.Nome() != "Rogério Atualizado" {
					t.Errorf("Nome esperado 'Rogério Atualizado', recebido '%s'", atual.Nome())
				}
				if atual.Login() != "rogerioatualizado" {
					t.Errorf("Login esperado 'rogerioatualizado', recebido '%s'", atual.Login())
				}
				if atual.Email() != "rogerioatualizado@email.com" {
					t.Errorf("Email esperado 'rogerioatualizado@email.com', recebido '%s'", atual.Email())
				}
				if atual.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", atual.Permissao())
				}
				if !atual.Status() {
					t.Error("Status esperado 'true', recebido 'false'")
				}
				if atual.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", atual.Avatar())
				}
				if atual.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if atual.AtualizadoEm().Before(atual.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não pode ser anterior a CriadoEm(%v)", atual.AtualizadoEm(), atual.CriadoEm())
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Nome() != "Rogério" {
					t.Errorf("Nome esperado 'Rogério', recebido '%s'", original.Nome())
				}
				if original.Login() != "rogerio" {
					t.Errorf("Login esperado 'rogerio', recebido '%s'", original.Login())
				}
				if original.Email() != "rogerio@email.com" {
					t.Errorf("Email esperado 'rogerio@email.com', recebido '%s'", original.Email())
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
				if !original.Status() {
					t.Error("Status esperado 'true', recebido 'false'")
				}
				if original.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", original.Avatar())
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Nome() != "Rogério" {
					t.Errorf("Nome esperado 'Rogério', recebido '%s'", original.Nome())
				}
				if original.Login() != "rogerio" {
					t.Errorf("Login esperado 'rogerio', recebido '%s'", original.Login())
				}
				if original.Email() != "rogerio@email.com" {
					t.Errorf("Email esperado 'rogerio@email.com', recebido '%s'", original.Email())
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
				if original.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", original.Avatar())
				}
				if original.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if original.AtualizadoEm().Before(original.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não pode ser anterior a CriadoEm(%v)", original.AtualizadoEm(), original.CriadoEm())
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Nome() != "Rogério" {
					t.Errorf("Nome esperado 'Rogério', recebido '%s'", original.Nome())
				}
				if original.Login() != "rogerio" {
					t.Errorf("Login esperado 'rogerio', recebido '%s'", original.Login())
				}
				if original.Email() != "rogerio@email.com" {
					t.Errorf("Email esperado 'rogerio@email.com', recebido '%s'", original.Email())
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
				if original.Avatar() != nil {
					t.Errorf("Avatar esperado 'nil', recebido '%v'", original.Avatar())
				}
				if original.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if original.AtualizadoEm().Before(original.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não pode ser anterior a CriadoEm(%v)", original.AtualizadoEm(), original.CriadoEm())
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

func TestUsuarioService_BuscarPorID(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario.ID() != usuarioTesteID {
					t.Errorf("esperava ID '%s', recebeu '%s'", usuarioTesteID, usuario.ID())
				}
				if usuario.Nome() != "Rogério" {
					t.Errorf("esperava nome 'Rogério', recebeu '%s'", usuario.Nome())
				}
				if usuario.Login() != "rogerio" {
					t.Errorf("esperava login 'rogerio', recebeu '%s'", usuario.Login())
				}
				if usuario.Email() != "rogerio@email.com" {
					t.Errorf("esperava email 'rogerio@email.com', recebeu '%s'", usuario.Email())
				}
				if usuario.Permissao() != usr.PermADM {
					t.Errorf("esperava permissão 'ADM', recebeu '%s'", usuario.Permissao())
				}
				if !usuario.Status() {
					t.Error("esperava usuário ativo ao criar")
				}
				if usuario.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if usuario.CriadoEm().IsZero() {
					t.Error("esperava campo CriadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf(
						"AtualizadoEm não pode ser anterior a CriadoEm (CriadoEm=%v, AtualizadoEm=%v)",
						usuario.CriadoEm(),
						usuario.AtualizadoEm(),
					)
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

// TestUsuarioService_BuscarPorLogin2 testa o método BuscarPorLogin do serviço de usuário.
func TestUsuarioService_BuscarPorLogin(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario.ID() != usuarioTesteID {
					t.Errorf("esperava ID '%s', recebeu '%s'", usuarioTesteID, usuario.ID())
				}
				if usuario.Nome() != "Rogério" {
					t.Errorf("esperava nome 'Rogério', recebeu '%s'", usuario.Nome())
				}
				if usuario.Login() != usuarioTesteLogin {
					t.Errorf("esperava login '%s', recebeu '%s'", usuarioTesteLogin, usuario.Login())
				}
				if usuario.Email() != "rogerio@email.com" {
					t.Errorf("esperava email 'rogerio@email.com', recebeu '%s'", usuario.Email())
				}
				if usuario.Permissao() != usr.PermADM {
					t.Errorf("esperava permissão 'ADM', recebeu '%s'", usuario.Permissao())
				}
				if !usuario.Status() {
					t.Error("esperava usuário ativo ao criar")
				}
				if usuario.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if usuario.CriadoEm().IsZero() {
					t.Error("esperava campo CriadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if usuario.AtualizadoEm().Before(usuario.CriadoEm()) {
					t.Errorf(
						"AtualizadoEm não pode ser anterior a CriadoEm (CriadoEm=%v, AtualizadoEm=%v)",
						usuario.CriadoEm(),
						usuario.AtualizadoEm(),
					)
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

// TestUsuarioService_BuscarCategoriaPermissaoPorID2 testa o método BuscarCategoriaPermissaoPorID do serviço de usuário.
func TestUsuarioService_AtualizarPermissao(t *testing.T) {
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
				fakeRepo := repo.(*fakeUsuarioRepository)
				atualizado, existe := fakeRepo.usuarios[usuario.ID()]
				if !existe {
					t.Error("deveia existir usuário atualizado no repositório")
					return
				}
				if atualizado.Permissao() != usr.PermUSR {
					t.Errorf("Permissão esperada 'USR', recebida '%s'", atualizado.Permissao())
				}
				if atualizado.AtualizadoEm().IsZero() {
					t.Error("esperava campo AtualizadoEm preenchido, mas está zerado")
				}
				if atualizado.AtualizadoEm().Before(atualizado.CriadoEm()) {
					t.Errorf("AtualizadoEm(%v) não pode ser anterior a CriadoEm(%v)", atualizado.AtualizadoEm(), atualizado.CriadoEm())
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				original, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if original.Permissao() != usr.PermADM {
					t.Errorf("Permissão esperada 'ADM', recebida '%s'", original.Permissao())
				}
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

// TestUsuarioService_BuscarCategoriaPermissaoPorID2 testa o método BuscarCategoriaPermissaoPorID do serviço de usuário.
func TestUsuarioService_AtualizarUltimoLogin(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				atual, existe := fakeRepo.usuarios[usuario.ID()]
				if !existe {
					t.Error("deveia existir usuário atualizado no repositório")
					return
				}
				if atual.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if atual.UltimoLogin().Before(atual.CriadoEm()) {
					t.Errorf("UltimoLogin(%v) não pode ser anterior a CriadoEm(%v)", atual.UltimoLogin(), atual.CriadoEm())
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao persistir: '%v'", usuario)
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

// TestUsuarioService_Desativar2 testa o método Desativar do serviço de usuário.
func TestUsuarioService_Desativar(t *testing.T) {
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				atual, existe := fakeRepo.usuarios[usuario.ID()]
				if !existe {
					t.Error("deveia existir usuário atualizado no repositório")
					return
				}
				if atual.Status() {
					t.Errorf("esperava usuário desativado, mas recebeu Status=%v", atual.Status())
				}
				if atual.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if atual.UltimoLogin().Before(atual.CriadoEm()) {
					t.Errorf("UltimoLogin(%v) não pode ser anterior a CriadoEm(%v)", atual.UltimoLogin(), atual.CriadoEm())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para desativar",
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao persistir: '%v'", usuario)
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

// TestUsuarioService_Ativar2 testa o método Ativar do serviço de usuário.
func TestUsuarioService_Ativar(t *testing.T) {
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				atual, existe := fakeRepo.usuarios[usuario.ID()]
				if !existe {
					t.Error("deveia existir usuário atualizado no repositório")
					return
				}
				if !atual.Status() {
					t.Errorf("esperava usuário ativado, recebeu Status=%v", atual.Status())
				}
				if atual.UltimoLogin().IsZero() {
					t.Error("esperava campo UltimoLogin preenchido, mas está zerado")
				}
				if atual.UltimoLogin().Before(atual.CriadoEm()) {
					t.Errorf("UltimoLogin(%v) não pode ser anterior a CriadoEm(%v)", atual.UltimoLogin(), atual.CriadoEm())
				}
			},
		},
		{
			nome: "erro ao buscar usuario no repositório para ativar",
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
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
				repo.usuarios[novoUsuario.ID()] = novoUsuario
				repo.erroAoAtualizar = errFakeRepo
				return repo
			},
			id:           usuarioTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo usr.Repository, usuario *usr.Usuario) {
				t.Helper()

				fakeRepo := repo.(*fakeUsuarioRepository)
				_, existe := fakeRepo.usuarios[usuarioTesteID]
				if !existe {
					t.Error("deve existir usuário original no repositório sem atualização")
					return
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando há erro ao persistir: '%v'", usuario)
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if usuario != nil {
					t.Errorf("não esperava usuário retornado quando usuário não é encontrado: '%v'", usuario)
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

// TestUsuarioService_Listar2 testa o método Listar do serviço de usuário.
func TestUsuarioService_Listar(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository2)
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

// TestUsuarioService_VerificarPermissao2 testa o método VerificarPermissao do serviço de usuário.
func TestUsuarioService_VerificarPermissao(t *testing.T) {
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

				fakeRepo := repo.(*fakeUsuarioRepository)
				if len(fakeRepo.usuarios) != 1 {
					t.Errorf("esperava 1 usuário no repositório, tem %d", len(fakeRepo.usuarios))
				}
				if ok {
					t.Error("não esperava valor ok quando usuário não é encontrado")
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

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeAtendimentoRepository é uma implementação fake do repositório de atendimentos para testes.
type fakeAtendimentoRepository struct {
	atendimentos    map[string]*atd.Atendimento
	erroAoCriar     error
	erroAoAtualizar error
	erroAoBuscar    error
	erroAoListar    error
	totalListar     int
}

func novoFakeAtendimentoRepository() *fakeAtendimentoRepository {
	return &fakeAtendimentoRepository{
		atendimentos: make(map[string]*atd.Atendimento),
	}
}

// Criar adiciona um novo atendimento ao repositório fake.
func (f *fakeAtendimentoRepository) Criar(ctx context.Context, a atd.Atendimento) (*atd.Atendimento, error) {
	if f.erroAoCriar != nil {
		return nil, f.erroAoCriar
	}
	f.atendimentos[a.ID()] = &a
	return &a, nil
}

// Atualizar atualiza um atendimento existente no repositório fake.
func (f *fakeAtendimentoRepository) Atualizar(ctx context.Context, id string, a atd.Atendimento) (*atd.Atendimento, error) {
	if f.erroAoAtualizar != nil {
		return nil, f.erroAoAtualizar
	}
	f.atendimentos[id] = &a
	return &a, nil
}

// BuscarPorID busca um atendimento pelo ID no repositório fake.
func (f *fakeAtendimentoRepository) BuscarPorID(ctx context.Context, id string) (*atd.Atendimento, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	a, existe := f.atendimentos[id]
	if !existe {
		return nil, mysql.ErrAtendimentoNaoEncontrado
	}
	return a, nil
}

// BuscarPorChamadoEAtribuidoID busca um atendimento pelo ID do chamado e do atribuído no repositório fake.
func (f *fakeAtendimentoRepository) BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	for _, a := range f.atendimentos {
		if a.ChamadoID() == chamadoID && a.AtribuidoID() == atribuidoID {
			return a, nil
		}
	}
	return nil, mysql.ErrAtendimentoNaoEncontrado
}

// Listar lista atendimentos com base no filtro fornecido no repositório fake.
func (f *fakeAtendimentoRepository) Listar(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	atendimentos := make([]atd.Atendimento, 0, len(f.atendimentos))
	for _, a := range f.atendimentos {
		atendimentos = append(atendimentos, *a)
	}

	total := f.totalListar
	if total == 0 {
		total = len(atendimentos)
	}

	return atendimentos, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

const atendimentoTesteID = "atendimento-123"

// novoAtendimentoTeste cria um novo atendimento para testes.
func novoAtendimentoTeste() *atd.Atendimento {
	a, _ := atd.Novo(
		atendimentoTesteID,
		usuarioTesteID,
		chamadoTesteID,
	)
	return a
}

// novoCriarAtendimentoParamsTeste cria parâmetros de criação de atendimento para testes.
func novoCriarAtendimentoParamsTeste() atd.CriarParams {
	return atd.CriarParams{
		AtribuidoID: usuarioTesteID,
		ChamadoID:   chamadoTesteID,
	}
}

// novoAtualizarAtendimentoParamsTeste cria parâmetros de atualização de atendimento para testes.
func novoAtualizarAtendimentoParamsTeste() atd.AtualizarParams {
	return atd.AtualizarParams{
		AtribuidoID: usuarioTesteID,
		ChamadoID:   chamadoTesteID,
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestAtendimentoService_Criar testa o método Criar do serviço de atendimento.
func TestAtendimentoService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repoSetup func() atd.Repository
		params    atd.CriarParams
		wantErr   bool
		assertFn  func(t *testing.T, repo atd.Repository, a *atd.Atendimento)
	}{
		{
			name: "Criação bem-sucedida",
			geradorID: &fakeGeradorID{
				id: "atd-123",
			},
			repoSetup: func() atd.Repository {
				return novoFakeAtendimentoRepository()
			},
			params:  novoCriarAtendimentoParamsTeste(),
			wantErr: false,
			assertFn: func(t *testing.T, repo atd.Repository, a *atd.Atendimento) {
				t.Helper()

				if a.ID() != "atd-123" {
					t.Errorf("ID esperado 'atd-123', obtido '%s'", a.ID())
				}

				fake := repo.(*fakeAtendimentoRepository)
				if len(fake.atendimentos) != 1 {
					t.Errorf("Esperado 1 atendimento no repositório, obtido %d", len(fake.atendimentos))
				}
			},
		},
		{
			name: "Erro ao criar atendimento no repositório",
			geradorID: &fakeGeradorID{
				id: "atd-456",
			},
			repoSetup: func() atd.Repository {
				repo := novoFakeAtendimentoRepository()
				repo.erroAoCriar = errors.New("erro ao criar atendimento")
				return repo
			},
			params:  novoCriarAtendimentoParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo atd.Repository, a *atd.Atendimento) {
				t.Helper()

				// repositório não deve conter atendimentos
				fake := repo.(*fakeAtendimentoRepository)
				if len(fake.atendimentos) != 0 {
					t.Errorf("Esperado 0 atendimentos no repositório, obtido %d", len(fake.atendimentos))
				}
			},
		},
		{
			name: "Erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errors.New("erro ao gerar ID"),
			},
			repoSetup: func() atd.Repository {
				return novoFakeAtendimentoRepository()
			},
			params:  novoCriarAtendimentoParamsTeste(),
			wantErr: true,
			assertFn: func(t *testing.T, repo atd.Repository, a *atd.Atendimento) {
				t.Helper()

				// repositório não deve conter atendimentos
				fake := repo.(*fakeAtendimentoRepository)
				if len(fake.atendimentos) != 0 {
					t.Errorf("Esperado 0 atendimentos no repositório, obtido %d", len(fake.atendimentos))
				}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			chamadoRepo := novoFakeChamadoRepository()
			chamadoRepo.chamados["chm-456"] = novoChamadoTeste()
			
			categoriaPermissaoRepo := novoFakeCategoriaPermissaoRepository()
			categoriaPermissaoRepo.categoriasPermissoes[chaveCompostaCategoriaPermissao] = novoCategoriaPermissaoTeste()
			
			usuarioRepo := novoFakeUsuarioRepository()
			novoUsuario := novoUsuarioTeste()
			usuarioRepo.usuarios[novoUsuario.ID()] = novoUsuario
			
			service := NovoAtendimentoService(tt.geradorID, repo, chamadoRepo, categoriaPermissaoRepo, usuarioRepo)

			a, err := service.Criar(ctx, tt.params)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, repo, a)
			}
		})
	}
}

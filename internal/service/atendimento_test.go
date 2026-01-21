package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeAtendimentoRepository é uma implementação falsa do repositório de atendimentos, usada para testes.
type fakeAtendimentoRepository struct {
	buscarFn                       func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error)
	atualizarFn                    func(ctx context.Context, atendimento atd.Atendimento) (*atd.Atendimento, error)
	criarFn                        func(ctx context.Context, atendimento atd.Atendimento) (*atd.Atendimento, error)
	buscarPorIDFn                  func(ctx context.Context, id string) (*atd.Atendimento, error)
	buscarPorChamadoEAtribuidoIDFn func(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error)
	listarFn                       func(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error)
}

// BuscarPorChamadoEAtribuidoID busca um atendimento pelo ID do chamado e pelo ID do usuário atribuído.
func (f *fakeAtendimentoRepository) BuscarPorChamadoEAtribuidoID(
	ctx context.Context,
	chamadoID, atribuidoID string,
) (*atd.Atendimento, error) {

	if f.buscarPorChamadoEAtribuidoIDFn == nil {
		return nil, nil
	}

	return f.buscarPorChamadoEAtribuidoIDFn(ctx, chamadoID, atribuidoID)
}

// Atualizar chama a função de atualização definida na struct falsa.
func (f *fakeAtendimentoRepository) Atualizar(ctx context.Context, id string, a atd.Atendimento) (*atd.Atendimento, error) {
	return f.atualizarFn(ctx, a)
}

// Criar chama a função de criação definida na struct falsa.
func (f *fakeAtendimentoRepository) Criar(ctx context.Context, a atd.Atendimento) (*atd.Atendimento, error) {
	return f.criarFn(ctx, a)
}

// BuscarPorID chama a função de busca por ID definida na struct falsa.
func (f *fakeAtendimentoRepository) BuscarPorID(ctx context.Context, id string) (*atd.Atendimento, error) {
	return f.buscarPorIDFn(ctx, id)
}

// Listar chama a função de listagem definida na struct falsa.
func (f *fakeAtendimentoRepository) Listar(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestAtendimentoService_BuscarPorID testa o método BuscarPorID do AtendimentoService.
func TestAtendimentoService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	atendimentoOK := &atd.Atendimento{}

	tests := []struct {
		name    string
		repo    *fakeAtendimentoRepository
		want    *atd.Atendimento
		wantErr bool
	}{
		{
			name: "buscar atendimento com sucesso",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					return atendimentoOK, nil
				},
			},
			want:    atendimentoOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar atendimento",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					return nil, errors.New("erro no banco")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "atendimento nao encontrado",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					return nil, errors.New("atendimento nao encontrado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAtendimentoService(nil, tt.repo, nil, nil, nil)

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

// TestAtendimentoService_BuscarPorChamadoEAtribuidoID testa o método BuscarPorChamadoEAtribuidoID do AtendimentoService.
func TestAtendimentoService_BuscarPorChamadoEAtribuidoID(t *testing.T) {
	ctx := context.Background()

	atendimentoOK := &atd.Atendimento{}

	tests := []struct {
		name    string
		repo    *fakeAtendimentoRepository
		want    *atd.Atendimento
		wantErr bool
	}{
		{
			name: "buscar atendimento com sucesso",
			repo: &fakeAtendimentoRepository{
				buscarFn: func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error) {
					return atendimentoOK, nil
				},
			},
			want:    atendimentoOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar atendimento",
			repo: &fakeAtendimentoRepository{
				buscarFn: func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error) {
					return nil, errors.New("erro no banco")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "atendimento nao encontrado",
			repo: &fakeAtendimentoRepository{
				buscarFn: func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error) {
					return nil, errors.New("atendimento nao encontrado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAtendimentoService(nil, tt.repo, nil, nil, nil)

			got, err := service.BuscarPorChamadoEAtribuidoID(ctx, "chamado-123", "user-123")

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

// TestAtendimentoService_Criar testa o método Criar do AtendimentoService.
func TestAtendimentoService_Criar(t *testing.T) {
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
		name                   string
		geradorID              dmn.GeradorID
		repo                   atd.Repository
		chamadoRepo            chm.Repository
		usuarioRepo            usr.Repository
		categoriaPermissaoRepo cpm.Repository
		wantErr                bool
	}{
		{
			name: "criar atendimento com sucesso",
			geradorID: &fakeGeradorID{
				id: "atendimento-123",
			},
			repo: &fakeAtendimentoRepository{
				criarFn: func(ctx context.Context, a atd.Atendimento) (*atd.Atendimento, error) {
					return &a, nil
				},
				buscarPorChamadoEAtribuidoIDFn: func(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error) {
					return nil, nil // NÃO existe atendimento ativo
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return chm.Novo(
						"chamado-123",
						"categoria-123",
						"subcategoria-123",
						"criador-123",
						"Título do chamado",
						"Descrição do chamado",
					)
				},
				atualizarFn: func(ctx context.Context, id string, chamado chm.Chamado) (*chm.Chamado, error) {
					return &chamado, nil
				},
			},
			usuarioRepo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usr.Novo(
						id,
						"Técnico Teste",
						"tecnico.teste@example.com",
						usr.NovoEmail("tecnico@gmail.com"),
						usr.PermTEC,
						nil,
					)
				},
			},
			categoriaPermissaoRepo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
					return &cpm.CategoriaPermissao{}, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao criar atendimento",
			geradorID: &fakeGeradorID{
				id: "atendimento-123",
			},
			repo: &fakeAtendimentoRepository{
				criarFn: func(ctx context.Context, atendimento atd.Atendimento) (*atd.Atendimento, error) {
					return nil, errors.New("erro no banco")
				},
				buscarPorChamadoEAtribuidoIDFn: func(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error) {
					return nil, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return chm.Novo(
						"chamado-123",
						"categoria-123",
						"subcategoria-123",
						"criador-123",
						"Título do chamado",
						"Descrição do chamado",
					)
				},
				atualizarFn: func(ctx context.Context, id string, chamado chm.Chamado) (*chm.Chamado, error) {
					return &chamado, nil
				},
			},
			usuarioRepo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usr.Novo(
						id,
						"Técnico Teste",
						"tecnico.teste@example.com",
						usr.NovoEmail("tecnico@gmail.com"),
						usr.PermTEC,
						nil,
					)
				},
			},
			categoriaPermissaoRepo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
					return &cpm.CategoriaPermissao{}, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAtendimentoService(
				tt.geradorID,
				tt.repo,
				tt.chamadoRepo,
				tt.categoriaPermissaoRepo,
				tt.usuarioRepo,
			)

			params := atd.CriarParams{
				AtribuidoID: "user-123",
				ChamadoID:   "chamado-123",
			}

			_, err := service.Criar(ctx, params)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestAtendimentoService_Atualizar testa o método Atualizar do AtendimentoService.
func TestAtendimentoService_Atualizar(t *testing.T) {
	ctx := context.Background()

	claims := &auth.Claims{
		ID:        "atribuido-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name                   string
		repo                   atd.Repository
		chamadoRepo            chm.Repository
		usuarioRepo            usr.Repository
		categoriaPermissaoRepo cpm.Repository
		params                 atd.AtualizarParams
		wantErr                bool
	}{
		{
			name: "atualizar atendimento com sucesso",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					atd, _ := atd.Novo(
						id,
						"atribuido-123",
						"chamado-123",
					)
					return atd, nil
				},
				atualizarFn: func(ctx context.Context, atendimento atd.Atendimento) (*atd.Atendimento, error) {
					return &atendimento, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					chamado, err := chm.Novo(
						"chamado-123",
						"categoria-123",
						"subcategoria-123",
						"criador-123",
						"Título do chamado",
						"Descrição do chamado",
					)
					if err != nil {
						return nil, err
					}
					return chamado, nil
				},
			},
			usuarioRepo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usr.Novo(
						id,
						"Técnico Teste",
						"tecnico.teste@example.com",
						usr.NovoEmail("tecnico@gmail.com"),
						usr.PermTEC,
						nil,
					)
				},
			},
			categoriaPermissaoRepo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
					cpm, _ := cpm.Novo(
						"categoria-123",
						"atribuido-123",
						usr.PermTEC,
					)
					return cpm, nil
				},
			},
			params: atd.AtualizarParams{
				AtribuidoID: "atribuido-123",
				ChamadoID:   "chamado-123",
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar atendimento",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					return nil, errors.New("erro no banco")
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					chamado, err := chm.Novo(
						"chamado-123",
						"categoria-123",
						"subcategoria-123",
						"criador-123",
						"Título do chamado",
						"Descrição do chamado",
					)
					if err != nil {
						return nil, err
					}
					return chamado, nil
				},
			},
			usuarioRepo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usr.Novo(
						id,
						"Técnico",
						"tec@example.com",
						usr.NovoEmail("tec@example.com"),
						usr.PermTEC,
						nil,
					)
				},
			},
			categoriaPermissaoRepo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
					cpm, _ := cpm.Novo(
						"categoria-123",
						"atribuido-123",
						usr.PermTEC,
					)
					return cpm, nil
				},
			},
			params: atd.AtualizarParams{
				AtribuidoID: "atribuido-123",
				ChamadoID:   "chamado-123",
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar atendimento",
			repo: &fakeAtendimentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*atd.Atendimento, error) {
					atd, _ := atd.Novo(
						id,
						"chamado-123",
						"atribuido-123",
					)
					return atd, nil
				},
				atualizarFn: func(ctx context.Context, atendimento atd.Atendimento) (*atd.Atendimento, error) {
					return nil, errors.New("erro no banco")
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					chamado, err := chm.Novo(
						"chamado-123",
						"categoria-123",
						"subcategoria-123",
						"criador-123",
						"Título do chamado",
						"Descrição do chamado",
					)
					if err != nil {
						return nil, err
					}
					return chamado, nil
				},
			},
			usuarioRepo: &fakeUsuarioRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*usr.Usuario, error) {
					return usr.Novo(
						id,
						"Técnico",
						"tec@example.com",
						usr.NovoEmail("tec@example.com"),
						usr.PermTEC,
						nil,
					)
				},
			},
			categoriaPermissaoRepo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, categoriaID, usuarioID string) (*cpm.CategoriaPermissao, error) {
					cpm, _ := cpm.Novo(
						"categoria-123",
						"atribuido-123",
						usr.PermTEC,
					)
					return cpm, nil
				},
			},
			params: atd.AtualizarParams{
				AtribuidoID: "atribuido-123",
				ChamadoID:   "chamado-123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAtendimentoService(
				nil,
				tt.repo,
				tt.chamadoRepo,
				tt.categoriaPermissaoRepo,
				tt.usuarioRepo,
			)

			_, err := service.Atualizar(ctx, "id-123", tt.params)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("esperava erro, mas recebeu nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestAtendimentoService_Listar testa o método Listar do AtendimentoService.
func TestAtendimentoService_Listar(t *testing.T) {
	ctx := context.Background()

	a, _ := atd.Novo(
		"atendimento-123",
		"atribuido-123",
		"chamado-123",
	)

	atendimentos := []atd.Atendimento{*a}

	tests := []struct {
		name     string
		repo     atd.Repository
		filtro   atd.Filtro
		wantTotal int
		wantErr  bool
	}{
		{
			name: "listar atendimentos com sucesso",
			repo: &fakeAtendimentoRepository{
				listarFn: func(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
					return atendimentos, 1, nil
				},
			},
			filtro:   atd.Filtro{},
			wantTotal: 1,
			wantErr:  false,
		},
		{
			name: "erro ao listar atendimentos",
			repo: &fakeAtendimentoRepository{
				listarFn: func(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
					return nil, 0, errors.New("erro no banco")
				},
			},
			filtro:   atd.Filtro{},
			wantTotal: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAtendimentoService(nil, tt.repo, nil, nil, nil)

			got, total, _, err := service.Listar(ctx, tt.filtro)

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
				t.Fatalf("esperava total '%d', recebeu '%d'", tt.wantTotal, total)
			}

			if len(got) != len(atendimentos) {
				t.Fatalf("esperava %d atendimentos, recebeu %d", len(atendimentos), len(got))
			}
		})
	}
}
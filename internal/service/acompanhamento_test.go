package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

type fakeAcompanhamentoRepository struct {
	criarFn              func(ctx context.Context, a acp.Acompanhamento) (*acp.Acompanhamento, error)
	atualizarFn          func(ctx context.Context, id string, a acp.Acompanhamento) (*acp.Acompanhamento, error)
	buscarPorIDFn        func(ctx context.Context, id string) (*acp.Acompanhamento, error)
	buscarPorChamadoIDFn func(ctx context.Context, id string) ([]acp.Acompanhamento, error)
	deletarFn            func(ctx context.Context, id string) error
	listarFn             func(ctx context.Context, f acp.Filtro) ([]acp.Acompanhamento, int, error)
}

// Criar chama a função fake definida.
func (f *fakeAcompanhamentoRepository) Criar(ctx context.Context, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	return f.criarFn(ctx, a)
}

// Atualizar chama a função fake definida.
func (f *fakeAcompanhamentoRepository) Atualizar(ctx context.Context, id string, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	return f.atualizarFn(ctx, id, a)
}

// BuscarPorID chama a função fake definida.
func (f *fakeAcompanhamentoRepository) BuscarPorID(ctx context.Context, id string) (*acp.Acompanhamento, error) {
	return f.buscarPorIDFn(ctx, id)
}

// BuscarPorChamadoID chama a função fake definida.
func (f *fakeAcompanhamentoRepository) BuscarPorChamadoID(ctx context.Context, id string) ([]acp.Acompanhamento, error) {
	return f.buscarPorChamadoIDFn(ctx, id)
}

// Deletar chama a função fake definida.
func (f *fakeAcompanhamentoRepository) Deletar(ctx context.Context, id string) error {
	return f.deletarFn(ctx, id)
}

// Listar chama a função fake definida.
func (f *fakeAcompanhamentoRepository) Listar(ctx context.Context, ftr acp.Filtro) ([]acp.Acompanhamento, int, error) {
	return f.listarFn(ctx, ftr)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestAcompanhamentoService_BuscarPorID testa o método BuscarPorID do AcompanhamentoService.
func TestAcompanhamentoService_BuscarPorID(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	acompanhamentoOK, _ := acp.Novo(
		"acomp-123",
		"chamado-123",
		"user-123",
		"Conteúdo do acompanhamento",
		usr.PermTEC,
	)

	chamadoOK, _ := chm.Novo(
		"chamado-123",
		"categoria-1",
		"subcategoria-1",
		"user-123",
		"Assunto do chamado",
		"Descrição do chamado",
	)

	atendimentoOK, _ := atd.Novo(
		"atend-123",
		"user-123",
		"chamado-123",
	)

	tests := []struct {
		name            string
		repo            acp.Repository
		chamadoRepo     chm.Repository
		atendimentoRepo atd.Repository
		want            *acp.Acompanhamento
		wantErr         bool
	}{
		{
			name: "sucesso ao buscar acompanhamento por ID",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					if id != "acomp-123" {
						return nil, errors.New("acompanhamento não encontrado")
					}
					return acompanhamentoOK, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return atendimentoOK, nil
				},
			},
			want:    acompanhamentoOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar acompanhamento por ID - acompanhamento não encontrado",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					return nil, errors.New("acompanhamento não encontrado")
				},
			},
			chamadoRepo:     &fakeChamadoRepository{},
			atendimentoRepo: &fakeAtendimentoRepository{},
			want:            nil,
			wantErr:         true,
		},
		{
			name: "erro ao buscar acompanhamento por ID - sem permissão de leitura",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					return acompanhamentoOK, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return nil, errors.New("atendimento não encontrado")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				nil,
				tt.repo,
				tt.chamadoRepo,
				tt.atendimentoRepo,
			)

			got, err := service.BuscarPorID(ctx, "acomp-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got.ID() != tt.want.ID() {
				t.Errorf("BuscarPorID() got = %v, want %v", got.ID(), tt.want.ID())
			}
		})
	}
}

// TestAcompanhamentoService_BuscarPorChamadoID testa o método BuscarPorChamadoID do AcompanhamentoService.
func TestAcompanhamentoService_BuscarPorChamadoID(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	chamadoOK, _ := chm.Novo(
		"chamado-123",
		"categoria-1",
		"subcategoria-1",
		"user-123",
		"Assunto do chamado",
		"Descrição do chamado",
	)

	acompanhamento1, _ := acp.Novo(
		"acomp-123",
		"chamado-123",
		"user-123",
		"Conteúdo do acompanhamento 1",
		usr.PermTEC,
	)

	acompanhamento2, _ := acp.Novo(
		"acomp-124",
		"chamado-123",
		"user-124",
		"Conteúdo do acompanhamento 2",
		usr.PermTEC,
	)

	tests := []struct {
		name            string
		repo            acp.Repository
		chamadoRepo     chm.Repository
		atendimentoRepo atd.Repository
		want            []acp.Acompanhamento
		wantErr         bool
	}{
		{
			name: "sucesso ao buscar acompanhamentos por ID do chamado",
			repo: &fakeAcompanhamentoRepository{
				buscarPorChamadoIDFn: func(ctx context.Context, id string) ([]acp.Acompanhamento, error) {
					if id != "chamado-123" {
						return nil, errors.New("nenhum acompanhamento encontrado")
					}
					return []acp.Acompanhamento{*acompanhamento1, *acompanhamento2}, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return &atd.Atendimento{}, nil
				},
			},
			want:    []acp.Acompanhamento{*acompanhamento1, *acompanhamento2},
			wantErr: false,
		},
		{
			name: "erro ao buscar acompanhamentos por ID do chamado - chamado não encontrado",
			repo: &fakeAcompanhamentoRepository{
				buscarPorChamadoIDFn: func(ctx context.Context, id string) ([]acp.Acompanhamento, error) {
					return nil, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{},
			want:            nil,
			wantErr:         true,
		},
		{
			name: "erro ao buscar acompanhamentos por ID do chamado - sem permissão de leitura",
			repo: &fakeAcompanhamentoRepository{
				buscarPorChamadoIDFn: func(ctx context.Context, id string) ([]acp.Acompanhamento, error) {
					return []acp.Acompanhamento{*acompanhamento1, *acompanhamento2}, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return nil, errors.New("atendimento não encontrado")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				nil,
				tt.repo,
				tt.chamadoRepo,
				tt.atendimentoRepo,
			)

			got, err := service.BuscarPorChamadoID(ctx, "chamado-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorChamadoID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(got) != len(tt.want) {
				t.Errorf("BuscarPorChamadoID() got = %v, want %v", len(got), len(tt.want))
			}
		})
	}
}

// TestAcompanhamentoService_Criar testa o método Criar do AcompanhamentoService.
func TestAcompanhamentoService_Criar(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermTEC.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	chamadoOK, _ := chm.Novo(
		"chamado-123",
		"categoria-1",
		"subcategoria-1",
		"user-123",
		"Assunto do chamado",
		"Descrição do chamado",
	)

	atendimentoOK, _ := atd.Novo(
		"atend-123",
		"user-123",
		"chamado-123",
	)

	tests := []struct {
		name            string
		geradorID       dmn.GeradorID
		repo            acp.Repository
		chamadoRepo     chm.Repository
		atendimentoRepo atd.Repository
		wantErr         bool
	}{
		{
			name: "sucesso ao criar acompanhamento",
			geradorID: &fakeGeradorID{
				id: "acomp-123",
			},
			repo: &fakeAcompanhamentoRepository{
				criarFn: func(ctx context.Context, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
					return &a, nil
				},
			},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return atendimentoOK, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao criar acompanhamento - chamado não encontrado",
			geradorID: &fakeGeradorID{
				id: "acomp-123",
			},
			repo: &fakeAcompanhamentoRepository{},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{},
			wantErr:         true,
		},
		{
			name: "erro ao criar acompanhamento - sem permissão de criação",
			geradorID: &fakeGeradorID{
				id: "acomp-123",
			},
			repo: &fakeAcompanhamentoRepository{},
			chamadoRepo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						return nil, errors.New("chamado não encontrado")
					}
					return chamadoOK, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarPorChamadoEAtribuidoIDFn: func(
					ctx context.Context,
					chamadoID string,
					atribuidoID string,
				) (*atd.Atendimento, error) {
					return nil, errors.New("atendimento não encontrado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				tt.geradorID,
				tt.repo,
				tt.chamadoRepo,
				tt.atendimentoRepo,
			)

			params := acp.CriarParams{
				Conteudo:  "Conteúdo do acompanhamento",
				ChamadoID: "chamado-123",
			}

			_, err := service.Criar(ctx, params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestAcompanhamentoService_Atualizar testa o método Atualizar do AcompanhamentoService.
func TestAcompanhamentoService_Atualizar(t *testing.T) {
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
		name    string
		repo    acp.Repository
		params  acp.AtualizarParams
		wantErr bool
	}{
		{
			name: "sucesso ao atualizar acompanhamento",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					acompanhamento, _ := acp.Novo(
						"acomp-123",
						"chamado-123",
						"user-123",
						"Conteúdo do acompanhamento",
						usr.PermTEC,
					)
					return acompanhamento, nil
				},
				atualizarFn: func(ctx context.Context, id string, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
					return &a, nil
				},
			},
			params: acp.AtualizarParams{
				Conteudo: "Conteúdo atualizado do acompanhamento",
			},
			wantErr: false,
		},
		{
			name: "erro ao atualizar acompanhamento - acompanhamento não encontrado",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					return nil, errors.New("acompanhamento não encontrado")
				},
			},
			params: acp.AtualizarParams{
				Conteudo: "Conteúdo atualizado do acompanhamento",
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar acompanhamento - sem permissão de atualização",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					acompanhamento, _ := acp.Novo(
						"acomp-123",
						"chamado-123",
						"user-456", // Diferente do usuário autenticado
						"Conteúdo do acompanhamento",
						usr.PermTEC,
					)
					return acompanhamento, nil
				},
			},
			params: acp.AtualizarParams{
				Conteudo: "Conteúdo atualizado do acompanhamento",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				nil,
				tt.repo,
				&fakeChamadoRepository{},
				&fakeAtendimentoRepository{},
			)

			_, err := service.Atualizar(ctx, "acomp-123", tt.params)

			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestAcompanhamentoService_Deletar testa o método Deletar do AcompanhamentoService.
func TestAcompanhamentoService_Deletar(t *testing.T) {
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
		name    string
		repo    acp.Repository
		wantErr bool
	}{
		{
			name: "sucesso ao deletar acompanhamento",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					acompanhamento, _ := acp.Novo(
						"acomp-123",
						"chamado-123",
						"user-123",
						"Conteúdo do acompanhamento",
						usr.PermTEC,
					)
					return acompanhamento, nil
				},
				deletarFn: func(ctx context.Context, id string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao deletar acompanhamento - acompanhamento não encontrado",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					return nil, errors.New("acompanhamento não encontrado")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao deletar acompanhamento - sem permissão de deleção",
			repo: &fakeAcompanhamentoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*acp.Acompanhamento, error) {
					acompanhamento, _ := acp.Novo(
						"acomp-123",
						"chamado-123",
						"user-456", // Diferente do usuário autenticado
						"Conteúdo do acompanhamento",
						usr.PermTEC,
					)
					return acompanhamento, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				nil,
				tt.repo,
				&fakeChamadoRepository{},
				&fakeAtendimentoRepository{},
			)

			err := service.Deletar(ctx, "acomp-123")

			if (err != nil) != tt.wantErr {
				t.Errorf("Deletar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestAcompanhamentoService_Listar testa o método Listar do AcompanhamentoService.
func TestAcompanhamentoService_Listar(t *testing.T) {
	ctx := context.Background()

	a, _ := acp.Novo(
		"acomp-123",
		"chamado-123",
		"user-123",
		"Conteúdo do acompanhamento",
		usr.PermTEC,
	)

	acompanhamentos := []acp.Acompanhamento{*a}

	tests := []struct {
		name      string
		repo      acp.Repository
		filtro    acp.Filtro
		wantTotal int
		wantErr   bool
	}{
		{
			name: "sucesso ao listar acompanhamentos",
			repo: &fakeAcompanhamentoRepository{
				listarFn: func(ctx context.Context, f acp.Filtro) ([]acp.Acompanhamento, int, error) {
					if f.Pagina() <= 0 {
						t.Fatalf("esperava pagina normalizada, recebeu %d", f.Pagina())
					}
					if f.Limite() <= 0 {
						t.Fatalf("esperava limite normalizado, recebeu %d", f.Limite())
					}
					return acompanhamentos, len(acompanhamentos), nil
				},
			},
			filtro:    acp.Filtro{},
			wantTotal: len(acompanhamentos),
			wantErr:   false,
		},
		{
			name: "erro ao listar acompanhamentos",
			repo: &fakeAcompanhamentoRepository{
				listarFn: func(ctx context.Context, f acp.Filtro) ([]acp.Acompanhamento, int, error) {
					return nil, 0, errors.New("erro ao listar acompanhamentos")
				},
			},
			filtro:    acp.Filtro{},
			wantTotal: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoAcompanhamentoService(
				nil,
				tt.repo,
				&fakeChamadoRepository{},
				&fakeAtendimentoRepository{},
			)

			_, gotTotal, _, err := service.Listar(ctx, tt.filtro)

			if (err != nil) != tt.wantErr {
				t.Errorf("Listar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotTotal != tt.wantTotal {
				t.Errorf("Listar() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}

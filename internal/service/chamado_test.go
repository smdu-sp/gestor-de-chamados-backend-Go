package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeChamadoRepository é uma implementação falsa de chamado.Repository para testes.
type fakeChamadoRepository struct {
	buscarPorIDFn func(ctx context.Context, id string) (*chm.Chamado, error)
	criarFn       func(ctx context.Context, c chm.Chamado) (*chm.Chamado, error)
	atualizarFn   func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error)
	listarFn      func(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error)
}

// BuscarPorID chama a função fake definida.
func (f *fakeChamadoRepository) BuscarPorID(ctx context.Context, id string) (*chm.Chamado, error) {
	return f.buscarPorIDFn(ctx, id)
}

// Criar chama a função fake definida.
func (f *fakeChamadoRepository) Criar(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
	return f.criarFn(ctx, c)
}

// Atualizar chama a função fake definida.
func (f *fakeChamadoRepository) Atualizar(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
	return f.atualizarFn(ctx, id, c)
}

// Listar chama a função fake definida.
func (f *fakeChamadoRepository) Listar(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestChamadoService_BuscarPorID testa o método BuscarPorID do ChamadoService.
func TestChamadoService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	chamadoOK := &chm.Chamado{} // não importa o conteúdo aqui

	tests := []struct {
		name    string
		repo    chm.Repository
		want    *chm.Chamado
		wantErr bool
	}{
		{
			name: "buscar chamado com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					if id != "chamado-123" {
						t.Errorf("ID inesperado: got %v, want %v", id, "chamado-123")
					}
					return chamadoOK, nil
				},
			},
			want:    chamadoOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar chamado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "chamado não encontrado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, nil, nil)

			// Act
			got, err := service.BuscarPorID(ctx, "chamado-123")

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

// TestChamadoService_Criar testa o método Criar do ChamadoService.
func TestChamadoService_Criar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name             string
		geradorID        dmn.GeradorID
		repo             chm.Repository
		categoriaRepo    ctg.Repository
		subcategoriaRepo subc.Repository
		wantErr          bool
	}{
		{
			name: "criar chamado com sucesso",
			geradorID: &fakeGeradorID{
				id: "chamado-123",
			},
			repo: &fakeChamadoRepository{
				criarFn: func(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
					if c.ID() != "chamado-123" {
						t.Errorf("ID inesperado: got %v, want %v", c.ID(), "chamado-123")
					}
					return &c, nil
				},
			},
			categoriaRepo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*ctg.Categoria, error) {
					categoria, _ := ctg.Novo(id, "Categoria Teste")
					categoria.Ativar()
					return categoria, nil
				},
			},
			subcategoriaRepo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*subc.Subcategoria, error) {
					subcategoria, _ := subc.Novo(id, "Subcategoria Teste", "cat-123")
					subcategoria.Ativar()
					return subcategoria, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errors.New("falha ao gerar ID"),
			},
			repo: &fakeChamadoRepository{
				criarFn: func(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
					return nil, nil
				},
			},
			categoriaRepo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*ctg.Categoria, error) {
					categoria, _ := ctg.Novo(id, "Categoria Teste")
					categoria.Ativar()
					return categoria, nil
				},
			},
			subcategoriaRepo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*subc.Subcategoria, error) {
					subcategoria, _ := subc.Novo(id, "Subcategoria Teste", "cat-123")
					subcategoria.Ativar()
					return subcategoria, nil
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao criar chamado no repositório",
			geradorID: &fakeGeradorID{
				id: "chamado-123",
			},
			repo: &fakeChamadoRepository{
				criarFn: func(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
					return nil, errors.New("erro ao criar chamado")
				},
			},
			categoriaRepo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*ctg.Categoria, error) {
					categoria, _ := ctg.Novo(id, "Categoria Teste")
					categoria.Ativar()
					return categoria, nil
				},
			},
			subcategoriaRepo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*subc.Subcategoria, error) {
					subcategoria, _ := subc.Novo(id, "Subcategoria Teste", "cat-123")
					subcategoria.Ativar()
					return subcategoria, nil
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(tt.geradorID, tt.repo, tt.categoriaRepo, tt.subcategoriaRepo, nil, nil, nil)

			params := chm.CriarParams{
				Titulo:         "Chamado de teste",
				Descricao:      "Descrição do chamado de teste",
				CategoriaID:    "cat-123",
				SubcategoriaID: "subcat-123",
				CriadorID:      "user-123",
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

// TestChamadoService_Atualizar testa o método Atualizar do ChamadoService.
func TestChamadoService_Atualizar(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermUSR.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name    string
		repo    chm.Repository
		params  chm.AtualizarParams
		wantErr bool
	}{
		{
			name: "atualizar chamado com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					if c.Descricao() != "Descrição atualizada" {
						t.Fatalf("Descrição inesperada: got %v, want %v", c.Descricao(), "Descrição atualizada")
					}
					return &c, nil
				},
			},
			params: chm.AtualizarParams{
				Descricao: ptrString("Descrição atualizada"),
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar chamado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			params: chm.AtualizarParams{
				Descricao: ptrString("Descrição atualizada"),
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar chamado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					return nil, errors.New("erro ao atualizar chamado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, nil, nil)

			// Act
			_, err := service.Atualizar(ctx, "chamado-123", tt.params)

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

// TestChamadoService_ArquivarDesarquivar testa os métodos Arquivar e Desarquivar do ChamadoService.
func TestChamadoService_ArquivarDesarquivar(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermUSR.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name    string
		repo    chm.Repository
		acao    string // "arquivar" ou "desarquivar"
		wantErr bool
	}{
		{
			name: "arquivar chamado com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					if !c.Arquivado() {
						t.Fatalf("Esperava chamado arquivado, mas não está")
					}
					return &c, nil
				},
			},
			acao:    "arquivar",
			wantErr: false,
		},
		{
			name: "desarquivar chamado com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					c.Arquivar()
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					if c.Arquivado() {
						t.Fatalf("Esperava chamado desarquivado, mas está arquivado")
					}
					return &c, nil
				},
			},
			acao:    "desarquivar",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, nil, nil)

			// Act
			var err error
			switch tt.acao {
			case "arquivar":
				_, err = service.Arquivar(ctx, "chamado-123")

			case "desarquivar":
				_, err = service.Desarquivar(ctx, "chamado-123")

			default:
				t.Fatalf("Ação desconhecida: %s", tt.acao)
			}

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

// TestChamadoService_AtualizarStatus testa o método AtualizarStatus do ChamadoService.
func TestChamadoService_AtualizarStatus(t *testing.T) {
	ctx := context.Background()
	claims := &auth.Claims{
		ID:        "user-123",
		Login:     "testuser",
		Nome:      "Test User",
		Email:     "test@example.com",
		Permissao: usr.PermUSR.String(),
	}
	ctx = auth.ContextoComClaims(ctx, claims)

	tests := []struct {
		name       string
		repo       chm.Repository
		novoStatus chm.StatusChamado
		wantErr    bool
	}{
		{
			name: "atualizar status com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					if c.Status() != chm.StatusAberto {
						t.Fatalf("Status inesperado: got %v, want %v", c.Status(), chm.StatusAberto)
					}
					return &c, nil
				},
			},
			novoStatus: chm.StatusAtribuido,
			wantErr:    false,
		},
		{
			name: "erro ao buscar chamado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			novoStatus: chm.StatusAtribuido,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, nil, nil)

			params := chm.AtualizarStatusParams{
				Status: chm.StatusAberto,
			}

			// Act
			_, err := service.AtualizarStatus(ctx, "chamado-123", params)

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

// TestChamadoService_AtualizarSolucao testa o método AtualizarSolucao do ChamadoService.
func TestChamadoService_AtualizarSolucao(t *testing.T) {
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
		name            string
		repo            chm.Repository
		atendimentoRepo atd.Repository
		solucao         string
		wantErr         bool
	}{
		{
			name: "atualizar solução com sucesso",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					c, _ := chm.Novo(
						"chamado-123",
						"cat-123",
						"subcat-123",
						"user-123",
						"Chamado de teste",
						"Descrição do chamado de teste",
					)
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
					if c.Solucao() == nil || *c.Solucao() != "Solução atualizada" {
						t.Fatalf("Solução inesperada: got %v, want %v", c.Solucao(), "Solução atualizada")
					}
					return &c, nil
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarFn: func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error) {
					atendimento, _ := atd.Novo("atendimento-123", chamadoID, usuarioID)
					return atendimento, nil
				},
			},
			solucao: "Solução atualizada",
			wantErr: false,
		},
		{
			name: "erro ao buscar chamado",
			repo: &fakeChamadoRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*chm.Chamado, error) {
					return nil, errors.New("chamado não encontrado")
				},
			},
			atendimentoRepo: &fakeAtendimentoRepository{
				buscarFn: func(ctx context.Context, chamadoID, usuarioID string) (*atd.Atendimento, error) {
					atendimento, _ := atd.Novo("atendimento-123", chamadoID, usuarioID)
					return atendimento, nil
				},
			},
			solucao: "Solução atualizada",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, tt.atendimentoRepo, nil)

			params := chm.AtualizarSolucaoParams{
				Solucao: tt.solucao,
			}

			// Act
			_, err := service.AtualizarSolucao(ctx, "chamado-123", params)

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

// TestChamadoService_Listar testa o método Listar do ChamadoService.
func TestChamadoService_Listar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    chm.Repository
		want    []chm.Chamado
		wantErr bool
	}{
		{
			name: "listar chamados com sucesso",
			repo: &fakeChamadoRepository{
				listarFn: func(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
					chamados := []chm.Chamado{}
					return chamados, len(chamados), nil
				},
			},
			want:    []chm.Chamado{},
			wantErr: false,
		},
		{
			name: "erro ao listar chamados",
			repo: &fakeChamadoRepository{
				listarFn: func(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
					return nil, 0, errors.New("erro ao listar chamados")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoChamadoService(nil, tt.repo, nil, nil, nil, nil, nil)

			filtro := chm.Filtro{}

			// Act
			got, _, _, err := service.Listar(ctx, filtro)

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

			if len(got) != len(tt.want) {
				t.Fatalf("esperava %d chamados, recebeu %d", len(tt.want), len(got))
			}
		})
	}
}

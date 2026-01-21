package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeCategoriaRepository é uma implementação falsa de categoria.Repository para testes.
type fakeCategoriaRepository struct {
	buscarPorIDFn func(ctx context.Context, id string) (*categoria.Categoria, error)
	atualizarFn   func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error)
	buscarNomeFn  func(ctx context.Context, nome string) (*categoria.Categoria, error)
	criarFn       func(ctx context.Context, c categoria.Categoria) (*categoria.Categoria, error)
	listarFn      func(ctx context.Context, filtro categoria.Filtro) ([]categoria.Categoria, int, error)
}

// BuscarPorID chama a função fake definida.
func (f *fakeCategoriaRepository) BuscarPorID(ctx context.Context, id string) (*categoria.Categoria, error) {
	return f.buscarPorIDFn(ctx, id)
}

// Atualizar chama a função fake definida.
func (f *fakeCategoriaRepository) Atualizar(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
	return f.atualizarFn(ctx, id, c)
}

// BuscarPorNome chama a função fake definida.
func (f *fakeCategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*categoria.Categoria, error) {
	return f.buscarNomeFn(ctx, nome)
}

// Criar chama a função fake definida.
func (f *fakeCategoriaRepository) Criar(ctx context.Context, c categoria.Categoria) (*categoria.Categoria, error) {
	return f.criarFn(ctx, c)
}

// Listar chama a função fake definida.
func (f *fakeCategoriaRepository) Listar(ctx context.Context, filtro categoria.Filtro) ([]categoria.Categoria, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestCategoriaService_BuscarPorID testa o método BuscarPorID do CategoriaService.
func TestCategoriaService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	categoriaOK := &categoria.Categoria{}

	tests := []struct {
		name    string
		repo    categoria.Repository
		want    *categoria.Categoria
		wantErr bool
	}{
		{
			name: "buscar categoria com sucesso",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					if id != "qualquer-id" {
						t.Errorf("esperava id 'qualquer-id', recebeu '%s'", id)
					}
					return categoriaOK, nil
				},
			},
			want:    categoriaOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "categoria não encontrada",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return nil, errors.New("categoria não encontrada")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := NovoCategoriaService(nil, tt.repo)

			got, err := svc.BuscarPorID(ctx, "qualquer-id")
			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuscarPorID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCategoriaService_BuscarPorNome testa o método BuscarPorNome do CategoriaService.
func TestCategoriaService_BuscarPorNome(t *testing.T) {
	ctx := context.Background()

	categoriaOK := &categoria.Categoria{}

	tests := []struct {
		name    string
		repo    categoria.Repository
		want    *categoria.Categoria
		wantErr bool
	}{
		{
			name: "buscar categoria com sucesso",
			repo: &fakeCategoriaRepository{
				buscarNomeFn: func(ctx context.Context, nome string) (*categoria.Categoria, error) {
					if nome != "qualquer-nome" {
						t.Errorf("esperava nome 'qualquer-nome', recebeu '%s'", nome)
					}
					return categoriaOK, nil
				},
			},
			want:    categoriaOK,
			wantErr: false,
		},
		{
			name: "erro ao buscar categoria",
			repo: &fakeCategoriaRepository{
				buscarNomeFn: func(ctx context.Context, nome string) (*categoria.Categoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
		{
			name: "categoria não encontrada",
			repo: &fakeCategoriaRepository{
				buscarNomeFn: func(ctx context.Context, nome string) (*categoria.Categoria, error) {
					return nil, errors.New("categoria não encontrada")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(nil, tt.repo)

			got, err := service.BuscarPorNome(ctx, "qualquer-nome")
			if (err != nil) != tt.wantErr {
				t.Errorf("BuscarPorNome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("BuscarPorNome() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCategoriaService_Criar testa o método Criar do CategoriaService.
func TestCategoriaService_Criar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repo      categoria.Repository
		wantErr   bool
	}{
		{
			name: "criar categoria com sucesso",
			geradorID: &fakeGeradorID{
				id: "novo-id",
			},
			repo: &fakeCategoriaRepository{
				criarFn: func(ctx context.Context, c categoria.Categoria) (*categoria.Categoria, error) {
					if c.ID() != "novo-id" {
						t.Fatalf("esperava id 'novo-id', recebeu '%s'", c.ID())
					}
					if c.Nome() != "Nova Categoria" {
						t.Fatalf("esperava nome 'Nova Categoria', recebeu '%s'", c.Nome())
					}
					return &c, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao gerar ID",
			geradorID: &fakeGeradorID{
				err: errors.New("erro ao gerar ID"),
			},
			repo: &fakeCategoriaRepository{
				criarFn: func(ctx context.Context, c categoria.Categoria) (*categoria.Categoria, error) {
					t.Fatalf("não esperava chamar o criar do repositório")
					return nil, nil
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao criar categoria",
			geradorID: &fakeGeradorID{
				id: "novo-id",
			},
			repo: &fakeCategoriaRepository{
				criarFn: func(ctx context.Context, c categoria.Categoria) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao criar categoria")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(tt.geradorID, tt.repo)

			criarParams := categoria.CriarParams{
				Nome: "Nova Categoria",
			}

			_, err := service.Criar(ctx, criarParams)

			if tt.wantErr && err == nil {
				t.Errorf("esperava erro, mas não recebeu nenhum")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}

// TestCategoriaService_Atualizar testa o método Atualizar do CategoriaService.
func TestCategoriaService_Atualizar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    categoria.Repository
		params categoria.AtualizarParams
		wantErr bool
	}{
		{
			name: "atualizar categoria com sucesso",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					c, _ := categoria.Novo(id, "Categoria Antiga")
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					if c.Nome() != "Categoria Atualizada" {
						t.Fatalf("esperava nome 'Categoria Atualizada', recebeu '%s'", c.Nome())
					}
					return &c, nil
				},
			},
			params: categoria.AtualizarParams{
				Nome: ptr("Categoria Atualizada"),
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao buscar categoria")
				},
			},
			params: categoria.AtualizarParams{
				Nome: ptr("Categoria Atualizada"),
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					c, _ := categoria.Novo(id, "Categoria Antiga")
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao atualizar categoria")
				},
			},
			params: categoria.AtualizarParams{
				Nome: ptr("Categoria Atualizada"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(nil, tt.repo)

			_, err := service.Atualizar(ctx, "categoria-id", tt.params)

			if tt.wantErr && err == nil {
				t.Errorf("esperava erro, mas não recebeu nenhum")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}

// TestCategoriaService_Ativar testa o método Atualizar do CategoriaService.
func TestCategoriaService_Ativar(t *testing.T) {
	ctx := context.Background()

	categoriaAtiva, _ := categoria.Novo("categoria-id", "Categoria Ativa")

	tests := []struct {
		name    string
		repo    categoria.Repository
		wantErr bool
	}{
		{
			name: "ativar categoria com sucesso",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return categoriaAtiva, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					return &c, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao buscar categoria")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					c, _ := categoria.Novo(id, "Categoria Inativa")
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao atualizar categoria")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(nil, tt.repo)

			_, err := service.Ativar(ctx, "categoria-id")

			if tt.wantErr && err == nil {
				t.Errorf("esperava erro, mas não recebeu nenhum")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}

// TestCategoriaService_Desativar testa o método Desativar do CategoriaService.
func TestCategoriaService_Desativar(t *testing.T) {
	ctx := context.Background()

	categoriaAtiva, _ := categoria.Novo("categoria-id", "Categoria Ativa")

	tests := []struct {
		name    string
		repo    categoria.Repository
		wantErr bool
	}{
		{
			name: "desativar categoria com sucesso",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return categoriaAtiva, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					return &c, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao buscar categoria")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao atualizar categoria",
			repo: &fakeCategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*categoria.Categoria, error) {
					c, _ := categoria.Novo(id, "Categoria Ativa")
					return c, nil
				},
				atualizarFn: func(ctx context.Context, id string, c categoria.Categoria) (*categoria.Categoria, error) {
					return nil, errors.New("erro ao atualizar categoria")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(nil, tt.repo)

			_, err := service.Desativar(ctx, "categoria-id")

			if tt.wantErr && err == nil {
				t.Errorf("esperava erro, mas não recebeu nenhum")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("não esperava erro, mas recebeu: %v", err)
			}
		})
	}
}

// TestCategoriaService_Listar testa o método Listar do CategoriaService.
func TestCategoriaService_Listar(t *testing.T) {
	ctx := context.Background()

	c, _ := categoria.Novo("cat-1", "Categoria 1")
	c2, _ := categoria.Novo("cat-2", "Categoria 2")

	categorias := []categoria.Categoria{*c, *c2}

	tests := []struct {
		name       string
		repo       categoria.Repository
		filtro categoria.Filtro
		wantTotal  int
		wantErr    bool
	}{
		{
			name: "listar categorias com sucesso e normalizar filtro",
			repo: &fakeCategoriaRepository{
				listarFn: func(ctx context.Context, filtro categoria.Filtro) ([]categoria.Categoria, int, error) {
					if filtro.Pagina() <= 0 {
						t.Errorf("esperava página maior que 0, recebeu %d", filtro.Pagina())
					}
					if filtro.Limite() <= 0 {
						t.Errorf("esperava limite maior que 0, recebeu %d", filtro.Limite())
					}
					return categorias, len(categorias), nil
				},
			},
			wantTotal: len(categorias),
			wantErr:   false,
		},
		{
			name: "erro ao listar categorias",
			repo: &fakeCategoriaRepository{
				listarFn: func(ctx context.Context, filtro categoria.Filtro) ([]categoria.Categoria, int, error) {
					return nil, 0, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaService(nil, tt.repo)

			categorias, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			if tt.wantErr && err == nil {
				t.Errorf("esperava erro, mas não recebeu nenhum")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("não esperava erro, mas recebeu: %v", err)
			}

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("esperava total %d, recebeu %d", tt.wantTotal, total)
				}
				if len(categorias) != tt.wantTotal {
					t.Errorf("esperava %d categorias, recebeu %d", tt.wantTotal, len(categorias))
				}
				if filtroRetornado.Pagina() <= 0 {
					t.Errorf("esperava página maior que 0, recebeu %d", filtroRetornado.Pagina())
				}
				if filtroRetornado.Limite() <= 0 {
					t.Errorf("esperava limite maior que 0, recebeu %d", filtroRetornado.Limite())
				}
			}
		})
	}
}
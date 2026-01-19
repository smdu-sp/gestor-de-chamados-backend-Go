package service

import (
	"context"
	"errors"
	"testing"

	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	sub "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeSubcategoriaRepository é uma implementação fake de sub.Repository para testes.
type fakeSubcategoriaRepository struct {
	buscarPorIDFn   func(ctx context.Context, id string) (*sub.Subcategoria, error)
	buscarPorNomeFn func(ctx context.Context, nome string) (*sub.Subcategoria, error)
	criarFn         func(ctx context.Context, subcategoria sub.Subcategoria) (*sub.Subcategoria, error)
	atualizarFn     func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error)
	listarFn				func(ctx context.Context, f sub.Filtro) ([]sub.Subcategoria, int, error)
}

// BuscarPorID chama a função fake definida em buscarPorIDFn.
func (f *fakeSubcategoriaRepository) BuscarPorID(ctx context.Context, id string) (*sub.Subcategoria, error) {
	return f.buscarPorIDFn(ctx, id)
}

// BuscarPorNome chama a função fake definida em buscarPorNomeFn.
func (f *fakeSubcategoriaRepository) BuscarPorNome(ctx context.Context, nome string) (*sub.Subcategoria, error) {
	return f.buscarPorNomeFn(ctx, nome)
}

// Criar chama a função fake definida em criarFn.
func (f *fakeSubcategoriaRepository) Criar(ctx context.Context, s sub.Subcategoria) (*sub.Subcategoria, error) {
	return f.criarFn(ctx, s)
}

func (f *fakeSubcategoriaRepository) Atualizar(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
	return f.atualizarFn(ctx, id, c)
}

func (f *fakeSubcategoriaRepository) Listar(ctx context.Context, filtro sub.Filtro) ([]sub.Subcategoria, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestSubcategoriaService_BuscarPorID testa o método BuscarPorID do SubcategoriaService.
func TestSubcategoriaService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	subc, _ := sub.Novo("sub-1", "Hardware", "cat-1")

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repo      sub.Repository
		wantErr   bool
	}{
		{
			name: "buscar subcategoria com sucesso",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subc, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar subcategoria",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoSubcategoriaService(tt.geradorID, tt.repo)

			// Act
			_, err := service.BuscarPorID(ctx, "sub-1")

			// Assert
			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_BuscarPorNome testa o método BuscarPorNome do SubcategoriaService.
func TestSubcategoriaService_BuscarPorNome(t *testing.T) {
	ctx := context.Background()

	subc, _ := sub.Novo("sub-1", "Hardware", "cat-1")

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repo      sub.Repository
		wantErr   bool
	}{
		{
			name: "buscar subcategoria por nome com sucesso",
			repo: &fakeSubcategoriaRepository{
				buscarPorNomeFn: func(ctx context.Context, nome string) (*sub.Subcategoria, error) {
					return subc, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar subcategoria por nome",
			repo: &fakeSubcategoriaRepository{
				buscarPorNomeFn: func(ctx context.Context, nome string) (*sub.Subcategoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoSubcategoriaService(tt.geradorID, tt.repo)

			// Act
			_, err := service.BuscarPorNome(ctx, "Hardware")

			// Assert
			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_Criar testa o método Criar do SubcategoriaService.
func TestSubcategoriaService_Criar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		geradorID dmn.GeradorID
		repo      sub.Repository
		params    sub.CriarParams
		wantErr   bool
	}{
		{
			name: "criar subcategoria com sucesso",
			geradorID: &fakeGeradorID{
				id: "sub-123",
			},
			repo: &fakeSubcategoriaRepository{
				criarFn: func(ctx context.Context, c sub.Subcategoria) (*sub.Subcategoria, error) {
					return &c, nil
				},
			},
			params: sub.CriarParams{
				Nome:        "Hardware",
				CategoriaID: "cat-1",
			},
			wantErr: false,
		},
		{
			name: "erro ao gerar id",
			geradorID: &fakeGeradorID{
				err: errors.New("erro ao gerar id"),
			},
			repo:    &fakeSubcategoriaRepository{},
			params:  sub.CriarParams{},
			wantErr: true,
		},
		{
			name: "erro de validacao do dominio",
			geradorID: &fakeGeradorID{
				id: "sub-123",
			},
			repo: &fakeSubcategoriaRepository{},
			params: sub.CriarParams{
				Nome:        "a", // inválido
				CategoriaID: "",
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar no repositorio",
			geradorID: &fakeGeradorID{
				id: "sub-123",
			},
			repo: &fakeSubcategoriaRepository{
				criarFn: func(ctx context.Context, c sub.Subcategoria) (*sub.Subcategoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			params: sub.CriarParams{
				Nome:        "Hardware",
				CategoriaID: "cat-1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			service := NovoSubcategoriaService(tt.geradorID, tt.repo)

			// Act
			_, err := service.Criar(ctx, tt.params)

			// Assert
			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_Atualizar testa o método Atualizar do SubcategoriaService.
func TestSubcategoriaService_Atualizar(t *testing.T) {
	ctx := context.Background()

	nomeAtualizado := "Infraestrutura"
	statusDesativado := false

	tests := []struct {
		name    string
		repo    sub.Repository
		params  sub.AtualizarParams
		wantErr bool
	}{
		{
			name: "atualizar subcategoria com sucesso",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subcategoriaValida(), nil
				},
				atualizarFn: func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
					return &c, nil
				},
			},
			params: sub.AtualizarParams{
				Nome:   &nomeAtualizado,
				Status: &statusDesativado,
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar subcategoria",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return nil, errors.New("nao encontrada")
				},
			},
			params:  sub.AtualizarParams{},
			wantErr: true,
		},
		{
			name: "erro de validacao ao atualizar dados",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subcategoriaValida(), nil
				},
			},
			params: sub.AtualizarParams{
				Nome: ptr("a"), // inválido
			},
			wantErr: true,
		},
		{
			name: "erro ao salvar atualizacao",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subcategoriaValida(), nil
				},
				atualizarFn: func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
					return nil, errors.New("erro no banco")
				},
			},
			params: sub.AtualizarParams{
				Nome: &nomeAtualizado,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoSubcategoriaService(nil, tt.repo)

			_, err := service.Atualizar(ctx, "sub-1", tt.params)

			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_Desativar testa o método Desativar do SubcategoriaService.
func TestSubcategoriaService_Desativar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    sub.Repository
		wantErr bool
	}{
		{
			name: "desativar subcategoria com sucesso",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subcategoriaValida(), nil
				},
				atualizarFn: func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
					if c.Status() {
						t.Fatal("esperava subcategoria desativada")
					}
					return &c, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar subcategoria",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return nil, errors.New("erro banco")
				},
			},
			wantErr: true,
		},
		{
			name: "erro ao persistir desativacao",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					return subcategoriaValida(), nil
				},
				atualizarFn: func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
					return nil, errors.New("erro banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoSubcategoriaService(nil, tt.repo)

			_, err := service.Desativar(ctx, "sub-1")

			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_Ativar testa o método Ativar do SubcategoriaService.
func TestSubcategoriaService_Ativar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    sub.Repository
		wantErr bool
	}{
		{
			name: "ativar subcategoria com sucesso",
			repo: &fakeSubcategoriaRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*sub.Subcategoria, error) {
					subc := subcategoriaValida()
					subc.Desativar()
					return subc, nil
				},
				atualizarFn: func(ctx context.Context, id string, c sub.Subcategoria) (*sub.Subcategoria, error) {
					if !c.Status() {
						t.Fatal("esperava subcategoria ativada")
					}
					return &c, nil
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoSubcategoriaService(nil, tt.repo)

			_, err := service.Ativar(ctx, "sub-1")

			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}
		})
	}
}

// TestSubcategoriaService_Listar testa o método Listar do SubcategoriaService.
func TestSubcategoriaService_Listar(t *testing.T) {
	ctx := context.Background()

	filtro := sub.Filtro{}
	filtro.Normalizar()

	tests := []struct {
		name    string
		repo    sub.Repository
		filtro  sub.Filtro
		wantErr bool
	}{
		{
			name: "listar subcategorias com sucesso",
			filtro: filtro,
			repo: &fakeSubcategoriaRepository{
				listarFn: func(ctx context.Context, f sub.Filtro) ([]sub.Subcategoria, int, error) {
					// Assert implícito: filtro já normalizado
					if f.Pagina() <= 0 {
						t.Fatal("esperava pagina normalizada")
					}
					if f.Limite() <= 0 {
						t.Fatal("esperava limite normalizado")
					}

					return []sub.Subcategoria{
						*subcategoriaValida(),
					}, 1, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao listar subcategorias",
			filtro: filtro,
			repo: &fakeSubcategoriaRepository{
				listarFn: func(ctx context.Context, f sub.Filtro) ([]sub.Subcategoria, int, error) {
					return nil, 0, errors.New("erro banco")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoSubcategoriaService(nil, tt.repo)

			subc, total, filtro, err := service.Listar(ctx, tt.filtro)

			if tt.wantErr && err == nil {
				t.Fatal("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr {
				if err != nil {
					t.Fatalf("erro inesperado: %v", err)
				}
				if total == 0 {
					t.Fatal("esperava total maior que zero")
				}
				if len(subc) == 0 {
					t.Fatal("esperava lista preenchida")
				}
				if filtro.Pagina() <= 0 || filtro.Limite() <= 0 {
					t.Fatal("esperava filtro normalizado no retorno")
				}
			}
		})
	}
}


// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// subcategoriaValida retorna uma subcategoria válida para uso em testes.
func subcategoriaValida() *sub.Subcategoria {
	subc, _ := sub.Novo(
		"sub-1",
		"Hardware",
		"cat-1",
	)
	return subc
}

package service

import (
	"context"
	"errors"
	"testing"

	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeCategoriaPermissaoRepository é uma implementação falsa do repositório de categorias de permissão, usada para testes.
type fakeCategoriaPermissaoRepository struct {
	criarFn       func(ctx context.Context, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error)
	buscarPorIDFn func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error)
	atualizarFn   func(ctx context.Context, categoriaID, usuarioID string, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error)
	deletarFn     func(ctx context.Context, ctgID, usrID string) error
	listarFn      func(ctx context.Context, f cpm.Filtro) ([]cpm.CategoriaPermissao, int, error)
}

// Criar chama a função de criação definida na struct falsa.
func (f *fakeCategoriaPermissaoRepository) Criar(ctx context.Context, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	return f.criarFn(ctx, cpm)
}

// BuscarPorID chama a função de busca por ID definida na struct falsa.
func (f *fakeCategoriaPermissaoRepository) BuscarPorID(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
	return f.buscarPorIDFn(ctx, ctgID, usrID)
}

// Atualizar chama a função de atualização definida na struct falsa.
func (f *fakeCategoriaPermissaoRepository) Atualizar(ctx context.Context, categoriaID, usuarioID string, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
	return f.atualizarFn(ctx, categoriaID, usuarioID, cpm)
}

// Deletar chama a função de deleção definida na struct falsa.
func (f *fakeCategoriaPermissaoRepository) Deletar(ctx context.Context, ctgID, usrID string) error {
	return f.deletarFn(ctx, ctgID, usrID)
}

// Listar chama a função de listagem definida na struct falsa.
func (f *fakeCategoriaPermissaoRepository) Listar(ctx context.Context, filtro cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestCategoriaPermissaoService_Criar testa o método Criar do CategoriaPermissaoService.
func TestCategoriaPermissaoService_Criar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		repo        cpm.Repository
		criarParams cpm.CriarParams
		wantErr     bool
	}{
		{
			name: "Criação bem-sucedida de categoria de permissão",
			repo: &fakeCategoriaPermissaoRepository{
				criarFn: func(ctx context.Context, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
					return &cpm, nil
				},
			},
			criarParams: cpm.CriarParams{
				CategoriaID: "ctg-123",
				UsuarioID:   "usr-456",
				Permissao:   usr.PermADM,
			},
			wantErr: false,
		},
		{
			name: "Falha na criação devido a categoria/permissão já existente",
			repo: &fakeCategoriaPermissaoRepository{
				criarFn: func(ctx context.Context, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("categoria/permissão já existe")
				},
			},
			criarParams: cpm.CriarParams{
				CategoriaID: "ctg-123",
				UsuarioID:   "usr-456",
				Permissao:   usr.PermADM,
			},
			wantErr: true,
		},
		{
			name: "Falha na criação devido a erro no banco de dados",
			repo: &fakeCategoriaPermissaoRepository{
				criarFn: func(ctx context.Context, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("erro no banco de dados")
				},
			},
			criarParams: cpm.CriarParams{
				CategoriaID: "ctg-123",
				UsuarioID:   "usr-456",
				Permissao:   usr.PermADM,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaPermissaoService(tt.repo)

			_, err := service.Criar(ctx, tt.criarParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("Criar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestCategoriaPermissaoService_BuscarPorID testa o método BuscarPorID do CategoriaPermissaoService.
func TestCategoriaPermissaoService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	categoriaPermissaoOK := &cpm.CategoriaPermissao{}

	tests := []struct {
		name    string
		repo    cpm.Repository
		want    *cpm.CategoriaPermissao
		wantErr bool
	}{
		{
			name: "Busca bem-sucedida de categoria de permissão por ID",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return categoriaPermissaoOK, nil
				},
			},
			want:    categoriaPermissaoOK,
			wantErr: false,
		},
		{
			name: "Falha na busca devido a categoria/permissão não encontrada",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("categoria/permissão não encontrada")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Falha na busca devido a erro no banco de dados",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("erro no banco de dados")
				},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaPermissaoService(tt.repo)

			got, err := service.BuscarPorID(ctx, "ctg-123", "usr-456")
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

// TestCategoriaPermissaoService_Atualizar testa o método Atualizar do CategoriaPermissaoService.
func TestCategoriaPermissaoService_Atualizar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name            string
		repo            cpm.Repository
		atualizarParams cpm.AtualizarParams
		wantErr         bool
	}{
		{
			name: "Atualização bem-sucedida de categoria de permissão",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					cpm, _ := cpm.Novo(ctgID, usrID, usr.PermADM)
					return cpm, nil
				},
				atualizarFn: func(ctx context.Context, categoriaID, usuarioID string, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
					return &cpm, nil
				},
			},
			atualizarParams: cpm.AtualizarParams{
				Permissao: usr.PermADM,
			},
			wantErr: false,
		},
		{
			name: "Falha na atualização devido a categoria/permissão não encontrada",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("categoria/permissão não encontrada")
				},
			},
			atualizarParams: cpm.AtualizarParams{
				Permissao: usr.PermADM,
			},
			wantErr: true,
		},
		{
			name: "Falha na atualização devido a erro no banco de dados",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					cpm, _ := cpm.Novo(ctgID, usrID, usr.PermADM)
					return cpm, nil
				},
				atualizarFn: func(ctx context.Context, categoriaID, usuarioID string, cpm cpm.CategoriaPermissao) (*cpm.CategoriaPermissao, error) {
					return nil, errors.New("erro no banco de dados")
				},
			},
			atualizarParams: cpm.AtualizarParams{
				Permissao: usr.PermADM,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaPermissaoService(tt.repo)

			_, err := service.Atualizar(ctx, "ctg-123", "usr-456", tt.atualizarParams)
			if (err != nil) != tt.wantErr {
				t.Errorf("Atualizar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestCategoriaPermissaoService_Deletar testa o método Deletar do CategoriaPermissaoService.
func TestCategoriaPermissaoService_Deletar(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		repo    cpm.Repository
		wantErr bool
	}{
		{
			name: "Deleção bem-sucedida de categoria de permissão",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return &cpm.CategoriaPermissao{}, nil
				},
				deletarFn: func(ctx context.Context, ctgID, usrID string) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "Falha na deleção devido a categoria/permissão não encontrada",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return &cpm.CategoriaPermissao{}, nil
				},
				deletarFn: func(ctx context.Context, ctgID, usrID string) error {
					return errors.New("categoria/permissão não encontrada")
				},
			},
			wantErr: true,
		},
		{
			name: "Falha na deleção devido a erro no banco de dados",
			repo: &fakeCategoriaPermissaoRepository{
				buscarPorIDFn: func(ctx context.Context, ctgID, usrID string) (*cpm.CategoriaPermissao, error) {
					return &cpm.CategoriaPermissao{}, nil
				},
				deletarFn: func(ctx context.Context, ctgID, usrID string) error {
					return errors.New("erro no banco de dados")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaPermissaoService(tt.repo)
			err := service.Deletar(ctx, "ctg-123", "usr-456")
			if (err != nil) != tt.wantErr {
				t.Errorf("Deletar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// TestCategoriaPermissaoService_Listar testa o método Listar do CategoriaPermissaoService.
func TestCategoriaPermissaoService_Listar(t *testing.T) {
	ctx := context.Background()

	cpm1, _ := cpm.Novo("ctg-1", "usr-1", usr.PermADM)
	cpm2, _ := cpm.Novo("ctg-2", "usr-2", usr.PermDEV)

	categoriasPermissoes := []cpm.CategoriaPermissao{*cpm1, *cpm2}

	tests := []struct {
		name    string
		repo    cpm.Repository
		filtro cpm.Filtro
		wantTotal int
		wantErr bool
	}{
		{
			name: "Listar todas as categorias de permissao e normalizar filtro",
			repo: &fakeCategoriaPermissaoRepository{
				listarFn: func(ctx context.Context, f cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
					if f.Pagina() <= 0 {
						t.Fatalf("esperava pagina normalizada, recebeu %d", f.Pagina())
					}
					if f.Limite() <= 0 {
						t.Fatalf("esperava limite normalizado, recebeu %d", f.Limite())
					}
					return categoriasPermissoes, len(categoriasPermissoes), nil
				},
			},
			wantTotal: len(categoriasPermissoes),
			wantErr: false,
		},
		{
			name: "Falha ao listar categorias de permissao devido a erro no banco de dados",
			repo: &fakeCategoriaPermissaoRepository{
				listarFn: func(ctx context.Context, f cpm.Filtro) ([]cpm.CategoriaPermissao, int, error) {
					return nil, 0, errors.New("erro no banco de dados")
				},
			},
			wantTotal: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoCategoriaPermissaoService(tt.repo)

			got, total, _, err := service.Listar(ctx, tt.filtro)
			if (err != nil) != tt.wantErr {
				t.Errorf("Listar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if total != tt.wantTotal {
				t.Errorf("Listar() total = %v, want %v", total, tt.wantTotal)
			}
			if !tt.wantErr && len(got) != tt.wantTotal {
				t.Errorf("Listar() got length = %v, want %v", len(got), tt.wantTotal)
			}
		})
	}
}
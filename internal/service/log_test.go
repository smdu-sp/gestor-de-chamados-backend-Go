package service

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
)

// --- Fake Repository -------------------------------------------------------------------------------------------------

// fakeLogRepository é uma implementação falsa de log.Repository para testes.
type fakeLogRepository struct {
	buscarPorIDFn func(ctx context.Context, id string) (*log.Log, error)
	criarFn       func(ctx context.Context, l log.Log) error
	listarFn      func(ctx context.Context, filtro log.LogFiltro) ([]log.Log, int, error)
}

// BuscarPorID chama a função fake definida.
func (f *fakeLogRepository) BuscarPorID(ctx context.Context, id string) (*log.Log, error) {
	return f.buscarPorIDFn(ctx, id)
}

// Criar chama a função fake definida.
func (f *fakeLogRepository) Criar(ctx context.Context, l log.Log) error {
	return f.criarFn(ctx, l)
}

func (f *fakeLogRepository) Listar(ctx context.Context, filtro log.LogFiltro) ([]log.Log, int, error) {
	return f.listarFn(ctx, filtro)
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestLogService_BuscarPorID testa o método BuscarPorID do LogService.
func TestLogService_BuscarPorID(t *testing.T) {
	ctx := context.Background()

	logValido := &log.Log{} // não importa o conteúdo aqui

	tests := []struct {
		name    string
		repo    log.Repository
		wantErr bool
	}{
		{
			name: "buscar log com sucesso",
			repo: &fakeLogRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*log.Log, error) {
					return logValido, nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao buscar log no repositorio",
			repo: &fakeLogRepository{
				buscarPorIDFn: func(ctx context.Context, id string) (*log.Log, error) {
					return nil, errors.New("log não encontrado")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NovoLogService(nil, tt.repo)

			result, err := service.BuscarPorID(ctx, "log-123")

			if tt.wantErr && err == nil {
				t.Fatalf("esperava erro, mas recebeu nil")
			}

			if !tt.wantErr && err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if !tt.wantErr && result == nil {
				t.Fatalf("esperava log, mas recebeu nil")
			}
		})
	}
}

// TestLogService_Criar testa o método Criar do LogService.
func TestLogService_Criar(t *testing.T) {
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
		name      string
		geradorID dmn.GeradorID
		repo      log.Repository
		wantErr   bool
	}{
		{
			name: "criar log com sucesso",
			geradorID: &fakeGeradorID{
				id: "log-123",
			},
			repo: &fakeLogRepository{
				criarFn: func(ctx context.Context, l log.Log) error {
					return nil
				},
			},
			wantErr: false,
		},
		{
			name: "erro ao criar log no repositorio",
			geradorID: &fakeGeradorID{
				id: "log-123",
			},
			repo: &fakeLogRepository{
				criarFn: func(ctx context.Context, l log.Log) error {
					return errors.New("erro ao criar log")
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Arrange
			service := NovoLogService(tt.geradorID, tt.repo)

			// Capturar apenas logs de erro (nível ERROR ou superior)
			var loggedError bool
			handler := slog.NewTextHandler(&testLogWriter{func(msg string) {
				loggedError = true
			}}, &slog.HandlerOptions{
				Level: slog.LevelError, // Captura apenas ERROR ou superior
			})
			slog.SetDefault(slog.New(handler))

			// Act
			service.Criar(ctx, log.Criar, "entidade", "detalhes")

			// Assert
			if tt.wantErr && !loggedError {
				t.Fatalf("esperava log de erro, mas nenhum foi registrado")
			}

			if !tt.wantErr && loggedError {
				t.Fatalf("não esperava log de erro, mas um foi registrado")
			}
		})
	}
}

// TestLogService_Listar testa o método Listar do LogService.
func TestLogService_Listar_Sucesso(t *testing.T) {
	ctx := context.Background()

	// Arrange
	repo := &fakeLogRepository{
		listarFn: func(ctx context.Context, f log.LogFiltro) ([]log.Log, int, error) {
			l, _ := log.Novo(
				"log-1",
				"user-123",
				log.Criar,
				"Subcategoria",
				"criou subcategoria",
			)
			return []log.Log{*l}, 1, nil
		},
	}

	service := NovoLogService(nil, repo)

	filtro := log.NovoLogFiltro(
		dmn.NovoPaginacao(0, 0), // valores não normalizados
		nil, nil, nil, nil, nil, nil,
	)

	// Act
	logs, total, filtroRetornado, err := service.Listar(ctx, filtro)

	// Assert
	if err != nil {
		t.Fatalf("erro inesperado: %v", err)
	}

	if total != 1 {
		t.Fatalf("esperava total 1, recebeu %d", total)
	}

	if len(logs) != 1 {
		t.Fatalf("esperava 1 log, recebeu %d", len(logs))
	}

	if filtroRetornado.Pagina() <= 0 {
		t.Fatalf("esperava filtro normalizado")
	}
}


// =====================================================================================================================
// FUNÇÕES AUXILIARES
// =====================================================================================================================

// testLogWriter é um writer personalizado para capturar logs durante os testes.
type testLogWriter struct {
	writeFn func(msg string)
}

// Write implementa a interface io.Writer.
func (w *testLogWriter) Write(p []byte) (n int, err error) {
	w.writeFn(string(p))
	return len(p), nil
}

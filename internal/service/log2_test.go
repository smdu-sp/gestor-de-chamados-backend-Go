package service

import (
	"context"
	"errors"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
)

// --- Fake Repository (simplificado com estado) ---------------------------------------------------

// fakeLogRepository2 é uma implementação falsa de log.Repository para testes.
type fakeLogRepository2 struct {
	logs         map[string]*log.Log
	erroAoCriar  error
	erroAoBuscar error
	erroAoListar error
	totalListar  int
}

// novoFakeLogRepository cria uma nova instância de fakeLogRepository2.
func novoFakeLogRepository2() *fakeLogRepository2 {
	return &fakeLogRepository2{
		logs: make(map[string]*log.Log),
	}
}

// Criar cria um log no repositório falso.
func (f *fakeLogRepository2) Criar(ctx context.Context, l log.Log) error {
	if f.erroAoCriar != nil {
		return f.erroAoCriar
	}
	f.logs[l.ID()] = &l
	return nil
}

// BuscarPorID busca um log por ID no repositório falso.
func (f *fakeLogRepository2) BuscarPorID(ctx context.Context, id string) (*log.Log, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	l, existe := f.logs[id]
	if !existe {
		return nil, errors.New("log não encontrado")
	}
	return l, nil
}

// Listar lista logs no repositório falso.
func (f *fakeLogRepository2) Listar(ctx context.Context, filtro log.LogFiltro) ([]log.Log, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	logs := make([]log.Log, 0, len(f.logs))
	for _, l := range f.logs {
		logs = append(logs, *l)
	}

	total := f.totalListar
	if total == 0 {
		total = len(logs)
	}

	return logs, total, nil
}

// =====================================================================================================================
// HELPERS
// =====================================================================================================================

func novoLogTeste(id, usuarioID, entidade, detalhes string, acao log.Acao) *log.Log {
	logTeste, _ := log.Novo(id, usuarioID, acao, entidade, detalhes)
	return logTeste
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// TestLogService_BuscarPorID testa o método BuscarPorID do LogService.
func TestLogService_BuscarPorID2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() log.Repository
		id        string
		wantErr   bool
		assertFn  func(t *testing.T, l *log.Log)
	}{
		{
			name: "buscar log com sucesso",
			repoSetup: func() log.Repository {
				repo := novoFakeLogRepository2()
				repo.logs["log-123"] = novoLogTeste("log-123", "user-456", "EntidadeX", "Detalhes do log", log.Criar)
				return repo
			},
			id:      "log-123",
			wantErr: false,
			assertFn: func(t *testing.T, l *log.Log) {
				t.Helper()
				if l.ID() != "log-123" {
					t.Errorf("esperado ID 'log-123', recebido '%s'", l.ID())
				}
				if l.Entidade() != "EntidadeX" {
					t.Errorf("esperado Entidade 'EntidadeX', recebido '%s'", l.Entidade())
				}
			},
		},
		{
			name: "erro ao buscar log",
			repoSetup: func() log.Repository {
				repo := novoFakeLogRepository2()
				repo.erroAoBuscar = errors.New("erro no banco")
				return repo
			},
			id:      "log-999",
			wantErr: true,
			assertFn: func(t *testing.T, l *log.Log) {
				t.Helper()

				if l != nil {
					t.Errorf("esperado log nil, recebido '%v'", l)
				}
			},
		},
		{
			name: "log não encontrado",
			repoSetup: func() log.Repository {
				return novoFakeLogRepository2()
			},
			id:      "log-inexistente",
			wantErr: true,
			assertFn: func(t *testing.T, l *log.Log) {
				t.Helper()

				if l != nil {
					t.Errorf("esperado log nil, recebido '%v'", l)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoLogService(nil, repo)

			l, err := service.BuscarPorID(ctx, tt.id)

			assertError(t, err, tt.wantErr)

			if tt.assertFn != nil {
				tt.assertFn(t, l)
			}
		})
	}
}

// TestLogService_Criar testa o método Criar do LogService.
func TestLogService_Criar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name 		string
		geradorID  dmn.GeradorID
		repoSetup func() log.Repository
		usuarioID 	string
		acao 		log.Acao
		entidade 	string
		detalhes 	string
		wantErr 	bool
		assertFn 	func(t *testing.T, repo log.Repository)
	}{
		{
			name: "criar log com sucesso",
			geradorID: &fakeGeradorID{
				id: "log-123",
			},
			repoSetup: func() log.Repository {
				return novoFakeLogRepository2()
			},
			usuarioID: "user-456",
			acao:      log.Criar,
			entidade:  "EntidadeX",
			detalhes:  "Detalhes do log",
			wantErr:   false,
			assertFn: func(t *testing.T, repo log.Repository) {
				t.Helper()
				fakeRepo := repo.(*fakeLogRepository2)
				l, exists := fakeRepo.logs["log-123"]
				if !exists {
					t.Fatalf("esperado log com ID 'log-123' não encontrado")
				}
				if l.Entidade() != "EntidadeX" {
					t.Errorf("esperado Entidade 'EntidadeX', recebido '%s'", l.Entidade())
				}
			},
		},
		{
			name: "erro ao criar log",
			geradorID: &fakeGeradorID{
				id: "log-123",
			},
			repoSetup: func() log.Repository {
				repo := novoFakeLogRepository2()
				repo.erroAoCriar = errors.New("erro no banco")
				return repo
			},
			usuarioID: "user-456",
			acao:      log.Criar,
			entidade:  "EntidadeX",
			detalhes:  "Detalhes do log",
			wantErr:   true,
			assertFn:  func(t *testing.T, repo log.Repository) {
				t.Helper()
				fakeRepo := repo.(*fakeLogRepository2)
				_, exists := fakeRepo.logs["log-123"]
				if exists {
					t.Fatalf("não esperado log com ID 'log-123' foi criado")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoLogService(tt.geradorID, repo)

			// Contexto com usuário autenticado
			claims := &auth.Claims{
				ID: tt.usuarioID,
			}
			ctxComUsuario := auth.ContextoComClaims(ctx, claims)

			service.Criar(ctxComUsuario, tt.acao, tt.entidade, tt.detalhes)

			if tt.assertFn != nil {
				tt.assertFn(t, repo)
			}
		})
	}
}

// TestLogService_Listar2 testa o método Listar do LogService.
func TestLogService_Listar2(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		repoSetup func() log.Repository
		filtro    log.LogFiltro
		wantTotal int
		wantErr   bool
		assertFn  func(t *testing.T, logs []log.Log, filtroRetornado log.LogFiltro)
	}{
		{
			name: "listar logs com sucesso e normalizar filtro",
			filtro: log.LogFiltro{},
			repoSetup: func() log.Repository {
				repo := novoFakeLogRepository2()
				repo.logs["log-1"] = novoLogTeste("log-1", "user-1", "EntidadeA", "Detalhes A", log.Criar)
				repo.logs["log-2"] = novoLogTeste("log-2", "user-2", "EntidadeB", "Detalhes B", log.Atualizar)
				return repo
			},
			wantTotal: 2,
			wantErr:   false,
			assertFn: func(t *testing.T, logs []log.Log, filtroRetornado log.LogFiltro) {
				t.Helper()
				if len(logs) != 2 {
					t.Errorf("esperado 2 logs, recebido %d", len(logs))
				}

				// Verificar normalização do filtro
				if filtroRetornado.Pagina() <= 0 {
					t.Errorf("esperado Página 1, recebido %d", filtroRetornado.Pagina())
				}
				if filtroRetornado.Limite() <= 0 {
					t.Errorf("esperado Limite 10, recebido %d", filtroRetornado.Limite())
				}
			},
		},
		{
			name: "erro ao listar logs",
			filtro: log.LogFiltro{},
			repoSetup: func() log.Repository {
				repo := novoFakeLogRepository2()
				repo.erroAoListar = errors.New("erro no banco")
				return repo
			},
			wantTotal: 0,
			wantErr:   true,
			assertFn: func(t *testing.T, logs []log.Log, filtroRetornado log.LogFiltro) {
				t.Helper()
				if len(logs) != 0 {
					t.Errorf("esperado 0 logs, recebido %d", len(logs))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := tt.repoSetup()
			service := NovoLogService(nil, repo)

			logs, total, filtroRetornado, err := service.Listar(ctx, tt.filtro)

			assertError(t, err, tt.wantErr)

			if !tt.wantErr {
				if total != tt.wantTotal {
					t.Errorf("esperado total %d, recebido %d", tt.wantTotal, total)
				}

				if tt.assertFn != nil {
					tt.assertFn(t, logs, filtroRetornado)
				}
			}
		})
	}
}
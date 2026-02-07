package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

// =====================================================================================================================
// REPOSITÓRIO FALSO
// =====================================================================================================================

// fakeLogRepository é uma implementação falsa de repositório de logs para testes.
type fakeLogRepository struct {
	logs         map[string]*log.Log
	erroAoCriar  error
	erroAoBuscar error
	erroAoListar error
}

// novoFakeLogRepository cria uma nova instância de fakeLogRepository.
func novoFakeLogRepository() *fakeLogRepository {
	return &fakeLogRepository{
		logs: make(map[string]*log.Log),
	}
}

// Criar cria um log no repositório falso.
func (f *fakeLogRepository) Criar(ctx context.Context, l log.Log) error {
	if f.erroAoCriar != nil {
		return f.erroAoCriar
	}
	log := l
	f.logs[l.ID()] = &log
	return nil
}

// BuscarPorID busca um log por ID no repositório falso.
func (f *fakeLogRepository) BuscarPorID(ctx context.Context, id string) (*log.Log, error) {
	if f.erroAoBuscar != nil {
		return nil, f.erroAoBuscar
	}
	l, existe := f.logs[id]
	if !existe {
		return nil, mysql.ErrLogNaoEncontrado
	}
	return l, nil
}

// Listar lista logs no repositório falso.
func (f *fakeLogRepository) Listar(ctx context.Context, filtro log.LogFiltro) ([]log.Log, int, error) {
	if f.erroAoListar != nil {
		return nil, 0, f.erroAoListar
	}

	// Converter o mapa para slice
	logs := make([]log.Log, 0, len(f.logs))
	for _, l := range f.logs {
		logs = append(logs, *l)
	}

	// Ordenação ORDER BY criado_em DESC
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].CriadoEm().After(logs[j].CriadoEm())
	})

	// Contar total
	total := len(logs)

	// Paginação
	if filtro.Limite() > 0 {
		inicio := filtro.Offset()
		fim := inicio + filtro.Limite()

		if inicio > total {
			return []log.Log{}, total, nil
		}

		if fim > total {
			fim = total
		}

		logs = logs[inicio:fim]
	}

	return logs, total, nil
}

// =====================================================================================================================
// FUNÇÕES AUXILIARES DE TESTE
// =====================================================================================================================

const logTesteID = "log-123"

// novoLogTeste cria um log de teste com valores fixos.
func novoLogTeste() *log.Log {
	l, _ := log.Novo(
		logTesteID,
		usuarioTesteID,
		log.Criar,
		"EntidadeTeste",
		"Detalhes do log de teste",
	)
	return l
}

// popularLogsTeste adiciona multiplos logs de teste ao repositório falso para testes.
func popularLogsTeste(repo *fakeLogRepository, quantidade int) {
	for i := 1; i <= quantidade; i++ {
		id := fmt.Sprintf("log-%02d", i)

		l, _ := log.Novo(
			id,
			fmt.Sprintf("usuario-%02d", i),
			log.Criar,
			fmt.Sprintf("Entidade-%02d", i),
			fmt.Sprintf("Detalhes do log %02d", i),
		)
		repo.logs[id] = l
	}
}

// =====================================================================================================================
// TESTES
// =====================================================================================================================

// Test_LogService_BuscarPorID testa o método BuscarPorID do LogService.
func Test_LogService_BuscarPorID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() log.Repository
		id                 string
		erroEsperado       error
		verificarResultado func(t *testing.T, repo log.Repository, log *log.Log)
	}{
		{
			nome: "buscar log por id com sucesso",
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				novoLog := novoLogTeste()
				repo.logs[novoLog.ID()] = novoLog
				return repo
			},
			id:           logTesteID,
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo log.Repository, log *log.Log) {
				t.Helper()

				verificarIDs(t, logTesteID, log.ID())
			},
		},
		{
			nome: "erro ao buscar log por id no repositório",
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				novoLog := novoLogTeste()
				repo.logs[novoLog.ID()] = novoLog
				repo.erroAoBuscar = errFakeRepo
				return repo
			},
			id:           logTesteID,
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo log.Repository, log *log.Log) {
				t.Helper()

				_, ok := repo.(*fakeLogRepository).logs[logTesteID]
				verificarOk(t, ok, "deveria existir log no repositório")
				verificarNulo(t, log, "não esperava log encontrado quando há erro no repositório")
			},
		},
		{
			nome: "log não encontrado por id",
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				novoLog := novoLogTeste()
				repo.logs[novoLog.ID()] = novoLog
				repo.erroAoBuscar = mysql.ErrLogNaoEncontrado
				return repo
			},
			id:           "log-inexistente",
			erroEsperado: mysql.ErrLogNaoEncontrado,
			verificarResultado: func(t *testing.T, repo log.Repository, log *log.Log) {
				t.Helper()

				_, ok := repo.(*fakeLogRepository).logs[logTesteID]
				verificarOk(t, ok, "deveria existir log no repositório")
				verificarNulo(t, log, "não esperava log encontrado quando log não existe")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoLogService(nil, repo)
			logEncontrado, erroRecebido := service.BuscarPorID(ctx, tt.id)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, logEncontrado)
			}
		})
	}
}

// Test_LogService_Criar testa o método Criar do LogService.
func Test_LogService_Criar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		geradorID          dmn.GeradorID
		prepararRepo       func() log.Repository
		usuarioID          string
		acao               log.Acao
		entidade           string
		detalhes           string
		verificarResultado func(t *testing.T, repo log.Repository)
	}{
		{
			nome: "criar log com sucesso",
			geradorID: &fakeGeradorID{
				id: logTesteID,
			},
			prepararRepo: func() log.Repository {
				return novoFakeLogRepository()
			},
			usuarioID:    usuarioTesteID,
			acao:         log.Criar,
			entidade:     "EntidadeX",
			detalhes:     "Detalhes do log",
			verificarResultado: func(t *testing.T, repo log.Repository) {
				t.Helper()

				logCriado, ok := repo.(*fakeLogRepository).logs[logTesteID]
				if logCriado.Acao() != log.Criar.String() {
					t.Errorf("Ação esperada %s, recebido %s", log.Criar, logCriado.Acao())
				}
				if logCriado.Entidade() != "EntidadeX" {
					t.Errorf("Entidade esperada 'EntidadeX', recebido '%s'", logCriado.Entidade())
				}
				if logCriado.Detalhes() != "Detalhes do log" {
					t.Errorf("Detalhes esperado 'Detalhes do log', recebidos '%s'", logCriado.Detalhes())
				}

				verificarOk(t, ok, "esperava log criado no repositório")
				verificarIDs(t, logTesteID, logCriado.ID())
				verificarIDs(t, usuarioTesteID, logCriado.UsuarioID())
			},
		},
		{
			nome: "erro ao salvar log no repositório",
			geradorID: &fakeGeradorID{
				id: logTesteID,
			},
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				repo.erroAoCriar = errFakeRepo
				return repo
			},
			usuarioID:    usuarioTesteID,
			acao:         log.Criar,
			entidade:     "EntidadeX",
			detalhes:     "Detalhes do log",
			verificarResultado: func(t *testing.T, repo log.Repository) {
				t.Helper()

				verificarNaoPersistido(t, len(repo.(*fakeLogRepository).logs))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoLogService(tt.geradorID, repo)
			claims := &auth.Claims{ID: tt.usuarioID}
			ctxComUsuario := auth.ContextoComClaims(ctx, claims)

			service.Criar(ctxComUsuario, tt.acao, tt.entidade, tt.detalhes)

			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo)
			}
		})
	}
}

// Test_LogService_Listar testa o método Listar do LogService.
func Test_LogService_Listar(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		nome               string
		prepararRepo       func() log.Repository
		filtro             log.LogFiltro
		erroEsperado       error
		verificarResultado func(t *testing.T, repo log.Repository, logs []log.Log, total int, filtroRetornado log.LogFiltro)
	}{
		{
			nome: "listar 10 logs paginados com sucesso - pagina 1 limite 5",
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				popularLogsTeste(repo, 10)
				return repo
			},
			filtro:       log.NovoLogFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil, nil, nil, nil),
			erroEsperado: nil,
			verificarResultado: func(t *testing.T, repo log.Repository, logs []log.Log, total int, filtroRetornado log.LogFiltro) {
				t.Helper()

				fakeRepo := repo.(*fakeLogRepository)
				verificarPaginacao(t, filtroRetornado.Paginacao(), 1, 5)
				verificarTotal(t, total, len(fakeRepo.logs))
				verificarLimite(t, logs, 5)
			},
		},
		{
			nome: "erro ao listar logs",
			prepararRepo: func() log.Repository {
				repo := novoFakeLogRepository()
				popularLogsTeste(repo, 10)
				repo.erroAoListar = errFakeRepo
				return repo
			},
			filtro:       log.NovoLogFiltro(dmn.NovoPaginacao(1, 5), nil, nil, nil, nil, nil, nil),
			erroEsperado: errFakeRepo,
			verificarResultado: func(t *testing.T, repo log.Repository, logs []log.Log, total int, filtroRetornado log.LogFiltro) {
				t.Helper()

				verificarPaginacao(t, filtroRetornado.Paginacao(), 0, 0)
				verificarTotal(t, total, 0)
				verificarLimite(t, logs, 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			t.Parallel()

			repo := tt.prepararRepo()
			service := NovoLogService(nil, repo)
			logs, total, filtroRetornado, erroRecebido := service.Listar(ctx, tt.filtro)

			verificarErro(t, erroRecebido, tt.erroEsperado)
			if tt.verificarResultado != nil {
				tt.verificarResultado(t, repo, logs, total, filtroRetornado)
			}
		})
	}
}

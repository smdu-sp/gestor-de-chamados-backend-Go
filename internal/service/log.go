package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
)

// LogService representa a camada de caso de uso para operações relacionadas a logs.
type LogService struct {
	id   dmn.GeradorID
	repo log.Repository
}

// NovoLogService cria uma nova instância de LogService.
func NovoLogService(id dmn.GeradorID, repo log.Repository) *LogService {
	return &LogService{id: id, repo: repo}
}

// Asserção de interface para garantir que LogService implementa LogService
var _ log.Service = (*LogService)(nil)

// BuscarPorID recebe um ID e retorna o log correspondente.
//
// Erros sentinela possíveis: ErrLogNaoEncontrado.
func (u *LogService) BuscarPorID(ctx context.Context, id string) (*log.Log, error) {
	log, err := u.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar log por ID: %w", err)
	}
	return log, nil
}

// Criar cria um novo log.
func (u *LogService) Criar(ctx context.Context, acao log.Acao, entidade, detalhes string) {
	// 1 - Extrair o usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		slog.Error("Erro ao extrair usuarioID do contexto", "error", err)
	}

	// 2 - Gerar ID
	id, err := u.id.NovoID()
	if err != nil {
		slog.Error("Erro ao gerar UUIDv7 para log", "error", err)
	}

	// 3 - Criar o novo log
	log, err := log.Novo(id, usrAutenticado.ID(), acao, entidade, detalhes)
	if err != nil {
		slog.Error("Erro ao criar novo log", "error", err)
	}

	// 4 - Salvar o log no repositório
	err = u.repo.Criar(ctx, *log)
	if err != nil {
		slog.Error("Erro ao salvar log", "error", err)
	}
	slog.Info("Log criado", "log", log.String())
}

// ListarLogs recebe um filtro e retorna uma lista de logs que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (u *LogService) Listar(ctx context.Context, f log.LogFiltro) ([]log.Log, int, log.LogFiltro, error) {
	f.Normalizar()
	logSlice, total, err := u.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, f, fmt.Errorf("listar logs: %w", err)
	}

	return logSlice, total, f, nil
}

package usecase

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/jwt"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth/middleware"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

// LogUsecase representa a camada de caso de uso para operações relacionadas a logs.
type LogUsecase struct {
	repository repository.LogRepository
}

// NewLogUsecase cria uma nova instância de LogUsecase.
func NewLogUsecase(repository repository.LogRepository) *LogUsecase {
	return &LogUsecase{repository: repository}
}

// BuscarLogPorID busca um log pelo seu ID.
func (u *LogUsecase) BuscarLogPorID(ctx context.Context, id string) (*model.Log, error) {
	if id == "" {
		return nil, utils.NewAppError(
			"[usecase.BuscarLogPorID]",
			utils.LevelInfo,
			"erro ao buscar log por id",
			model.ErrIDInvalido,
		)
	}

	log, err := u.repository.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[usecase.BuscarLogPorID]: %w", err)
	}
	return log, nil
}

// CriarLog cria um novo log.
func (u *LogUsecase) CriarLog(ctx context.Context, acao model.Acao, entidade, detalhes string) error {
	const metodo = "[usecase.CriarLog]: %w"

	usuarioID, err := ExtrairUsuarioIDDoContexto(ctx)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	id, err := utils.NewUUIDv7String()
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	log, err := model.NewLog(
		id,
		usuarioID,
		acao,
		entidade,
		detalhes,
	)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}

	err = u.repository.Salvar(ctx, log)
	if err != nil {
		return fmt.Errorf(metodo, err)
	}
	return nil
}

// ListarLogs lista logs com paginação e filtros opcionais.
func (u *LogUsecase) ListarLogs(ctx context.Context, filtro model.LogFiltro) ([]model.Log, int, model.LogFiltro, error) {
	if filtro.Pagina <= 0 {
		filtro.Pagina = 1
	}

	if filtro.Limite <= 0 || filtro.Limite > 100 {
		filtro.Limite = 10
	}

	logs, total, err := u.repository.Listar(ctx, filtro)
	if err != nil {
		return nil, 0, filtro, fmt.Errorf("[usecase.ListarLogs]: %w", err)
	}

	return logs, total, filtro, nil
}

// Metodos auxiliares

// ExtrairUsuarioIDDoContexto extrai o ID do usuário do contexto.
func ExtrairUsuarioIDDoContexto(ctx context.Context) (string, error) {
	claims, ok := ctx.Value(middleware.ChaveUsuario).(*jwt.Claims)
	if !ok || claims == nil {
		return "", utils.NewAppError(
			"[usecase.ExtrairUsuarioIDDoContexto]",
			utils.LevelInfo,
			"erro ao extrair usuário do contexto",
			middleware.ErrUsuarioNaoAutenticado,
		)
	}
	return claims.ID, nil
}

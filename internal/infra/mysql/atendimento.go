package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
)

// Erros sentinela do repositório de atendimentos.
var (
	ErrAtendimentoNaoEncontrado = errors.New("atendimento não encontrado no banco de dados")
	ErrTecnicoNaoEncontrado     = errors.New("técnico atribuído não encontrado no banco de dados")
)

// AtendimentoRepository é a implementação do repositório de atendimentos para MySQL.
type AtendimentoRepository struct {
	db *sql.DB
}

// NewAtendimentoRepository cria uma nova instância de AtendimentoRepository.
func NovoAtendimentoRepository(db *sql.DB) *AtendimentoRepository {
	return &AtendimentoRepository{db: db}
}

// Asserção de interface para garantir que AtendimentoRepository implementa AtendimentoRepository
var _ atd.Repository = (*AtendimentoRepository)(nil)

// BuscarPorID recebe o ID do atendimento e retorna o atendimento correspondente.
//
// Erros sentinela possíveis: ErrAtendimentoNaoEncontrado.
func (r *AtendimentoRepository) BuscarPorID(ctx context.Context, id string) (*atd.Atendimento, error) {
	query :=
		`SELECT id, atribuido_id, chamado_id, criado_em, atualizado_em
		 FROM atendimentos 
		 WHERE id = ?`

	atd, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if atd == nil {
		return nil, ErrAtendimentoNaoEncontrado
	}

	return atd, nil
}

// BuscarPorChamadoEAtribuidoID recebe o ID do chamado e do técnico atribuído, e retorna o atendimento correspondente.
//
// Erros sentinela possíveis: ErrAtendimentoNaoEncontrado.
func (r *AtendimentoRepository) BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, atribuidoID string) (*atd.Atendimento, error) {
	query :=
		`SELECT id, atribuido_id, chamado_id, criado_em, atualizado_em
		 FROM atendimentos 
		 WHERE chamado_id = ? AND atribuido_id = ?`

	atd, err := r.buscar(ctx, query, chamadoID, atribuidoID)
	if err != nil {
		return nil, fmt.Errorf("buscar por chamado e atribuído ID no repo: %w", err)
	}

	if atd == nil {
		return nil, ErrAtendimentoNaoEncontrado
	}

	return atd, nil
}

// Criar insere um novo atendimento no repositório.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrTecnicoNaoEncontrado.
func (r *AtendimentoRepository) Criar(ctx context.Context, a atd.Atendimento) (*atd.Atendimento, error) {
	existeChamado, err := ExisteChamadoPorID(ctx, r.db, a.ChamadoID())
	if err != nil {
		return nil, fmt.Errorf("criar no repo: %w", err)
	}
	if !existeChamado {
		return nil, ErrChamadoNaoEncontrado
	}

	existeTecnico, err := ExisteUsuarioPorID(ctx, r.db, a.AtribuidoID())
	if err != nil {
		return nil, fmt.Errorf("criar no repo: %w", err)
	}
	if !existeTecnico {
		return nil, ErrTecnicoNaoEncontrado
	}

	query := `
		INSERT INTO atendimentos (
	 	id, atribuido_id, chamado_id, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(ctx, query,
		a.ID(),
		a.AtribuidoID(),
		a.ChamadoID(),
		a.CriadoEm(),
		a.AtualizadoEm(),
	)
	if err != nil {
		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	atdCriado, err := r.BuscarPorID(ctx, a.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return atdCriado, nil
}

// Atualizar recebe o ID do atendimento e os dados atualizados, e atualiza o atendimento no repositório.
func (r *AtendimentoRepository) Atualizar(ctx context.Context, id string, a atd.Atendimento) (*atd.Atendimento, error) {
	query := `
	  UPDATE atendimentos
		SET atribuido_id = ?, chamado_id = ?, atualizado_em = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		a.AtribuidoID(),
		a.ChamadoID(),
		a.AtualizadoEm(),
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	atdAtualizado, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return atdAtualizado, nil
}

// Listar recebe um filtro e retorna uma lista de atendimentos que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (r *AtendimentoRepository) Listar(ctx context.Context, filtro atd.Filtro) ([]atd.Atendimento, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(filtro)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	atdSlice, err := scanRows(rows, r.scanAtendimento)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return atdSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e seus argumentos, executa a consulta e retorna um atendimento.
func (r *AtendimentoRepository) buscar(ctx context.Context, query string, args ...any) (*atd.Atendimento, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	atd, err := r.scanAtendimento(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return atd, nil
}

// ExisteAtendimentoPorID recebe o ID do atendimento e verifica se ele existe no banco de dados.
func ExisteAtendimentoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM atendimentos WHERE id = ?)"
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// scanAtendimento recebe um scanner (row ou rows) e retorna um atendimento.
func (r *AtendimentoRepository) scanAtendimento(scanner scanner) (*atd.Atendimento, error) {
	var atdDB atd.AtendimentoDB
	err := scanner.Scan(
		&atdDB.ID,
		&atdDB.AtribuidoID,
		&atdDB.ChamadoID,
		&atdDB.CriadoEm,
		&atdDB.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return atd.CarregarDoBD(atdDB), nil
}

// construirQueryListar recebe um filtro e constrói a query SQL correspondente, retornando a query e os argumentos.
func (r *AtendimentoRepository) construirQueryListar(filtro atd.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
		id, atribuido_id, chamado_id, criado_em, atualizado_em
		FROM atendimentos
		WHERE 1=1
	`)

	if filtro.ChamadoID() != nil {
		query.WriteString(" AND chamado_id = ?")
		args = append(args, *filtro.ChamadoID())
	}

	if filtro.AtribuidoID() != nil {
		query.WriteString(" AND atribuido_id = ?")
		args = append(args, *filtro.AtribuidoID())
	}

	// Ordenação e paginação
	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())

	return query.String(), args
}

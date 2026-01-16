package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
)

// Erros sentinela do repositório de acompanhamentos.
var ErrAcompanhamentoNaoEncontrado = errors.New("acompanhamento não encontrado no banco de dados")

// AcompanhamentoRepository é a implementação do repositório de acompanhamento para MySQL.
type AcompanhamentoRepository struct {
	db *sql.DB
}

// NewAcompanhamentoRepository cria uma nova instância do repositório de acompanhamento para MySQL.
func NovoAcompanhamentoRepository(db *sql.DB) *AcompanhamentoRepository {
	return &AcompanhamentoRepository{db: db}
}

// Asserção de interface para garantir que AcompanhamentoRepository implementa Repository
var _ acp.Repository = (*AcompanhamentoRepository)(nil)

// BuscarPorID recebe um ID de acompanhamento e retorna o acompanhamento correspondente.
//
// Erros sentinela possíveis: ErrAcompanhamentoNaoEncontrado.
func (r *AcompanhamentoRepository) BuscarPorID(ctx context.Context, id string) (*acp.Acompanhamento, error) {
	query :=
		`SELECT id, conteudo, chamado_id, usuario_id, 
		 remetente, criado_em, atualizado_em 
		 FROM acompanhamentos 
		 WHERE id = ?`

	acp, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if acp == nil {
		return nil, ErrAcompanhamentoNaoEncontrado
	}

	return acp, nil
}

// BuscarPorChamadoID recebe um ID de chamado e retorna uma lista de acompanhamentos correspondentes.
func (r *AcompanhamentoRepository) BuscarPorChamadoID(ctx context.Context, chamadoID string) ([]acp.Acompanhamento, error) {
	query := `
		SELECT 
			id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		FROM acompanhamentos
		WHERE chamado_id = ?
		ORDER BY criado_em ASC
	`
	rows, err := r.db.QueryContext(ctx, query, chamadoID)
	if err != nil {
		return nil, fmt.Errorf("buscar por chamado ID no repo: %w", err)
	}
	defer rows.Close()

	var acpSlice []acp.Acompanhamento
	for rows.Next() {
		acompanhamento, err := r.scanAcompanhamento(rows)
		if err != nil {
			return nil, fmt.Errorf("buscar por chamado ID no repo: %w", err)
		}
		acpSlice = append(acpSlice, *acompanhamento)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("buscar por chamado ID no repo: %w", err)
	}

	return acpSlice, nil
}

// Criar recebe um acompanhamento e o persiste no banco de dados.
func (r *AcompanhamentoRepository) Criar(ctx context.Context, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	query := `
		INSERT INTO acompanhamentos (
			id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		a.ID(),
		a.Conteudo(),
		a.ChamadoID(),
		a.UsuarioID(),
		a.Remetente(),
		a.CriadoEm(),
		a.AtualizadoEm(),
	)
	if err != nil {
		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	acpCriado, err := r.BuscarPorID(ctx, a.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return acpCriado, nil
}

// Atualizar recebe um ID de acompanhamento e os novos dados para atualizar o acompanhamento correspondente no banco de dados.
func (r *AcompanhamentoRepository) Atualizar(ctx context.Context, id string, a acp.Acompanhamento) (*acp.Acompanhamento, error) {
	query := `
		UPDATE acompanhamentos
		SET conteudo = ?, chamado_id = ?, usuario_id = ?, remetente = ?, atualizado_em = ?
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		a.Conteudo(),
		a.ChamadoID(),
		a.UsuarioID(),
		a.Remetente(),
		a.AtualizadoEm(),
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	acpAtualizado, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return acpAtualizado, nil
}

// Deletar recebe um ID de acompanhamento e remove o acompanhamento correspondente do banco de dados.
//
// Erros sentinela possíveis: ErrAcompanhamentoNaoEncontrado.
func (r *AcompanhamentoRepository) Deletar(ctx context.Context, id string) error {
	existe, err := ExisteAcompanhamentoPorID(ctx, r.db, id)
	if err != nil {
		return fmt.Errorf("deletar no repo: %w", err)
	}
	if !existe {
		return ErrAcompanhamentoNaoEncontrado
	}

	query := `DELETE FROM acompanhamentos WHERE id = ?`

	_, err = r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("deletar no repo: %w", err)
	}

	return nil
}

// Listar recebe filtros e retorna uma lista de acompanhamentos que correspondem aos critérios.
func (r *AcompanhamentoRepository) Listar(ctx context.Context, filtro acp.Filtro) ([]acp.Acompanhamento, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(filtro)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	acpSlice, err := scanRows(rows, r.scanAcompanhamento)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return acpSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e argumentos, executa a query e retorna um acompanhamento.
func (r *AcompanhamentoRepository) buscar(ctx context.Context, query string, args ...any) (*acp.Acompanhamento, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	acp, err := r.scanAcompanhamento(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return acp, nil
}

// ExisteAcompanhamentoPorID recebe um ID de acompanhamento e verifica se ele existe no banco de dados.
func ExisteAcompanhamentoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM acompanhamentos WHERE id = ?)"
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// scanAcompanhamento recebe um scanner (row ou rows) e retorna um acompanhamento escaneado.
func (r *AcompanhamentoRepository) scanAcompanhamento(scanner scanner) (*acp.Acompanhamento, error) {
	var acpDB acp.AcompanhamentoDB
	err := scanner.Scan(
		&acpDB.ID,
		&acpDB.Conteudo,
		&acpDB.ChamadoID,
		&acpDB.UsuarioID,
		&acpDB.Remetente,
		&acpDB.CriadoEm,
		&acpDB.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}
	return acp.CarregarDoDB(acpDB), nil
}

// construirQueryListar recebe filtros e constrói a query SQL correspondente.
func (r *AcompanhamentoRepository) construirQueryListar(filtro acp.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(`
		SELECT SQL_CALC_FOUND_ROWS
			id, conteudo, chamado_id, usuario_id, remetente, criado_em, atualizado_em
		FROM acompanhamentos
		WHERE 1=1
	`)

	if filtro.ChamadoID() != nil {
		query.WriteString(" AND chamado_id = ?")
		args = append(args, *filtro.ChamadoID())
	}

	if filtro.UsuarioID() != nil {
		query.WriteString(" AND usuario_id = ?")
		args = append(args, *filtro.UsuarioID())
	}

	query.WriteString(" ORDER BY criado_em ASC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
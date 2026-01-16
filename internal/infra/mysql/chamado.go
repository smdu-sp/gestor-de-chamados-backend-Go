package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
)

var ErrChamadoNaoEncontrado = errors.New("chamado não encontrado no banco de dados")

// ChamadoRepository implementa a interface ChamadoRepository para MySQL.
type ChamadoRepository struct {
	db *sql.DB
}

// NovoChamadoRepository cria uma nova instância de ChamadoRepository.
func NovoChamadoRepository(db *sql.DB) *ChamadoRepository {
	return &ChamadoRepository{db: db}
}

// Asserção de interface para garantir que ChamadoRepository implementa ChamadoRepository
var _ chm.Repository = (*ChamadoRepository)(nil)

// BuscarPorID recebe um ID e retorna o chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado.
func (r *ChamadoRepository) BuscarPorID(ctx context.Context, id string) (*chm.Chamado, error) {
	query :=
		`SELECT id, titulo, descricao, status, arquivado, criado_em, 
		 atualizado_em, solucionado_em, solucao, fechado_em, 
		 categoria_id, subcategoria_id, criador_id
		 FROM chamados 
		 WHERE id=?`

	chm, err := r.buscar(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("buscar por ID no repo: %w", err)
	}

	if chm == nil {
		return nil, ErrChamadoNaoEncontrado
	}

	return chm, nil
}

// Criar recebe um chamado e o insere no repositório.
func (r *ChamadoRepository) Criar(ctx context.Context, c chm.Chamado) (*chm.Chamado, error) {
	query := 
		`INSERT INTO chamados (
		 id, titulo, descricao, status, arquivado, categoria_id, 
		 subcategoria_id, criador_id, criado_em, atualizado_em
		 ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.ExecContext(ctx,	query, 
		c.ID(), 
		c.Titulo(), 
		c.Descricao(), 
		c.Status(), 
		c.Arquivado(), 
		c.CategoriaID(), 
		c.SubcategoriaID(), 
		c.CriadorID(),
		c.CriadoEm(),
		c.AtualizadoEm(),
	)
	if err != nil {
		return nil, fmt.Errorf("criar no repo: %w", err)
	}

	chmCriado, err := r.BuscarPorID(ctx, c.ID())
	if err != nil {
		return nil, fmt.Errorf("buscar após criar no repo: %w", err)
	}

	return chmCriado, nil
}

// Atualizar recebe um ID e um chamado atualizado, e aplica as mudanças no repositório.
func (r *ChamadoRepository) Atualizar(ctx context.Context, id string, c chm.Chamado) (*chm.Chamado, error) {
	query :=
		`UPDATE chamados 
		 SET titulo=?, descricao=?, status=?, arquivado=?, categoria_id=?, 
		 subcategoria_id=?, atualizado_em=?
		 WHERE id=?`

	_, err := r.db.ExecContext(ctx, query, 
		c.Titulo(), 
		c.Descricao(), 
		c.Status(), 
		c.Arquivado(), 
		c.CategoriaID(), 
		c.SubcategoriaID(), 
		c.AtualizadoEm(), 
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("atualizar no repo: %w", err)
	}

	chmAtualizado, err := r.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar após atualizar no repo: %w", err)
	}

	return chmAtualizado, nil
}

// Listar recebe um filtro e retorna uma lista de chamados que correspondem aos critérios do filtro,
// juntamente com o total de registros encontrados.
func (r *ChamadoRepository) Listar(ctx context.Context, filtro chm.Filtro) ([]chm.Chamado, int, error) {
	// Constrói a query base
	query, args := r.construirQueryListar(filtro)

	// Executa a query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}
	defer rows.Close()

	// Scaneia os resultados
	chmSlice, err := scanRows(rows, r.scanChamado)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	// Obtém o total de registros
	total, err := obterTotalRegistros(ctx, r.db)
	if err != nil {
		return nil, 0, fmt.Errorf("listar no repo: %w", err)
	}

	return chmSlice, total, nil
}

// =====================================================================================================================
// Métodos auxiliares
// =====================================================================================================================

// buscar recebe uma query e argumentos, executa a busca e retorna o chamado correspondente.
func (r *ChamadoRepository) buscar(ctx context.Context, query string, args ...any) (*chm.Chamado, error) {
	row := r.db.QueryRowContext(ctx, query, args...)
	chm, err := r.scanChamado(row)
	if err != nil {
		return nil, fmt.Errorf("buscar: %w", err)
	}
	return chm, nil
}

// ExisteChamadoPorID recebe um ID e verifica se um chamado com esse ID existe no repositório.
func ExisteChamadoPorID(ctx context.Context, db *sql.DB, id string) (bool, error) {
	var existe bool
	query := `SELECT EXISTS(SELECT 1 FROM chamados WHERE id=?)`
	err := db.QueryRowContext(ctx, query, id).Scan(&existe)
	if err != nil {
		return false, fmt.Errorf("verificar existência: %w", err)
	}

	return existe, nil
}

// scanChamado recebe um scanner (row ou rows) e retorna um chamado.
func (r *ChamadoRepository) scanChamado(scanner scanner) (*chm.Chamado, error) {
	var chmDB chm.ChamadoDB
	err := scanner.Scan(
		&chmDB.ID,
		&chmDB.Titulo,
		&chmDB.Descricao,
		&chmDB.Status,
		&chmDB.CriadoEm,
		&chmDB.AtualizadoEm,
		&chmDB.SolucionadoEm,
		&chmDB.Solucao,
		&chmDB.FechadoEm,
		&chmDB.CategoriaID,
		&chmDB.SubcategoriaID,
		&chmDB.CriadorID,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("scan: %w", err)
	}

	return chm.CarregarDoDB(chmDB), nil
}

// construirQueryListar recebe um filtro e constrói a query SQL correspondente.
func (r *ChamadoRepository) construirQueryListar(filtro chm.Filtro) (string, []any) {
	var query strings.Builder
	args := []any{}

	query.WriteString(
		`SELECT SQL_CALC_FOUND_ROWS
		id, titulo, descricao, status, criado_em, 
		atualizado_em, solucionado_em, solucao, fechado_em, 
		categoria_id, subcategoria_id, criador_id
		FROM chamados WHERE 1=1`,
	)

	if filtro.Busca() != nil && *filtro.Busca() != "" {
		padrao := "%" + *filtro.Busca() + "%"
		query.WriteString(" AND (titulo LIKE ? OR descricao LIKE ?)")
		args = append(args, padrao, padrao)
	}

	if filtro.Status() != nil && *filtro.Status() != "" {
		query.WriteString(" AND status = ?")
		args = append(args, *filtro.Status())
	}

	if filtro.CategoriaID() != nil && *filtro.CategoriaID() != "" {
		query.WriteString(" AND categoria_id = ?")
		args = append(args, *filtro.CategoriaID())
	}

	if filtro.SubcategoriaID() != nil && *filtro.SubcategoriaID() != "" {
		query.WriteString(" AND subcategoria_id = ?")
		args = append(args, *filtro.SubcategoriaID())
	}

	if filtro.CriadorID() != nil && *filtro.CriadorID() != "" {
		query.WriteString(" AND criador_id = ?")
		args = append(args, *filtro.CriadorID())
	}

	query.WriteString(" ORDER BY criado_em DESC LIMIT ? OFFSET ?")
	args = append(args, filtro.Limite(), filtro.Offset())
	return query.String(), args
}
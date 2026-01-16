package service

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

var (
	ErrCategoriaInativa                  = fmt.Errorf("categoria está inativa")
	ErrSubcategoriaInativa               = fmt.Errorf("subcategoria está inativa")
	ErrSubcategoriaNaoPertenceCategoria  = fmt.Errorf("subcategoria não pertence à categoria")
	ErrUsuarioNaoCriadorChamado          = fmt.Errorf("usuário não é o criador do chamado")
	ErrUsuarioNaoTecnicoAtribuidoChamado = fmt.Errorf("usuário não é técnico atribuído ao chamado")
)

// ChamadoService representa a camada de caso de uso para operações relacionadas a chamados.
type ChamadoService struct {
	id                     dmn.GeradorID
	repo                   chm.Repository
	categoriaRepo          ctg.Repository
	subcategoriaRepo       subc.Repository
	categoriaPermissaoRepo cpm.Repository
	atendimentoRepo        atd.Repository
	usuarioRepo            usr.Repository
}

// NovoChamadoService cria uma nova instância de ChamadoService.
func NovoChamadoService(
	id dmn.GeradorID,
	repo chm.Repository,
	categoriaRepo ctg.Repository,
	subcategoriaRepo subc.Repository,
	categoriaPermissaoRepo cpm.Repository,
	atendimentoRepo atd.Repository,
	usuarioRepo usr.Repository,
) *ChamadoService {
	return &ChamadoService{
		id:                     id,
		repo:                   repo,
		categoriaRepo:          categoriaRepo,
		subcategoriaRepo:       subcategoriaRepo,
		categoriaPermissaoRepo: categoriaPermissaoRepo,
		atendimentoRepo:        atendimentoRepo,
		usuarioRepo:            usuarioRepo,
	}
}

// Asserção de interface para garantir que ChamadoService implementa ChamadoService
var _ chm.Service = (*ChamadoService)(nil)

// BuscarPorID recebe um ID e retorna o chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado.
func (c *ChamadoService) BuscarPorID(ctx context.Context, id string) (*chm.Chamado, error) {
	chm, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar chamado por ID: %w", err)
	}
	return chm, nil
}

// Criar recebe os parâmetros de criação e cria um novo chamado.
//
// Erros sentinela possíveis: ErrCategoriaInativa, ErrSubcategoriaInativa, 
// ErrSubcategoriaNaoPertenceCategoria.
func (c *ChamadoService) Criar(ctx context.Context, criar chm.CriarParams) (*chm.Chamado, error) {
	// 1 - validar os parâmetros de criação
	err := c.validarCategoriaESubcategoria(ctx, criar.CategoriaID, criar.SubcategoriaID)
	if err != nil {
		return nil, fmt.Errorf("criar chamado: %w", err)
	}

	// 2 - gerar um novo ID para o chamado
	id, err := c.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar chamado: %w", err)
	}

	// 3 - criar a entidade chamado
	chamado, err := chm.Novo(
		id,
		criar.CategoriaID,
		criar.SubcategoriaID,
		criar.CriadorID,
		criar.Titulo,
		criar.Descricao,
	)
	if err != nil {
		return nil, err
	}

	// 4 - persistir o chamado
	chamadoCriado, err := c.repo.Criar(ctx, *chamado)
	if err != nil {
		return nil, fmt.Errorf("criar chamado: %w", err)
	}

	return chamadoCriado, nil
}

// Atualizar recebe um ID e parâmetros de atualização, e aplica as mudanças no chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrCategoriaInativa, ErrosValidacao,
// ErrSubcategoriaInativa, ErrSubcategoriaNaoPertenceCategoria, ErrUsuarioNaoCriadorChamado.
func (c *ChamadoService) Atualizar(ctx context.Context, id string, atualizar chm.AtualizarParams) (*chm.Chamado, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("atualizar chamado: %w", err)
	}

	// 2 - Buscar o chamado existente
	chamado, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar chamado: %w", err)
	}

	// 3 - Validar se o usuário é o criador do chamado
	if err := validarCriador(chamado, usrAutenticado); err != nil {
		return nil, fmt.Errorf("atualizar chamado: %w", err)
	}

	// 4 - Validar os parâmetros de atualização
	if atualizar.CategoriaID != nil && atualizar.SubcategoriaID != nil {
		err = c.validarCategoriaESubcategoria(ctx, *atualizar.CategoriaID, *atualizar.SubcategoriaID)
		if err != nil {
			return nil, fmt.Errorf("atualizar chamado: %w", err)
		}
	}

	// 5 - Atualizar os dados do chamado
	if err := chamado.AtualizarDados(atualizar); err != nil {
		return nil, err
	}

	// 6 - Persistir o chamado atualizado
	chamadoAtualizado, err := c.repo.Atualizar(ctx, id, *chamado)
	if err != nil {
		return nil, fmt.Errorf("atualizar chamado: %w", err)
	}

	return chamadoAtualizado, nil
}

// Arquivar recebe um ID e arquiva o chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrUsuarioNaoCriadorChamado.
func (c *ChamadoService) Arquivar(ctx context.Context, id string) (*chm.Chamado, error) {
	// 1 - obter o ID do usuário autenticado a partir do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("arquivar chamado: %w", err)
	}

	// 2 - buscar o chamado existente
	chamado, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("arquivar chamado: %w", err)
	}

	// 3 - validar se o usuário é o criador do chamado
	if err := validarCriador(chamado, usrAutenticado); err != nil {
		return nil, fmt.Errorf("arquivar chamado: %w", err)
	}

	// 4 - marcar o chamado como arquivado
	chamado.Arquivar()

	// 5 - persistir o chamado atualizado
	chamadoAtualizado, err := c.repo.Atualizar(ctx, id, *chamado)
	if err != nil {
		return nil, fmt.Errorf("arquivar chamado: %w", err)
	}

	return chamadoAtualizado, nil
}

// Desarquivar recebe um ID e desarquiva o chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrUsuarioNaoCriadorChamado.
func (c *ChamadoService) Desarquivar(ctx context.Context, id string) (*chm.Chamado, error) {
	// 1 - obter o ID do usuário autenticado a partir do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("desarquivar chamado: %w", err)
	}

	// 2 - buscar o chamado existente
	chamado, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("desarquivar chamado: %w", err)
	}

	// 3 - validar se o usuário é o criador do chamado
	if err := validarCriador(chamado, usrAutenticado); err != nil {
		return nil, fmt.Errorf("desarquivar chamado: %w", err)
	}

	// 4 - marcar o chamado como desarquivado
	chamado.Desarquivar()

	// 5 - persistir o chamado atualizado
	chamadoAtualizado, err := c.repo.Atualizar(ctx, id, *chamado)
	if err != nil {
		return nil, fmt.Errorf("desarquivar chamado: %w", err)
	}

	return chamadoAtualizado, nil
}

// AtualizarStatus recebe um ID e parâmetros de atualização de status, e aplica as mudanças no chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrUsuarioNaoTecnicoAtribuidoChamado.
func (c *ChamadoService) AtualizarStatus(ctx context.Context, id string, params chm.AtualizarStatusParams) (*chm.Chamado, error) {
	// 1 - obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("atualizar status do chamado: %w", err)
	}

	// 2 - buscar o chamado existente
	chamado, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar status do chamado: %w", err)
	}

	// 3 - validar autorização do usuário técnico ou administrador
	if err := c.validarAutorizacaoTecnica(ctx, chamado, usrAutenticado); err != nil {
		return nil, fmt.Errorf("atualizar status do chamado: %w", err)
	}

	// 4 - atualizar o status do chamado
	if err := chamado.AtualizarStatus(params.Status, params.Solucao); err != nil {
		return nil, err
	}

	// 5 - persistir o chamado atualizado
	chamadoAtualizado, err := c.repo.Atualizar(ctx, id, *chamado)
	if err != nil {
		return nil, fmt.Errorf("atualizar status do chamado: %w", err)
	}

	return chamadoAtualizado, nil
}

// AtualizarSolucao recebe um ID e parâmetros de atualização de solução, 
// e aplica as mudanças no chamado correspondente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrUsuarioNaoTecnicoAtribuidoChamado.
func (c *ChamadoService) AtualizarSolucao(ctx context.Context, id string, a chm.AtualizarSolucaoParams) (*chm.Chamado, error) {
	// 1 - obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("atualizar status do chamado: %w", err)
	}

	// 2 - buscar o chamado existente
	chamado, err := c.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar solução do chamado: %w", err)
	}

	// 3 - validar autorização do usuário técnico ou administrador
	if err := c.validarAutorizacaoTecnica(ctx, chamado, usrAutenticado); err != nil {
		return nil, fmt.Errorf("atualizar solução do chamado: %w", err)
	}

	// 4 - atualizar a solução do chamado
	if err := chamado.AtualizarSolucao(a.Solucao); err != nil {
		return nil, err
	}

	// 5 - persistir o chamado atualizado
	chamadoAtualizado, err := c.repo.Atualizar(ctx, id, *chamado)
	if err != nil {
		return nil, fmt.Errorf("atualizar solução do chamado: %w", err)
	}

	return chamadoAtualizado, nil
}

// Listar recebe um filtro e retorna a lista de chamados correspondentes, 
// o total de registros e o filtro aplicado.
func (c *ChamadoService) Listar(ctx context.Context, f chm.Filtro) ([]chm.Chamado, int, chm.Filtro, error) {
	f.Normalizar()
	chmSlice, total, err := c.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, f, fmt.Errorf("listar chamados: %w", err)
	}

	return chmSlice, total, f, nil
}

// =====================================================================================================================
// FUNÇÕES DE VALIDAÇÃO
// =====================================================================================================================

// validarCategoriaESubcategoria recebe IDs de categoria e subcategoria, 
// e valida se são válidos.
//
// Erros sentinela possíveis: ErrCategoriaNaoEncontrada, ErrCategoriaInativa, ErrSubcategoriaInativa, ErrSubcategoriaNaoPertenceCategoria.
func (c *ChamadoService) validarCategoriaESubcategoria(ctx context.Context, categoriaID, subcategoriaID string) error {

	// 1 - verifica se a categoria existe e esta ativa
	ctg, err := c.categoriaRepo.BuscarPorID(ctx, categoriaID)
	if err != nil {
		return err
	}
	if !ctg.Status() {
		return ErrCategoriaInativa
	}

	// 2 - verifica se a subcategoria existe e esta ativa
	subctg, err := c.subcategoriaRepo.BuscarPorID(ctx, subcategoriaID)
	if err != nil {
		return err
	}
	if !subctg.Status() {
		return ErrSubcategoriaInativa
	}

	// 3 - verifica se a subcategoria pertence à categoria
	if subctg.CategoriaID() != ctg.ID() {
		return ErrSubcategoriaNaoPertenceCategoria
	}

	return nil
}

// validarPermissaoCategoria recebe um ID de categoria e um ID de usuário, 
// e valida se o usuário tem permissão para acessar a categoria.
//
// Erros sentinela possíveis: ErrPermissaoAcessoCategoriaNegada.
func (c *ChamadoService) validarPermissaoCategoria(ctx context.Context, categoriaID, usuarioID string) error {
	_, err := c.categoriaPermissaoRepo.BuscarPorID(ctx, categoriaID, usuarioID)
	if err != nil {
		if err == mysql.ErrCategoriaPermissaoNaoEncontrada {
			return ErrPermissaoAcessoCategoriaNegada
		}
		return fmt.Errorf("validar permissão de categoria: %w", err)
	}

	return nil
}

// validarCriador recebe um chamado e um usuário autenticado, e valida se o usuário 
// é o criador do chamado.
//
// Erros sentinela possíveis: ErrUsuarioNaoCriadorChamado.
func validarCriador(chamado *chm.Chamado, usuario *auth.UsuarioAutenticado) error {
	if chamado.CriadorID() != usuario.ID() {
		return ErrUsuarioNaoCriadorChamado
	}
	return nil
}

// validarAutorizacaoTecnica recebe um chamado e um usuário autenticado, 
// e valida se o usuário é um técnico atribuído
// ou um administrador com permissão para a categoria do chamado.
//
// Erros sentinela possíveis: ErrUsuarioNaoTecnicoAtribuidoChamado.
func (c *ChamadoService) validarAutorizacaoTecnica(
	ctx context.Context,
	chamado *chm.Chamado,
	usuario *auth.UsuarioAutenticado,
) error {

	if usuario.EhTEC() {
		if _, err := c.atendimentoRepo.BuscarPorChamadoEAtribuidoID(ctx, chamado.ID(), usuario.ID()); err != nil {
			return ErrUsuarioNaoTecnicoAtribuidoChamado
		}
	}

	if usuario.EhADM() {
		if err := c.validarPermissaoCategoria(ctx, chamado.CategoriaID(), usuario.ID()); err != nil {
			return err
		}
	}

	return nil
}

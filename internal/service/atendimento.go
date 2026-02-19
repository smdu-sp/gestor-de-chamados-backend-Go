package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

var (
	ErrTecnicoAtribuidoNaoEncontrado        = errors.New("técnico atribuído não encontrado")
	ErrTecnicoAtribuidoDesativado           = errors.New("técnico atribuído está desativado")
	ErrChamadoStatusInvalidoParaAtendimento = errors.New("chamado está em estado inválido para atendimento (arquivado, fechado ou rejeitado)")
	ErrTecnicoNaoPodeAtribuirParaOutro      = errors.New("técnico só pode atribuir atendimento para si mesmo")
	ErrUsuarioNaoTemPermissaoTecnico        = errors.New("usuário não tem permissão de técnico")
	ErrAtendimentoAtivoJaExiste             = errors.New("já existe um atendimento ativo para este chamado")
	ErrPermissaoAcessoCategoriaNegada       = errors.New("permissão de acesso à categoria negada")
)

// AtendimentoService representa a camada de serviço para operações relacionadas a atendimentos.
type AtendimentoService struct {
	id                     dmn.GeradorID
	repo                   atd.Repository
	chamadoRepo            chm.Repository
	categoriaPermissaoRepo cpm.Repository
	usuarioRepo            usr.Repository
}

// NovoAtendimentoService cria uma nova instância de AtendimentoService.
func NovoAtendimentoService(
	id dmn.GeradorID,
	repo atd.Repository,
	chamadoRepo chm.Repository,
	categoriaPermissao cpm.Repository,
	usuarioRepo usr.Repository,
) *AtendimentoService {
	return &AtendimentoService{
		id:                     id,
		repo:                   repo,
		chamadoRepo:            chamadoRepo,
		categoriaPermissaoRepo: categoriaPermissao,
		usuarioRepo:            usuarioRepo,
	}
}

// Asserção de interface para garantir que AtendimentoService implementa AtendimentoService
var _ atd.Service = (*AtendimentoService)(nil)

// BuscarPorID recebe o ID do atendimento e retorna o atendimento correspondente.
//
// Erros sentinela possíveis: ErrAtendimentoNaoEncontrado.
func (a *AtendimentoService) BuscarPorID(ctx context.Context, id string) (*atd.Atendimento, error) {
	atd, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar atendimento por ID: %w", err)
	}
	return atd, nil
}

// BuscarPorChamadoEAtribuidoID recebe o ID do chamado e do técnico atribuído, e retorna o atendimento correspondente.
//
// Erros sentinela possíveis: ErrAtendimentoNaoEncontrado.
func (a *AtendimentoService) BuscarPorChamadoEAtribuidoID(ctx context.Context, chamadoID, tecAtribuidoID string) (*atd.Atendimento, error) {
	atd, err := a.repo.BuscarPorChamadoEAtribuidoID(ctx, chamadoID, tecAtribuidoID)
	if err != nil {
		return nil, fmt.Errorf("buscar atendimento por chamado e atribuido ID: %w", err)
	}
	return atd, nil
}

// Criar recebe os parâmetros para criar um novo atendimento e retorna o atendimento criado.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrTecnicoNaoEncontrado, ErrTecnicoAtribuidoDesativado,
// ErrUsuarioNaoTemPermissaoTecnico, ErrTecnicoNaoPodeAtribuirParaOutro, ErrPermissaoAcessoCategoriaNegada,
// ErrAtendimentoAtivoJaExiste.
func (a *AtendimentoService) Criar(ctx context.Context, criar atd.CriarParams) (*atd.Atendimento, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("criar atendimento: obter claims: %w", err)
	}

	// 2 - Buscar e validar chamado
	chamadoAtual, err := a.validarChamado(ctx, criar.ChamadoID)
	if err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 3 - Buscar e validar técnico atribuído
	tecnico, err := a.validarTecnico(ctx, criar.AtribuidoID, chamadoAtual.CategoriaID())
	if err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 4 - Validar permissões baseadas no tipo de usuário
	if err := a.validarPermissoesAtribuicao(ctx, usrAutenticado, tecnico, chamadoAtual); err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 5 - Verificar se já existe atendimento ativo
	if err := a.validarAtendimentoDuplicado(ctx, criar.ChamadoID, criar.AtribuidoID); err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 6 - Gerar ID
	id, err := a.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 7 - Criar entidade
	atendimento, err := atd.Novo(id, criar.AtribuidoID, criar.ChamadoID)
	if err != nil {
		return nil, err
	}

	// 8 - Persistir
	atendimentoCriado, err := a.repo.Criar(ctx, *atendimento)
	if err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	// 9 - atualizar status do chamado para "atribuido"
	chamadoStatusAtribuido := chamadoAtual.StatusAtribuido()

	// 10 - Persistir alteração do chamado
	if _, err := a.chamadoRepo.Atualizar(ctx, chamadoAtual.ID(), chamadoStatusAtribuido); err != nil {
		return nil, fmt.Errorf("criar atendimento: %w", err)
	}

	return atendimentoCriado, nil
}

// Atualizar atualiza um atendimento existente.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrTecnicoNaoEncontrado, ErrTecnicoAtribuidoDesativado,
// ErrUsuarioNaoTemPermissaoTecnico, ErrTecnicoNaoPodeAtribuirParaOutro, ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) Atualizar(ctx context.Context, id string, atualizar atd.AtualizarParams) (*atd.Atendimento, error) {
	// 1 - Obter claims do usuário autenticado que está atualizando o atendimento
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	// 2 - Buscar e validar chamado
	chamado, err := a.validarChamado(ctx, atualizar.ChamadoID)
	if err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	// 3 - Buscar e validar técnico atribuído
	tecnico, err := a.validarTecnico(ctx, atualizar.AtribuidoID, chamado.CategoriaID())
	if err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	// 4 - Validar permissões baseadas no tipo de usuário
	if err := a.validarPermissoesAtribuicao(ctx, usrAutenticado, tecnico, chamado); err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	// 5 - Buscar atendimento existente
	atendimentoAtual, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	// 6 - Atualizar dados do atendimento
	atendimentoAtualizado, err := atendimentoAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	// 7 - Persistir alterações
	atendimentoSalvo, err := a.repo.Atualizar(ctx, id, atendimentoAtualizado)
	if err != nil {
		return nil, fmt.Errorf("atualizar atendimento: %w", err)
	}

	return atendimentoSalvo, nil
}

// Listar recebe um filtro e retorna uma lista de atendimentos que correspondem aos critérios do filtro,
func (a *AtendimentoService) Listar(ctx context.Context, f atd.Filtro) ([]atd.Atendimento, int, atd.Filtro, error) {
	f.Normalizar()
	atdSlice, total, err := a.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, atd.Filtro{}, fmt.Errorf("listar atendimentos: %w", err)
	}
	return atdSlice, total, f, nil
}

// =====================================================================================================================
// FUNÇÕES DE VALIDAÇÃO
// =====================================================================================================================

// validarChamado recebe o ID do chamado, busca e valida o chamado.
//
// Erros sentinela possíveis: ErrChamadoNaoEncontrado, ErrChamadoStatusInvalidoParaAtendimento.
func (a *AtendimentoService) validarChamado(ctx context.Context, chamadoID string) (*chm.Chamado, error) {
	chamado, err := a.chamadoRepo.BuscarPorID(ctx, chamadoID)
	if err != nil {
		if err == mysql.ErrChamadoNaoEncontrado {
			return nil, err
		}
		return nil, fmt.Errorf("validar chamado: %w", err)
	}

	if chamado.NaoPodeSerModificado() {
		return nil, ErrChamadoStatusInvalidoParaAtendimento
	}

	return chamado, nil
}

// validarTecnico recebe o ID do técnico e da categoria, busca e valida o usuário que será atribuído.
//
// Erros sentinela possíveis: ErrTecnicoAtribuidoNaoEncontrado, ErrUsuarioNaoTemPermissaoTecnico, ErrTecnicoAtribuidoDesativado, ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) validarTecnico(ctx context.Context, tecnicoID, categoriaID string) (*usr.Usuario, error) {
	tecnico, err := a.usuarioRepo.BuscarPorID(ctx, tecnicoID)
	if err != nil {
		if err == mysql.ErrUsuarioNaoEncontrado {
			return nil, ErrTecnicoAtribuidoNaoEncontrado
		}
		return nil, fmt.Errorf("validar técnico atribuído: %w", err)
	}

	if tecnico.Permissao() != usr.PermTEC {
		return nil, ErrUsuarioNaoTemPermissaoTecnico
	}

	if !tecnico.Status() {
		return nil, ErrTecnicoAtribuidoDesativado
	}

	if err := a.validarPermissaoCategoria(ctx, categoriaID, tecnicoID); err != nil {
		return nil, fmt.Errorf("validar técnico atribuído: %w", err)
	}

	return tecnico, nil
}

// validarPermissoesAtribuicao valida recebimento de atribuição de atendimento com base no tipo de usuário autenticado.
//
// Erros sentinela possíveis: ErrTecnicoNaoPodeAtribuirParaOutro, ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) validarPermissoesAtribuicao(
	ctx context.Context,
	usrAutenticado *auth.UsuarioAutenticado,
	tecnico *usr.Usuario,
	chamado *chm.Chamado,
) error {

	if usrAutenticado.EhTEC() {
		return a.validarAtribuicaoTec(ctx, usrAutenticado, tecnico, chamado)
	}

	if usrAutenticado.EhADM() {
		return a.validarAtribuicaoAdm(ctx, usrAutenticado, chamado)
	}

	if usrAutenticado.EhDEV() {
		// Desenvolvedor não tem restrições
		return nil
	}

	return nil
}

// validarAtribuicaoTec recebe o contexto, o usuário autenticado, o técnico atribuído e o chamado, e valida a atribuição feita por um técnico.
//
// Erros sentinela possíveis: ErrTecnicoNaoPodeAtribuirParaOutro, ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) validarAtribuicaoTec(
	ctx context.Context,
	usrAutenticado *auth.UsuarioAutenticado,
	tecnico *usr.Usuario,
	chamado *chm.Chamado,
) error {
	// Técnico só pode atribuir para si mesmo
	if usrAutenticado.ID() != tecnico.ID() {
		return ErrTecnicoNaoPodeAtribuirParaOutro
	}

	// Validar permissão de acesso à categoria
	if err := a.validarPermissaoCategoria(ctx, chamado.CategoriaID(), usrAutenticado.ID()); err != nil {
		return fmt.Errorf("criar atendimento: %w", err)
	}

	return nil
}

// validarAtribuicaoAdm recebe o contexto, o usuário autenticado e o chamado, e valida a atribuição feita por um administrador.
// 
// Erros sentinela possíveis: ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) validarAtribuicaoAdm(
	ctx context.Context,
	usrAutenticado *auth.UsuarioAutenticado,
	chamado *chm.Chamado,
) error {
	// Administrador deve ter permissão na categoria
	if err := a.validarPermissaoCategoria(ctx, chamado.CategoriaID(), usrAutenticado.ID()); err != nil {
		return fmt.Errorf("criar atendimento: %w", err)
	}

	return nil
}

// validarPermissaoCategoria recebe o ID da categoria e do usuário, e valida se o usuário tem permissão para acessar a categoria.
//
// Erros sentinela possíveis: ErrPermissaoAcessoCategoriaNegada.
func (a *AtendimentoService) validarPermissaoCategoria(ctx context.Context, categoriaID, usuarioID string) error {
	_, err := a.categoriaPermissaoRepo.BuscarPorID(ctx, categoriaID, usuarioID)
	if err != nil {
		if err == mysql.ErrCategoriaPermissaoNaoEncontrada {
			return ErrPermissaoAcessoCategoriaNegada
		}
		return fmt.Errorf("validar permissão de categoria: %w", err)
	}

	return nil
}

// validarAtendimentoDuplicado recebe o ID do chamado e do técnico atribuído, e verifica se já existe um atendimento ativo para esses IDs.
//
// Erros sentinela possíveis: ErrAtendimentoAtivoJaExiste.
func (a *AtendimentoService) validarAtendimentoDuplicado(ctx context.Context, chamadoID, atribuidoID string) error {
	atendimentoAtivo, err := a.repo.BuscarPorChamadoEAtribuidoID(ctx, chamadoID, atribuidoID)
	if err != nil && err != mysql.ErrAtendimentoNaoEncontrado {
		return fmt.Errorf("validar atendimento duplicado: %w", err)
	}

	if atendimentoAtivo != nil {
		return ErrAtendimentoAtivoJaExiste
	}

	return nil
}

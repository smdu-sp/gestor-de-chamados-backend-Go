package service

import (
	"context"
	"fmt"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
)

var (
	ErrUsuarioNaoEhAutorDoChamado              = fmt.Errorf("usuário autenticado não é o autor do chamado")
	ErrTecnicoNaoAtribuidoAoChamado            = fmt.Errorf("técnico autenticado não é o atribuído ao chamado")
	ErrUsuarioNaoEhAutorDoAcompanhamento       = fmt.Errorf("usuário autenticado não é o autor do acompanhamento")
	ErrChamadoStatusInvalidoParaAcompanhamento = fmt.Errorf("chamado está em estado inválido para acompanhamento (arquivado, fechado ou rejeitado)")
)

// AcompanhamentoService representa a camada de serviço para operações relacionadas a acompanhamentos.
type AcompanhamentoService struct {
	id              dmn.GeradorID
	repo            acp.Repository
	chamadoRepo     chm.Repository
	atendimentoRepo atd.Repository
}

// NovoAcompanhamentoService cria uma nova instância de AcompanhamentoService.
func NovoAcompanhamentoService(
	id dmn.GeradorID,
	repo acp.Repository,
	chamadoRepo chm.Repository,
	atendimentoRepo atd.Repository,
) *AcompanhamentoService {
	return &AcompanhamentoService{
		id:              id,
		repo:            repo,
		chamadoRepo:     chamadoRepo,
		atendimentoRepo: atendimentoRepo,
	}
}

// Asserção de interface para garantir que AcompanhamentoService implementa AcompanhamentoService
var _ acp.Service = (*AcompanhamentoService)(nil)

// BuscarPorID busca um acompanhamento pelo seu ID.
func (a *AcompanhamentoService) BuscarPorID(ctx context.Context, id string) (*acp.Acompanhamento, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("buscar acompanhamento por ID: %w", err)
	}

	// 2 - Buscar o acompanhamento
	acompanhamento, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar acompanhamento por ID: %w", err)
	}

	// 3 - Validar se o usuário autenticado tem permissão para ler o acompanhamento
	if err := a.validarPermissaoLeitura(ctx, usrAutenticado, acompanhamento.ChamadoID()); err != nil {
		return nil, fmt.Errorf("buscar acompanhamento por ID: %w", err)
	}

	return acompanhamento, nil
}

// BuscarPorChamadoID busca acompanhamentos pelo ID do chamado.
func (a *AcompanhamentoService) BuscarPorChamadoID(ctx context.Context, id string) ([]acp.Acompanhamento, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("buscar acompanhamentos por ID do chamado: %w", err)
	}

	// 2 - Verificar se o chamado existe
	chamado, err := a.chamadoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar acompanhamentos por ID do chamado: %w", err)
	}

	// 3 - Validar se o usuário autenticado tem permissão para ler os acompanhamentos do chamado
	if err := a.validarPermissaoLeitura(ctx, usrAutenticado, chamado.ID()); err != nil {
		return nil, fmt.Errorf("buscar acompanhamentos por ID do chamado: %w", err)
	}

	// 5 - Buscar os acompanhamentos
	acompanhamentos, err := a.repo.BuscarPorChamadoID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("buscar acompanhamentos por ID do chamado: %w", err)
	}
	
	return acompanhamentos, nil
}

// Criar cria um novo acompanhamento.
func (a *AcompanhamentoService) Criar(ctx context.Context, criar acp.CriarParams) (*acp.Acompanhamento, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("criar acompanhamento: %w", err)
	}

	// 2 - Validar o chamado
	chamado, err := a.validarChamado(ctx, criar.ChamadoID)
	if err != nil {
		return nil, fmt.Errorf("criar acompanhamento: %w", err)
	}

	// 3 - Validar permissão de criação
	if err := a.validarPermissaoCriacao(ctx, usrAutenticado, chamado); err != nil {
		return nil, fmt.Errorf("criar acompanhamento: %w", err)
	}

	// 4 - Criar o ID do acompanhamento
	id, err := a.id.NovoID()
	if err != nil {
		return nil, fmt.Errorf("criar acompanhamento: %w", err)
	}

	// 5 - Criar o acompanhamento
	acompanhamento, err := acp.Novo(
		id,
		criar.ChamadoID,
		usrAutenticado.ID(),
		criar.Conteudo,
		usrAutenticado.Permissao(),
	)
	if err != nil {
		return nil, err
	}

	return a.repo.Criar(ctx, *acompanhamento)
}

// Atualizar atualiza as informações de um acompanhamento existente.
func (a *AcompanhamentoService) Atualizar(ctx context.Context, id string, atualizar acp.AtualizarParams) (*acp.Acompanhamento, error) {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return nil, fmt.Errorf("atualizar acompanhamento: %w", err)
	}

	// 2 - Buscar o acompanhamento existente
	acompanhamentoAtual, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("atualizar acompanhamento: %w", err)
	}

	// 3 - Verificar se o usuário autenticado é o autor do acompanhamento
	if usrAutenticado.ID() != acompanhamentoAtual.UsuarioID() {
		return nil, ErrUsuarioNaoEhAutorDoAcompanhamento
	}

	// 4 - Atualizar os dados do acompanhamento
	acompanhamentoAtualizado, err := acompanhamentoAtual.ComDadosAtualizados(atualizar)
	if err != nil {
		return nil, err
	}

	// 5 - Persistir as alterações no repositório
	acompanhamentoSalvo, err := a.repo.Atualizar(ctx, id, acompanhamentoAtualizado)
	if err != nil {
		return nil, fmt.Errorf("atualizar acompanhamento: %w", err)
	}

	return acompanhamentoSalvo, nil
}

// Deletar remove um acompanhamento pelo ID.
func (a *AcompanhamentoService) Deletar(ctx context.Context, id string) error {
	// 1 - Obter usuário autenticado do contexto
	usrAutenticado, err := auth.UsuarioAutenticadoDoContexto(ctx)
	if err != nil {
		return fmt.Errorf("deletar acompanhamento: %w", err)
	}

	// 2 - Buscar o acompanhamento existente
	acompanhamento, err := a.repo.BuscarPorID(ctx, id)
	if err != nil {
		return fmt.Errorf("deletar acompanhamento: %w", err)
	}

	// 3 - Verificar se o usuário autenticado é o autor do acompanhamento
	if usrAutenticado.ID() != acompanhamento.UsuarioID() {
		return ErrUsuarioNaoEhAutorDoAcompanhamento
	}

	// 4 - Deletar o acompanhamento do repositório
	if err := a.repo.Deletar(ctx, id); err != nil {
		return fmt.Errorf("deletar acompanhamento: %w", err)
	}
	return nil
}

// Listar lista acompanhamentos com paginação e filtros opcionais.
func (a *AcompanhamentoService) Listar(ctx context.Context, f acp.Filtro) ([]acp.Acompanhamento, int, acp.Filtro, error) {
	f.Normalizar()
	acompanhamentos, total, err := a.repo.Listar(ctx, f)
	if err != nil {
		return nil, 0, f, fmt.Errorf("listar acompanhamentos: %w", err)
	}
	return acompanhamentos, total, f, nil
}

// =====================================================================================================================
// FUNÇÕES DE VALIDAÇÃO
// =====================================================================================================================

// validarChamado verifica se o chamado existe e está em um estado válido para acompanhamento.
func (a *AcompanhamentoService) validarChamado(ctx context.Context, chamadoID string) (*chm.Chamado, error) {
	chamado, err := a.chamadoRepo.BuscarPorID(ctx, chamadoID)
	if err != nil {
		return nil, err
	}

	if chamado.NaoPodeSerModificado() {
		return nil, ErrChamadoStatusInvalidoParaAcompanhamento
	}

	return chamado, nil
}

// validarPermissaoLeitura verifica se o usuário autenticado tem permissão para ler acompanhamentos do chamado.
func (a *AcompanhamentoService) validarPermissaoLeitura(
	ctx context.Context,
	usrAutenticado *auth.UsuarioAutenticado,
	chamadoID string,
) error {

	chamado, err := a.chamadoRepo.BuscarPorID(ctx, chamadoID)
	if err != nil {
		return err
	}

	if usrAutenticado.EhTEC() {
		if err := a.validarTecnicoAtribuido(ctx, chamadoID, usrAutenticado.ID()); err != nil {
			return err
		}
	}

	if usrAutenticado.ID() != chamado.CriadorID() {
		return ErrUsuarioNaoEhAutorDoChamado
	}

	return nil
}

// validarPermissaoCriacao verifica se o usuário autenticado tem permissão para criar um acompanhamento no chamado.
func (a *AcompanhamentoService) validarPermissaoCriacao(
	ctx context.Context,
	usrAutenticado *auth.UsuarioAutenticado,
	chamado *chm.Chamado,
) error {

	if usrAutenticado.EhTEC() {
		return a.validarTecnicoAtribuido(ctx, chamado.ID(), usrAutenticado.ID())
	}

	if usrAutenticado.ID() != chamado.CriadorID() {
		return ErrUsuarioNaoEhAutorDoChamado
	}

	return nil
}

// validarTecnicoAtribuido verifica se o técnico está atribuído ao chamado.
func (a *AcompanhamentoService) validarTecnicoAtribuido(
	ctx context.Context,
	chamadoID, tecnicoID string,
) error {

	_, err := a.atendimentoRepo.BuscarPorChamadoEAtribuidoID(ctx, chamadoID, tecnicoID)
	if err != nil {
		if err == mysql.ErrAtendimentoNaoEncontrado {
			return ErrTecnicoNaoAtribuidoAoChamado
		}
		return err
	}

	return nil
}

package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/model"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usecase"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/repository"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/interface/response"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/utils"
)

const (
	entidadeCategoriaPermissao = "CATEGORIA_PERMISSAO"
)

// CategoriaPermissaoHandler gerencia as requisições HTTP relacionadas a categorias e permissões.
type CategoriaPermissaoHandler struct {
	Usecase    usecase.CategoriaPermissaoUsecase
	UsecaseLog usecase.LogUsecase
}

// NewCategoriaPermissaoHandler cria uma nova instância de CategoriaPermissaoHandler.
func NewCategoriaPermissaoHandler(usecase usecase.CategoriaPermissaoUsecase, usecaseLog usecase.LogUsecase) *CategoriaPermissaoHandler {
	return &CategoriaPermissaoHandler{
		Usecase:    usecase,
		UsecaseLog: usecaseLog,
	}
}

// Criar godoc
// @Summary Cria uma nova categoria e permissão
// @Description Cria uma nova categoria e permissão no sistema.
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoria_permissao body model.CategoriaPermissao true "Dados da categoria e permissão"
// @Success 201 {object} model.CategoriaPermissao
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Router /categorias_permissoes [post]
// Criar categoriaPermissao
func (h *CategoriaPermissaoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPost) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	var categoriaPermissao model.CategoriaPermissao
	if err := json.NewDecoder(r.Body).Decode(&categoriaPermissao); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.CriarCategoriaPermissao(ctx, &categoriaPermissao); err != nil {
		switch {
			// requisições inválidas - 400
			case errors.As(err, &utils.ValidacaoErrors{}):
				response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao criar categoriaPermissao", err.Error())
				return

			// conflitos - 409
			case errors.Is(err, repository.ErrCategoriaPermissaoJaExiste):
				response.ErrorJSON(w, http.StatusConflict, "categoriaPermissao já cadastrada", err.Error())
				return

			// erros internos - 500
			case errors.Is(err, utils.ErrUUIDv7Generation),
				errors.Is(err, repository.ErrExecContext),
				errors.Is(err, repository.ErrRowsAffected):
				response.ErrorJSON(w, http.StatusInternalServerError, "erro ao criar categoriaPermissao", err.Error())
				return

			// erro de contexto - 408
			case errors.Is(err, context.DeadlineExceeded):
				response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao criar categoriaPermissao", err.Error())
				return

			// erro de contexto - 400
			case errors.Is(err, context.Canceled):
				response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao criar categoriaPermissao", err.Error())
				return

			// erro desconhecido - 500
			default:
				response.ErrorJSON(w, http.StatusInternalServerError, "erro desconhecido ao criar categoriaPermissao", err.Error())
				return
		}
	}

		err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoCriar,
		entidadeCategoriaPermissao,
		fmt.Sprintf("CategoriaPermissao criada via API: %s", categoriaPermissao.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, response.ToCategoriaPermissaoResponse(&categoriaPermissao))
}

// BuscarTudo godoc
// @Sumary Listar todas as categorias e permissões
// @Description Lista todas as categorias e permissões com base em filtros e paginação.
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param pagina query int false "Página"
// @Param limite query int false "Limite"
// @Param categoriaId query string false "ID da categoria"
// @Param usuarioId query string false "ID do usuário"
// @Param permissao query string false "permissão"
// @Success 200 {object} []model.CategoriaPermissao
// @Failure 400 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias_permissoes [get]
// BuscarTudo lista todas as categorias e permissões com base em filtros e paginação.
func (h *CategoriaPermissaoHandler) BuscarTudo(w http.ResponseWriter, r *http.Request) {
		if !metodoHttpValido(w, r, http.MethodGet) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	query := r.URL.Query()

	filtro := model.CategoriaPermissaoFiltro{}

	if pagina, err := strconv.Atoi(query.Get("pagina")); err == nil {
		filtro.Pagina = pagina
	}
	if limite, err := strconv.Atoi(query.Get("limite")); err == nil {
		filtro.Limite = limite
	}
	if categoriaID := query.Get("categoriaId"); categoriaID != "" {
		filtro.CategoriaID = &categoriaID
	}
	if usuarioID := query.Get("usuarioId"); usuarioID != "" {
		filtro.UsuarioID = &usuarioID
	}
	if permissao := query.Get("permissao"); permissao != "" {
		filtro.Permissao = &permissao
	}

	items, total, filtroCorrigido, err := h.Usecase.ListarCategoriaPermissao(ctx, filtro)
	if err != nil {
		switch {
			// erros internos - 500
		case errors.Is(err, repository.ErrQueryContext),
			errors.Is(err, repository.ErrScannerCategoriaPermissao),
			errors.Is(err, repository.ErrScan):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao listar categoria permissao", err.Error())
			return

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao listar categoria permissao", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao listar categoria permissao", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao listar categoria permissao", err.Error())
			return
		}
	}
		response.JSON(w, http.StatusOK, response.PageResponse[model.CategoriaPermissao]{
		Total:  total,
		Pagina: filtroCorrigido.Pagina,
		Limite: filtroCorrigido.Limite,
		Items:  items,
	})
}

// Atualizar godoc
// @Summary Atualiza uma categoria e permissão existente
// @Description Atualiza os dados de uma categoria e permissão existente no sistema.
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoriaId path string true "ID da categoria"
// @Param usuarioId path string true "ID do usuário"
// @Param categoria_permissao body model.CategoriaPermissao true "Dados atualizados da categoria e permissão"
// @Success 200 {object} map[string]string
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 409 {object} any
// @Failure 500 {object} any
// @Router /categorias_permissoes/{categoriaId}/usuarios/{usuarioId} [put]
// Atualizar modifica os dados de uma categoria e permissão existente.
func (h *CategoriaPermissaoHandler) Atualizar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodPut) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := obterParametroURL(r, "categoriaId")
	usuarioID := obterParametroURL(r, "usuarioId")

	var categoriaPermissao model.CategoriaPermissao
	if err := json.NewDecoder(r.Body).Decode(&categoriaPermissao); err != nil {
		response.ErrorJSON(w, http.StatusBadRequest, payloadInvalidoMsg, err.Error())
		return
	}

	if err := h.Usecase.AtualizarCategoriaPermissao(ctx, id, usuarioID, &categoriaPermissao); err != nil {
		switch {
			// requisições inválidas - 400
		case errors.Is(err, model.ErrCategoriaPermissaoIDInvalido),
			errors.Is(err, model.ErrNomeInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "dados inválidos ao atualizar categoriaPermissao", err.Error())

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrCategoriaPermissaoNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "ID inválido ao atualizar categoria", err.Error())

		// conflitos - 409
		case errors.Is(err, repository.ErrCategoriaPermissaoJaExiste):
			response.ErrorJSON(w, http.StatusConflict, "categoriaPermissao já cadastrada", err.Error())

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao atualizar categoriaPermissao", err.Error())

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao atualizar categoriaPermissao", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao atualizar categoriaPermissao", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao atualizar categoriaPermissao", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoAtualizar,
		entidadeCategoriaPermissao,
		fmt.Sprintf("CategoriaPermissao atualizada via API: %s", categoriaPermissao.String()),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "categoriaPermissao atualizada com sucesso"})
}

// Deletar godoc
// @Summary Deleta uma categoriaPermissao existente
// @Description Remove uma categoriaPermissao do sistema.
// @Tags CategoriaPermissao
// @Accept json
// @Produce json
// @Param categoriaId path string true "ID da categoria"
// @Param usuarioId path string true "ID do usuário"
// @Success 200 {object} map[string]string
// @Failure 400 {object} any
// @Failure 404 {object} any
// @Failure 405 {object} any
// @Failure 408 {object} any
// @Failure 500 {object} any
// @Router /categorias_permissoes/{categoriaId}/usuarios/{usuarioId} [delete]
// Deletar remove uma categoriaPermissao do sistema.
func (h *CategoriaPermissaoHandler) Deletar(w http.ResponseWriter, r *http.Request) {
	if !metodoHttpValido(w, r, http.MethodDelete) {
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), timeoutPadrao)
	defer cancel()

	id := obterParametroURL(r, "categoriaId")
	usuarioID := obterParametroURL(r, "usuarioId")

	if err := h.Usecase.DeletarCategoriaPermissao(ctx, id, usuarioID); err != nil {
		switch {
			// requisições inválidas - 400
		case errors.Is(err, model.ErrCategoriaPermissaoIDInvalido):
			response.ErrorJSON(w, http.StatusBadRequest, "ID inválido ao deletar categoriaPermissao", err.Error())

		// recurso não encontrado - 404
		case errors.Is(err, repository.ErrCategoriaPermissaoNaoEncontrada):
			response.ErrorJSON(w, http.StatusNotFound, "categoriaPermissao não encontrada ao deletar", err.Error())

		// erros internos - 500
		case errors.Is(err, repository.ErrExecContext),
			errors.Is(err, repository.ErrRowsAffected):
			response.ErrorJSON(w, http.StatusInternalServerError, "erro interno ao deletar categoriaPermissao", err.Error())

		// erros de contexto - 408
		case errors.Is(err, context.DeadlineExceeded):
			response.ErrorJSON(w, http.StatusRequestTimeout, "tempo de requisição excedido ao deletar categoriaPermissao", err.Error())
			return

		// erros de contexto - 400
		case errors.Is(err, context.Canceled):
			response.ErrorJSON(w, http.StatusBadRequest, "requisição cancelada ao deletar categoriaPermissao", err.Error())
			return

		// fallback de segurança - 500
		default:
			response.ErrorJSON(w, http.StatusInternalServerError, "erro inesperado ao deletar categoriaPermissao", err.Error())
			return
		}
	}

	err := h.UsecaseLog.CriarLog(
		ctx,
		model.AcaoDeletar,
		entidadeCategoriaPermissao,
		fmt.Sprintf("CategoriaPermissao deletada via API: categoriaID=%s, usuarioID=%s", id, usuarioID),
	)
	if err != nil {
		response.ErrorJSON(w, http.StatusInternalServerError, erroLogMsg, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "categoriaPermissao deletada com sucesso"})
}




// metodos auxiliares

// obterParametroURL extrai o valor de um parâmetro da URL da requisição.
func obterParametroURL(r *http.Request, nome string) string {
	vars := r.Context().Value("vars")
	if varsMap, ok := vars.(map[string]string); ok {
		return varsMap[nome]
	}
	return ""
}
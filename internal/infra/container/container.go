package container

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/auth"
	cfg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/config"
	dmn "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain"
	acp "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/acompanhamento"
	atd "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/atendimento"
	ctg "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria"
	cpm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/categoria_permissao"
	chm "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/chamado"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/log"
	subc "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/subcategoria"
	usr "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/domain/usuario"
	hdl "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/handler"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/http/router"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/id"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/migrations"
	mgt "github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/migrations"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/infra/mysql"
	"github.com/smdu-sp/gestor-de-chamados-backend-Go/internal/service"
)

// Container centraliza a criação e gerenciamento de dependências.
type Container struct {
	db             *sql.DB
	handlers       *router.Handlers
	jwtService     auth.JWTService
	usuarioService usr.Service
	migrador       mgt.Migrador
}

// NovoContainer cria e inicializa o container com todas as dependências.
func NovoContainer(config *cfg.Config, db *sql.DB) (*Container, error) {
	slog.Info("iniciando criação do container de dependências")

	if db == nil {
		return nil, fmt.Errorf("conexão com o banco de dados é nil")
	}

	// 1. Camada de Migrações (Infraestrutura)
	migrador := mgt.NovoMySQLMigrador(db, migrations.FS)

	// 2. Camada de Repositórios (Infraestrutura)
	repos := construirRepositories(db)

	// 3. Camada de Serviços (Domínio)
	svcs := construirServices(repos)

	// 4. Camada de Autenticação (Infraestrutura/Aplicação)
	authSvcs, err := construirAuthServices(config, svcs.usuario, svcs.log)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar serviços de autenticação: %w", err)
	}

	// 5. Camada de Handlers (Apresentação/Interface)
	handlers := construirHandlers(svcs, authSvcs)

	slog.Info("container de dependências criado com sucesso")

	return &Container{
		db:             db,
		handlers:       handlers,
		jwtService:     authSvcs.jwt,
		usuarioService: svcs.usuario,
		migrador:       migrador,
	}, nil
}

// ExecutarMigracoes executa as migrações pendentes do banco de dados.
// Deve ser chamado após a criação do container e antes de iniciar o servidor.
func (c *Container) ExecutarMigracoes(ctx context.Context) error {
	slog.Info("iniciando execução de migrações")
	
	if err := c.migrador.Migrar(ctx); err != nil {
		return fmt.Errorf("erro ao executar migrações: %w", err)
	}
	
	slog.Info("migrações executadas com sucesso")
	return nil
}

// Handlers retorna os handlers configurados.
// Método único de acesso às dependências de apresentação.
func (c *Container) Handlers() *router.Handlers {
	return c.handlers
}

// JWTService retorna o serviço JWT.
// Necessário para middlewares de autenticação no roteador.
func (c *Container) JWTService() auth.JWTService {
	return c.jwtService
}

// UsuarioService retorna o serviço de usuário.
// Necessário para middlewares de autenticação no roteador.
func (c *Container) UsuarioService() usr.Service {
	return c.usuarioService
}

// FecharRecursos encerra os recursos gerenciados pelo container.
func (c *Container) FecharRecursos(l *slog.Logger) {
	l.Info("iniciando fechamento de recursos")

	if c.db != nil {
		if err := c.db.Close(); err != nil {
			l.Error("erro ao fechar conexão com banco", "error", err)
		} else {
			l.Info("conexão com banco fechada com sucesso")
		}
	}

	l.Info("todos os recursos foram finalizados")
}

// ===================================================================
// Tipos privados para organização interna
// ===================================================================

// Alias: container.repositorios agrupa todos os repositórios.
type repositories struct {
	log                log.Repository
	usuario            usr.Repository
	chamado            chm.Repository
	categoria          ctg.Repository
	atendimento        atd.Repository
	subcategoria       subc.Repository
	acompanhamento     acp.Repository
	categoriaPermissao cpm.Repository
}

// Alias: container.servicos agrupa todos os serviços.
type services struct {
	log                log.Service
	usuario            usr.Service
	chamado            chm.Service
	categoria          ctg.Service
	atendimento        atd.Service
	subcategoria       subc.Service
	acompanhamento     acp.Service
	categoriaPermissao cpm.Service
}

// Alias: container.authServices agrupa todos os serviços de autenticação.
type authServices struct {
	jwt  auth.JWTService
	ldap auth.AuthExternoService
	auth auth.AuthInternoService
}

// ===================================================================
// Funções construtoras privadas
// ===================================================================

// construirRepositories cria todas as instâncias de repositórios.
func construirRepositories(db *sql.DB) *repositories {
	return &repositories{
		log:                mysql.NovoLogRepository(db),
		usuario:            mysql.NovoUsuarioRepository(db),
		chamado:            mysql.NovoChamadoRepository(db),
		categoria:          mysql.NovoCategoriaRepository(db),
		atendimento:        mysql.NovoAtendimentoRepository(db),
		subcategoria:       mysql.NovoSubcategoriaRepository(db),
		acompanhamento:     mysql.NovoAcompanhamentoRepository(db),
		categoriaPermissao: mysql.NovoCategoriaPermissaoRepository(db),
	}
}

// construirServices cria todas as instâncias de serviços.
func construirServices(r *repositories) *services {
	// Dependência para geração de IDs
	var id dmn.GeradorID = id.Gerador{}

	// 1- Serviços com dependência única
	logSvc := service.NovoLogService(id, r.log)
	usuarioSvc := service.NovoUsuarioService(id, r.usuario)
	categoriaSvc := service.NovoCategoriaService(id, r.categoria)
	subcategoriaSvc := service.NovoSubcategoriaService(id, r.subcategoria)
	acompanhamentoSvc := service.NovoAcompanhamentoService(id, r.acompanhamento, r.chamado, r.atendimento)
	categoriaPermissaoSvc := service.NovoCategoriaPermissaoService(r.categoriaPermissao)

	// 2- Serviços com múltiplas dependências
	chamadoSvc := service.NovoChamadoService(
		id,
		r.chamado,
		r.categoria,
		r.subcategoria,
		r.categoriaPermissao,
		r.atendimento,
		r.usuario,
	)

	atendimentoSvc := service.NovoAtendimentoService(
		id,
		r.atendimento,
		r.chamado,
		r.categoriaPermissao,
		r.usuario,
	)

	return &services{
		log:                logSvc,
		usuario:            usuarioSvc,
		chamado:            chamadoSvc,
		categoria:          categoriaSvc,
		atendimento:        atendimentoSvc,
		subcategoria:       subcategoriaSvc,
		acompanhamento:     acompanhamentoSvc,
		categoriaPermissao: categoriaPermissaoSvc,
	}
}

// construirAuthServices cria todas as instâncias de serviços de autenticação.
func construirAuthServices(
	config *cfg.Config,
	usuarioSvc usr.Service,
	logSvc log.Service,
) (*authServices, error) {

	// 1- Validações básicas
	if usuarioSvc == nil {
		return nil, fmt.Errorf("serviço de usuário é nil")
	}
	if logSvc == nil {
		return nil, fmt.Errorf("serviço de log é nil")
	}

	// 2- Parsear durações de TTL dos tokens
	tTTL, err := time.ParseDuration(config.TokenTTL())
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear TokenTTL: %w", err)
	}
	rtTTL, err := time.ParseDuration(config.RefreshTokenTTL())
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear RefreshTokenTTL: %w", err)
	}

	// 3- Segredos dos tokens
	tSegredo := []byte(config.TokenSecret())
	rtSegredo := []byte(config.RefreshTokenSecret())

	// 4- Criar instâncias dos serviços de autenticação
	ldapSvc := auth.NovoLDAPService(config)
	jwtSvc := auth.NovoTokenJWTService(tSegredo, rtSegredo, tTTL, rtTTL)
	authSvc := auth.NovoAuthInterno(usuarioSvc, jwtSvc, ldapSvc, logSvc, config)

	return &authServices{
		jwt:  jwtSvc,
		ldap: ldapSvc,
		auth: authSvc,
	}, nil
}

// construirHandlers cria todas as instâncias de handlers.
func construirHandlers(s *services, a *authServices) *router.Handlers {
	return &router.Handlers{
		Log:                hdl.NovoLogHandler(s.log),
		Auth:               auth.NovoAuthHandler(a.auth),
		Usuario:            hdl.NovoUsuarioHandler(s.usuario, a.auth, a.ldap, s.log),
		Chamado:            hdl.NovoChamadoHandler(s.chamado, s.log),
		Categoria:          hdl.NovoCategoriaHandler(s.categoria, s.log),
		Atendimento:        hdl.NovoAtendimentoHandler(s.atendimento, s.log),
		Subcategoria:       hdl.NovoSubcategoriaHandler(s.subcategoria, s.log),
		Acompanhamento:     hdl.NovoAcompanhamentoHandler(s.acompanhamento, s.log),
		CategoriaPermissao: hdl.NovoCategoriaPermissaoHandler(s.categoriaPermissao, s.log),
	}
}

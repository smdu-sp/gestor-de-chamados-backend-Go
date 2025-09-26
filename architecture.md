## Estrutura de Diretórios do Gestor de Chamados

```markdown
/gestor-de-chamados-backend-Go
│
├── /cmd/               # Entrypoints da aplicação
│ ├── /api/             # API HTTP REST
│ └── main.go           # Ponto de entrada da API
│
├── /internal/          # Código interno da aplicação (não exposto externamente)
│ ├── /config/          # Configuração (env, flags, setup)
│ │ └── config.go
│ │
│ ├── /domain/          # Regras de negócio (Enterprise Rules)
│ │ ├── /model/         # Entidades do domínio
│ │ ├── /repository/    # Interfaces de repositório
│ │ └── /usecase/       # Interfaces de casos de uso
│ │
│ ├── /usecase/         # Regras de aplicação (Application Business Rules)
│ │ ├── /user/          # Implementação de casos de uso de usuários
│ │ └── /auth/          # Casos de uso de autenticação/autorização
│ │
│ ├── /interface/       # Adaptadores externos (Interface Adapters)
│ │ ├── /handler/       # HTTP Handlers
│ │ ├── /presenter/     # Formatação de respostas, DTOs
│ │ ├── /gateway/       # Implementações de componentes externos (DB, APIs externas)
│ │ └── /router/        # Inicialização de rotas
│ │
│ ├── /infra/           # Infraestrutura
│ │ ├── /db/            # Conexões com BD
│ │ └── /repository/    # Implementações dos repositórios
│ │
│ ├── /auth/            # Infraestrutura de auth (JWT, OAuth, middleware)
│ │ ├── /middleware/    # Autenticação/autorização via middleware
│ │ ├── /jwt/           # Lógica de JWT
│ │ └── /provider/      # OAuth2, Auth0, etc.
│ │
│ ├── /middleware/      # Middlewares HTTP (logging, recovery, CORS, rate limiting)
│ │
│ ├── /email/           # Envio de emails e templates
│ │
│ ├── /job/             # Workers e tarefas assíncronas
│ │
│ └── /utils/           # Utilitários gerais e helpers
│
├── /pkg/               # Bibliotecas reutilizáveis por outros projetos
│ └── /logger/          # Logger configurável (zap, logrus etc.)
│
├── /migrations/        # Scripts de migração de banco de dados
│ └── V001_init.sql
│
├── /tests/             # Testes de integração, e2e, mocks
│ ├── /e2e/             # End-to-end tests
│ └── /mocks/           # Mocks e fakes para testes
│
├── /scripts/           # Scripts utilitários (reset DB, gerar dados fake, lint, etc.)
│ └── reset_db.sh
│
├── /docs/              # Documentação geral do sistema (Swagger e etc.)
│ ├── architecture.md
│ ├── docs.go
│ └── swagger.json
│
├── /deploy/            # Infraestrutura como código (Docker, K8s, Terraform, etc.)
│ ├── Dockerfile
│ ├── docker-compose.yml
│ └── k8s/
│     └── deployment.yaml
│
├── /hack/              # Scripts e ferramentas internas de dev (build helpers, local setup)
│ └── devtools.sh
│
├── .env                # Variáveis de ambiente locais
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```
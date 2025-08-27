# Gestor de Chamados - Backend

## API REST em Go (puro) com JWT + LDAP + RBAC

Este projeto é uma **API REST desenvolvida em Go (sem frameworks web)** utilizando apenas `net/http`, com suporte a:

* **Autenticação via LDAP ou Active Directory**
* **Emissão de tokens JWT**
* **Refresh Token**
* **Controle de acesso baseado em papéis (RBAC)**
* **Middlewares de autenticação e autorização**

---

## Estrutura do projeto

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/           # JWT, LDAP/AD, Middleware
│   ├── config/         # Configurações
│   ├── domain/         # Model, Repository, Service 
│   ├── http/           # Handlers, Routers
│   ├── response/       # JSON, ErrorJSON
│   └── util/           # Utilitários
├── ldap-init/          # LDIF de testes OpenLDAP
├── migrations/         # Migrations do BD
├── .env
├── docker-compose.yml  # Sobe banco e servido OpenLDAP para testes
├── go.mod
├── go.sum
└── README.md
```

---

## Tecnologias

* [Go 1.25.0](https://golang.org/)
* [net/http](https://pkg.go.dev/net/http)
* [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
* [go-ldap/ldap/v3](https://github.com/go-ldap/ldap)

---

## Configuração

Crie um arquivo `.env` na raiz do projeto:

```ini
# ============= App =============
ENVIRONMENT=production
PORT=8080
CORS_ORIGIN=http://localhost:8080 # modificar dominio em produção

# ============= DB (MySQL 8) ============
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=user
DB_PASS=troque-pela-senha-do-db
DB_NAME=troque-pelo-nome-do-db 

# ============= JWT =====================
JWT_SECRET=troque-esta-chave 
RT_SECRET=troque-esta-chave-refresh 
ACCESS_TTL=24h
REFRESH_TTL=168h # 7 dias

# ============= LDAP / OpenLDAP (Desenvolvimento) ============
# LDAP_SERVER=ldap:
# LDAP_BASE=dc=,dc=
# LDAP_USER=cn=,dc=,dc=
# LDAP_PASS=
# LDAP_LOGIN_ATTR=uid

# ============= Active Directory (Produção) ==================
LDAP_SERVER=
LDAP_DOMAIN=
LDAP_BASE=DC=,DC=
LDAP_USER=
LDAP_PASS=
LDAP_LOGIN_ATTR=sAMAccountName
```

> O sistema suporta autenticação via **OpenLDAP** (ideal para testes locais) ou **Active Directory** (produção).

---

# AD (exemplo)
```bash
ldapsearch -x -H ldap://10.10.65.242 \
  -D "usuario@rede.sp" -w "senha" \
  -b "DC=rede,DC=sp"
```

---

## Executando a aplicação

```bash
# Instalar dependências
go mod tidy

# Rodar API
go run ./cmd/api
```

A API estará disponível em:
`http://localhost:8080`

---

## Endpoints principais

### Autenticação

**POST /login**

```json
Request:
{ "login": "usuario", "password": "senha@" }

Response 200:
{
    "access_token": <access-token>,
    "refresh_token": <refresh_token>
}
```

* `401 Unauthorized`: credenciais inválidas

---

### Refresh token

**POST /refresh**

```json
Request:
{ "refreshToken": "<refresh-token>" }

Response 200:
{
  "access_token": "<novo-access>",
  "refresh_token": "<novo-refresh>"
}
```

---

## Testando rapidamente com `curl`

```bash
# Público
curl -s http://localhost:8080/

# Login
curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"usuario","password":"senha@"}'

# Refresh
curl -s -X POST http://localhost:8080/refresh \
  -H 'Content-Type: application/json' \
  -d '{"refreshToken":"<refresh-token>"}'
```


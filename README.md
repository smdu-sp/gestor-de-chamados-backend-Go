# Gestor de Chamados - Backend

Este projeto é uma **API REST desenvolvida em Go (sem frameworks web)** utilizando apenas `net/http`, com suporte a:

* **Autenticação via LDAP ou Active Directory**
* **Emissão de tokens JWT**
* **Refresh Token**
* **Controle de acesso baseado em papéis (RBAC)**
* **Middlewares de autenticação e autorização**

---

## Tecnologias

* [Go 1.25.0](https://golang.org/)
* [net/http](https://pkg.go.dev/net/http)
* [go-ldap/ldap/v3](https://github.com/go-ldap/ldap)

---

## Configuração

Antes de rodar o projeto, configure as variáveis de ambiente:

Copie o arquivo de exemplo:
```bash
  cp .env.example .env
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


# Gestor de Chamados - Backend

## API REST em Go (puro) com JWT + LDAP + RBAC

Este projeto é uma **API REST desenvolvida em Go (sem frameworks web)** utilizando apenas `net/http`, com suporte a:

* **Autenticação via LDAP**
* **Emissão de tokens JWT**
* **Controle de acesso baseado em papéis (RBAC)**
* **Middlewares de autenticação, autorização, recover e timeout**
* Repositório de usuários **in-memory** (com interface pronta para extensão em banco SQL, ex.: Postgres)

---

## Estrutura do projeto

```text
.
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── auth/           # JWT + RBAC
│   ├── config/         # Configurações
│   ├── handler/        # Handlers HTTP
│   ├── httpx/          # Middleware e roteamento
│   ├── ldapx/          # Cliente LDAP
│   ├── model/          # Modelos de domínio
│   ├── repository/     # Interfaces e memória
│   └── util/           # Utilitários
├── test/               # Testes
│   └── ldap/           # Servidor LDAP de testes 
├── .env
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
HTTP_ADDR=:8080
JWT_SECRET=troque-por-um-segredo-forte
JWT_ISSUER=authapi
JWT_TTL=2h

LDAP_URL=ldap://localhost:389
LDAP_BASE_DN=dc=minhaempresa,dc=local
LDAP_BIND_DN=cn=admin,dc=minhaempresa,dc=local
LDAP_BIND_PASS=admin
LDAP_USER_FILTER=(|(uid=%s)(sAMAccountName=%s)(mail=%s))

LDAP_TLS=false
LDAP_INSECURE=true
LDAP_ATTR_LOGIN=uid
LDAP_ATTR_NAME=cn
LDAP_ATTR_EMAIL=mail
LDAP_ATTR_AVATAR=
LDAP_ATTR_PERM=department
```
## Subindo servidor LDAP de teste com Docker

Para testar a API localmente, você pode usar um container **OpenLDAP** simples:

```bash
docker run --name ldap-server -p 389:389 \
  -e LDAP_ORGANISATION="Minha Empresa" \
  -e LDAP_DOMAIN="minhaempresa.local" \
  -e LDAP_ADMIN_PASSWORD="admin" \
  -d osixia/openldap:1.5.0
```

Credenciais padrão:

* **Bind DN:** `cn=admin,dc=minhaempresa,dc=local`
* **Senha:** `admin`

---

### Inserindo usuários de teste

Crie um arquivo `usuarios.ldif`:

```ldif
dn: uid=jdoe,dc=minhaempresa,dc=local
objectClass: inetOrgPerson
uid: jdoe
sn: Doe
cn: John Doe
mail: jdoe@minhaempresa.local
userPassword: s3nh@
```

Carregue o usuário no LDAP:

```bash
docker cp usuarios.ldif ldap-server:/usuarios.ldif
docker exec -it ldap-server ldapadd -x \
  -D "cn=admin,dc=minhaempresa,dc=local" -w admin \
  -f /usuarios.ldif
```

---

### Testando o LDAP

```bash
ldapsearch -x -H ldap://localhost:389 \
  -D "cn=admin,dc=minhaempresa,dc=local" -w admin \
  -b "dc=minhaempresa,dc=local"
```

Se o usuário `jdoe` aparecer, o ambiente LDAP está pronto para ser usado pela API.

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
{ "login": "jdoe", "password": "s3nh@" }

Response 200:
{
  "token": "<jwt>",
  "user": {
    "id": "123",
    "nome": "John Doe",
    "login": "jdoe",
    "email": "jdoe@empresa.com",
    "permissao": "USR"
  }
}
```

* `401 Unauthorized`: credenciais inválidas

---

### Usuário autenticado

**GET /me**
(Requer header `Authorization: Bearer <token>`)

Retorna os dados do usuário logado (claims do JWT).

---

### Rotas com RBAC

* **GET /admin/ping** → requer papel `ADM`
* **GET /suporte/ping** → requer papel `SUP` ou `DEV`
* **GET /** → rota pública (sem autenticação)

---

## Testando rapidamente com `curl`

```bash
# Público
curl -s http://localhost:8080/

# Login
curl -s -X POST http://localhost:8080/login \
  -H 'Content-Type: application/json' \
  -d '{"login":"jdoe","password":"s3nh@"}'

# Claims
TOKEN="<cole-o-token>"
curl -s http://localhost:8080/me -H "Authorization: Bearer $TOKEN"

# Rota admin
curl -s http://localhost:8080/admin/ping -H "Authorization: Bearer $TOKEN"
```


# Cria diretórios e arquivos conforme a estrutura da arquitetura limpa

# Função auxiliar para criar arquivos vazios
function New-EmptyFile($path) {
  if (-not (Test-Path $path)) {
    New-Item -ItemType File -Path $path | Out-Null
  }
}

# Estrutura de diretórios
$directories = @(
  "cmd/api",
  "internal/config",
  "internal/domain/model",
  "internal/domain/repository",
  "internal/domain/usecase",
  "internal/usecase/user",
  "internal/usecase/auth",
  "internal/interface/handler",
  "internal/interface/response",
  "internal/interface/gateway",
  "internal/interface/router",
  "internal/infra/db",
  "internal/infra/repository",
  "internal/auth/middleware",
  "internal/auth/jwt",
  "internal/auth/provider",
  "internal/middleware",
  "internal/validator",
  "internal/email",
  "internal/job",
  "internal/utils",
  "pkg/logger",
  "api",
  "migrations",
  "tests/e2e",
  "tests/mocks",
  "scripts",
  "docs",
  "deploy/k8s",
  "hack"
)

# Criando diretórios
foreach ($dir in $directories) {
  New-Item -ItemType Directory -Path $dir -Force | Out-Null
}

# Arquivos principais (vazios por enquanto)
$files = @(
  "cmd/api/main.go",
  "internal/config/config.go",
  "migrations/V001_init.sql",
  "scripts/reset_db.sh",
  "docs/architecture.md",
  "docs/api.md",
  "deploy/Dockerfile",
  "deploy/docker-compose.yml",
  "hack/devtools.sh",
  ".env",
  ".gitignore",
  "go.mod",
  "go.sum",
  "README.md"
)

# Criando arquivos vazios
foreach ($file in $files) {
  New-EmptyFile $file
}

Write-Host "Estrutura de projeto criada com sucesso!"

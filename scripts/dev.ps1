param(
  [Parameter(Mandatory = $true)]
  [ValidateSet('migrate', 'server', 'build', 'help')]
  [string]$Command
)

switch ($Command) {
  'migrate' {
    Write-Host "Executando migrations..." -ForegroundColor Cyan
    go run .\cmd\api migrate
  }
  'server' {
    Write-Host "Iniciando servidor..." -ForegroundColor Cyan
    go run .\cmd\api serve
  }
  'build' {
    Write-Host "Compilando binário..." -ForegroundColor Cyan
    if (!(Test-Path -Path "bin")) {
      New-Item -ItemType Directory -Path "bin"
    }
    go build -o bin\api.exe .\cmd\api
    if ($LASTEXITCODE -eq 0) {
      Write-Host "Binário: bin\api.exe" -ForegroundColor Green
    }
  }
  'help' {
    Write-Host @"
Comandos disponíveis:
  .\dev.ps1 migrate  - Executa migrations
  .\dev.ps1 server    - Inicia servidor
  .\dev.ps1 build    - Compila binário
  .\dev.ps1 help     - Exibe esta ajuda
"@ -ForegroundColor Yellow
  }
}
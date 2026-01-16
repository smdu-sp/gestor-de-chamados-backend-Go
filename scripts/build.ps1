Write-Host "Compilando binário..." -ForegroundColor Cyan

# Criar diretório bin se não existir
if (!(Test-Path -Path "bin")) {
  New-Item -ItemType Directory -Path "bin"
}

go build -o bin\api.exe .\cmd\api

if ($LASTEXITCODE -eq 0) {
  Write-Host "Binário criado em: bin\api.exe" -ForegroundColor Green
}
else {
  Write-Host "Erro ao compilar" -ForegroundColor Red
  exit 1
}
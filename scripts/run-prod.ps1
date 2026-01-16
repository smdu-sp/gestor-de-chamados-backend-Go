Write-Host "Modo Produção" -ForegroundColor Yellow

# 1. Build
.\scripts\build.ps1
if ($LASTEXITCODE -ne 0) { exit 1 }

# 2. Migrations
Write-Host "`nExecutando migrations..." -ForegroundColor Cyan
.\bin\api.exe migrate
if ($LASTEXITCODE -ne 0) { 
  Write-Host "Falha nas migrations" -ForegroundColor Red
  exit 1 
}

# 3. Serve
Write-Host "`nIniciando servidor..." -ForegroundColor Cyan
.\bin\api.exe serve
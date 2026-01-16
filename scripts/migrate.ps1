Write-Host "Executando migrations..." -ForegroundColor Cyan
go run .\cmd\api migrate

if ($LASTEXITCODE -eq 0) {
  Write-Host "Migrations executadas com sucesso!" -ForegroundColor Green
}
else {
  Write-Host "Erro ao executar migrations" -ForegroundColor Red
  exit 1
}
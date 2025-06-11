# Script de deploy para o Quiz do Amor - Windows PowerShell
# Para usar com Cloudflare e Docker

param(
    [switch]$SkipBuild,
    [switch]$SkipHealthCheck
)

Write-Host "🚀 Iniciando deploy do Quiz do Amor..." -ForegroundColor Green

# Verificar se Docker está rodando
try {
    docker info | Out-Null
} catch {
    Write-Host "❌ Docker não está rodando. Inicie o Docker Desktop e tente novamente." -ForegroundColor Red
    exit 1
}

# Verificar se .env existe
if (-not (Test-Path ".env")) {
    Write-Host "⚠️  Arquivo .env não encontrado. Criando arquivo de exemplo..." -ForegroundColor Yellow
    
    # Gerar JWT secret
    $jwtSecret = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((New-Guid).ToString() + (Get-Date).Ticks))
    
    @"
# Configurações de Produção
GIN_MODE=release
DATABASE_PATH=/app/data/quiz.db
ADMIN_PASSWORD=SuaSenhaSeguraAqui123!
JWT_SECRET=$jwtSecret
PORT=8080

# Configurações do Cloudflare (opcional)
CLOUDFLARE_API_TOKEN=seu_token_aqui
CLOUDFLARE_ZONE_ID=seu_zone_id_aqui
DOMAIN=seudominio.com
"@ | Out-File -FilePath ".env" -Encoding UTF8
    
    Write-Host "✅ Arquivo .env criado. Configure as variáveis antes de continuar." -ForegroundColor Green
    exit 1
}

if (-not $SkipBuild) {
    Write-Host "📦 Construindo imagem Docker..." -ForegroundColor Green
    docker build -t namorada-quiz:latest .
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Falha ao construir imagem Docker" -ForegroundColor Red
        exit 1
    }
}

Write-Host "🛑 Parando containers existentes..." -ForegroundColor Green
docker-compose down

Write-Host "🚀 Subindo nova versão..." -ForegroundColor Green
docker-compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Falha ao subir containers" -ForegroundColor Red
    exit 1
}

if (-not $SkipHealthCheck) {
    Write-Host "⏳ Aguardando aplicação inicializar..." -ForegroundColor Green
    Start-Sleep -Seconds 10

    # Health check
    Write-Host "🔍 Verificando saúde da aplicação..." -ForegroundColor Green
    try {
        $response = Invoke-WebRequest -Uri "http://localhost:8080/" -TimeoutSec 10
        if ($response.StatusCode -eq 200) {
            Write-Host "✅ Aplicação está rodando corretamente!" -ForegroundColor Green
        } else {
            throw "Status code: $($response.StatusCode)"
        }
    } catch {
        Write-Host "❌ Aplicação não respondeu ao health check: $_" -ForegroundColor Red
        Write-Host "📋 Logs do container:" -ForegroundColor Yellow
        docker-compose logs quiz-app
        exit 1
    }
}

Write-Host "🎉 Deploy concluído com sucesso!" -ForegroundColor Green
Write-Host "🌐 Aplicação disponível em: http://localhost:8080" -ForegroundColor Green

# Verificar se domain está configurado
$envContent = Get-Content ".env" | ForEach-Object {
    if ($_ -match "DOMAIN=(.+)" -and $Matches[1] -ne "seudominio.com") {
        Write-Host "🌍 Domínio configurado: https://$($Matches[1])" -ForegroundColor Green
    }
}

Write-Host "`n📝 Próximos passos:" -ForegroundColor Yellow
Write-Host "1. Configure seu proxy reverso (Nginx/Cloudflare Tunnel)"
Write-Host "2. Configure SSL/TLS"
Write-Host "3. Teste a aplicação thoroughly"
Write-Host "4. Configure backup do banco de dados"

Write-Host "`n🔐 Credenciais padrão:" -ForegroundColor Green
Write-Host "   Usuário: admin"
Write-Host "   Senha: Verifique a variável ADMIN_PASSWORD no .env"

Write-Host "`n📊 Comandos úteis:" -ForegroundColor Cyan
Write-Host "   Ver logs: docker-compose logs -f"
Write-Host "   Parar: docker-compose down"
Write-Host "   Reiniciar: docker-compose restart"
Write-Host "   Backup DB: docker-compose exec quiz-app cp /app/data/quiz.db /backup/"

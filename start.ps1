# Script de inicialização do Quiz do Amor

param(
    [switch]$Dev,
    [switch]$Build,
    [switch]$Deploy,
    [switch]$Clean
)

Write-Host "💕 Quiz do Amor - Script de Inicialização" -ForegroundColor Magenta

if ($Clean) {
    Write-Host "🧹 Limpando projeto..." -ForegroundColor Yellow
    docker-compose down -v
    docker system prune -f
    Remove-Item -Path "quiz.db" -Force -ErrorAction SilentlyContinue
    Write-Host "✅ Limpeza concluída!" -ForegroundColor Green
    exit 0
}

# Verificar dependências
Write-Host "🔍 Verificando dependências..." -ForegroundColor Cyan

# Verificar Go
try {
    $goVersion = go version
    Write-Host "✅ Go encontrado: $goVersion" -ForegroundColor Green
} catch {
    Write-Host "❌ Go não encontrado. Instale o Go primeiro." -ForegroundColor Red
    exit 1
}

# Verificar Docker
try {
    docker --version | Out-Null
    Write-Host "✅ Docker encontrado" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker não encontrado. Instale o Docker Desktop primeiro." -ForegroundColor Red
    exit 1
}

# Criar arquivo .env se não existir
if (-not (Test-Path ".env")) {
    Write-Host "📝 Criando arquivo .env..." -ForegroundColor Yellow
    
    $jwtSecret = [System.Convert]::ToBase64String([System.Text.Encoding]::UTF8.GetBytes((New-Guid).ToString()))
    
    @"
# Configurações de Desenvolvimento
GIN_MODE=debug
DATABASE_PATH=./quiz.db
ADMIN_PASSWORD=admin123
JWT_SECRET=$jwtSecret
PORT=8080

# Para produção, altere estes valores
DOMAIN=localhost
"@ | Out-File -FilePath ".env" -Encoding UTF8

    Write-Host "✅ Arquivo .env criado!" -ForegroundColor Green
}

if ($Dev) {
    Write-Host "🚀 Iniciando em modo desenvolvimento..." -ForegroundColor Green
    
    # Baixar dependências
    Write-Host "📦 Baixando dependências Go..." -ForegroundColor Cyan
    go mod download
    
    # Executar migrations/setup se necessário
    if (-not (Test-Path "quiz.db")) {
        Write-Host "🗄️  Criando banco de dados..." -ForegroundColor Cyan
    }
    
    # Iniciar servidor
    Write-Host "🌟 Iniciando servidor de desenvolvimento..." -ForegroundColor Green
    Write-Host "🌐 Acesse: http://localhost:8080" -ForegroundColor Cyan
    Write-Host "👤 Login: admin / admin123" -ForegroundColor Yellow
    Write-Host "" 
    Write-Host "Pressione Ctrl+C para parar" -ForegroundColor Gray
    
    go run ./cmd/server
}

if ($Build) {
    Write-Host "🔨 Construindo aplicação..." -ForegroundColor Green
    
    # Build Go
    go build -o quiz-server.exe ./cmd/server
    
    # Build Docker
    docker build -t namorada-quiz:latest .
    
    Write-Host "✅ Build concluído!" -ForegroundColor Green
}

if ($Deploy) {
    Write-Host "🚀 Iniciando deploy..." -ForegroundColor Green
    
    # Executar script de deploy
    if (Test-Path "deploy.ps1") {
        .\deploy.ps1
    } else {
        Write-Host "❌ Script de deploy não encontrado!" -ForegroundColor Red
        exit 1
    }
}

# Se nenhum parâmetro foi passado, mostrar ajuda
if (-not ($Dev -or $Build -or $Deploy -or $Clean)) {
    Write-Host ""
    Write-Host "📖 Uso:" -ForegroundColor Yellow
    Write-Host "  .\start.ps1 -Dev      # Inicia servidor de desenvolvimento"
    Write-Host "  .\start.ps1 -Build    # Constrói a aplicação"
    Write-Host "  .\start.ps1 -Deploy   # Faz deploy com Docker"
    Write-Host "  .\start.ps1 -Clean    # Limpa projeto e dados"
    Write-Host ""
    Write-Host "💡 Exemplos:" -ForegroundColor Cyan
    Write-Host "  .\start.ps1 -Dev                 # Desenvolvimento"
    Write-Host "  .\start.ps1 -Build -Deploy       # Build + Deploy"
    Write-Host "  .\start.ps1 -Clean -Dev          # Limpar + Dev"
}

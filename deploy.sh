#!/bin/bash

# Script de deploy para produção com Cloudflare

set -e

echo "🚀 Iniciando deploy do Quiz do Amor..."

# Verificar se as variáveis estão definidas
if [ -z "$CLOUDFLARE_ZONE_ID" ] || [ -z "$CLOUDFLARE_API_TOKEN" ] || [ -z "$DOMAIN" ]; then
    echo "❌ Erro: Defina as variáveis CLOUDFLARE_ZONE_ID, CLOUDFLARE_API_TOKEN e DOMAIN"
    echo "Exemplo:"
    echo "export CLOUDFLARE_ZONE_ID=sua_zone_id"
    echo "export CLOUDFLARE_API_TOKEN=seu_token"
    echo "export DOMAIN=quiz.seudominio.com"
    exit 1
fi

# Build da imagem Docker
echo "📦 Buildando imagem Docker..."
docker build -t valentine-quiz:latest .

# Deploy usando docker-compose
echo "🚢 Fazendo deploy da aplicação..."
docker-compose down
docker-compose up -d

# Aguardar a aplicação estar pronta
echo "⏳ Aguardando aplicação estar pronta..."
sleep 10

# Verificar se está funcionando
if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    echo "✅ Aplicação rodando localmente em http://localhost:8080"
else
    echo "❌ Erro: Aplicação não está respondendo"
    docker-compose logs
    exit 1
fi

# Configurar Cloudflare (opcional)
echo "🌐 Configurando Cloudflare..."
echo "Certifique-se de que:"
echo "1. Seu domínio $DOMAIN aponta para o IP do servidor"
echo "2. O SSL/TLS está configurado como 'Full (strict)'"
echo "3. As regras de Page Rules estão configuradas se necessário"

echo "🎉 Deploy concluído!"
echo "📱 Acesse: https://$DOMAIN"
echo "🔧 Admin: https://$DOMAIN/login?admin=1"

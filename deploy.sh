#!/bin/bash

# Script de deploy para produção com Cloudflare

set -e

echo "🚀 Iniciando deploy do Quiz do Amor..."

# Verificar se as variáveis estão definidas
if [ -z "$DOMAIN" ]; then
    echo "❌ Erro: Defina a variável DOMAIN"
    echo "Exemplo:"
    echo "export DOMAIN=quiz.seudominio.com"
    echo "export EMAIL=seu@email.com  # Para Let's Encrypt"
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
if curl -f -k https://localhost:443/ > /dev/null 2>&1; then
    echo "✅ Aplicação rodando com HTTPS em https://localhost"
elif curl -f http://localhost:80/ > /dev/null 2>&1; then
    echo "⚠️  Aplicação rodando apenas com HTTP em http://localhost"
    echo "💡 Execute './setup-ssl.sh' para configurar HTTPS"
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

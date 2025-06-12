#!/bin/bash
# Script para inicializar SSL com Let's Encrypt

set -e

# Verificar se as variáveis estão definidas
if [ -z "$DOMAIN" ]; then
    echo "❌ Erro: Defina a variável DOMAIN"
    echo "Exemplo:"
    echo "export DOMAIN=quiz.seudominio.com"
    echo "export EMAIL=seu@email.com"
    exit 1
fi

if [ -z "$EMAIL" ]; then
    echo "❌ Erro: Defina a variável EMAIL para Let's Encrypt"
    exit 1
fi

echo "🔧 Configurando SSL para: $DOMAIN"

# 1. Atualizar configuração do nginx com o domínio correto
echo "📝 Atualizando configuração do nginx..."
sed -i "s/quiz\.seudominio\.com/$DOMAIN/g" nginx/sites/quiz.conf
sed -i "s/nexuscode\.tech/$DOMAIN/g" nginx/sites/quiz.conf

# 2. Criar diretórios necessários
mkdir -p certbot/www
mkdir -p certbot/conf

# 3. Subir nginx primeiro (sem SSL)
echo "🚀 Iniciando nginx para validação ACME..."
docker-compose up -d nginx

# 4. Obter certificado inicial
echo "🔐 Obtendo certificado SSL..."
docker-compose run --rm certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

# 5. Reiniciar nginx com SSL
echo "🔄 Reiniciando nginx com SSL..."
docker-compose restart nginx

# 6. Subir aplicação
echo "🎉 Subindo aplicação completa..."
docker-compose up -d

echo "✅ SSL configurado com sucesso!"
echo "🌐 Acesse: https://$DOMAIN"

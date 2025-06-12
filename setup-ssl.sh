#!/bin/bash

# Script para configurar SSL/TLS com Let's Encrypt
# Execute este script DEPOIS de configurar o DNS do seu domínio

set -e

echo "🔒 Configurando HTTPS para o Quiz do Amor..."

# Verificar se as variáveis estão definidas
if [ -z "$DOMAIN" ]; then
    echo "❌ Erro: Defina a variável DOMAIN"
    echo "Exemplo: export DOMAIN=quiz.seudominio.com"
    exit 1
fi

if [ -z "$EMAIL" ]; then
    echo "❌ Erro: Defina a variável EMAIL para o Let's Encrypt"
    echo "Exemplo: export EMAIL=seu@email.com"
    exit 1
fi

echo "📋 Configurações:"
echo "  Domínio: $DOMAIN"
echo "  Email: $EMAIL"

# Atualizar configuração do Nginx com o domínio real
echo "⚙️ Atualizando configuração do Nginx..."
sed -i "s/quiz.seudominio.com/$DOMAIN/g" nginx/sites/quiz.conf
sed -i "s/server_name _;/server_name $DOMAIN;/g" nginx/sites/quiz.conf

# Criar configuração temporária para validação do domínio
echo "🏗️ Criando configuração temporária..."
cat > nginx/sites/temp.conf << EOF
server {
    listen 80;
    server_name $DOMAIN;
    
    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
        try_files \$uri =404;
    }
    
    location / {
        return 200 "OK - Quiz do Amor está sendo configurado...";
        add_header Content-Type text/plain;
    }
}
EOF

# Backup da configuração original
cp nginx/sites/quiz.conf nginx/sites/quiz.conf.backup

# Usar configuração temporária
mv nginx/sites/quiz.conf nginx/sites/quiz.conf.ssl
mv nginx/sites/temp.conf nginx/sites/quiz.conf

# Subir apenas Nginx e Certbot
echo "🚀 Iniciando Nginx..."
docker-compose up -d nginx

# Aguardar Nginx estar pronto
echo "⏳ Aguardando Nginx estar pronto..."
sleep 10

# Gerar certificado SSL
echo "🔐 Gerando certificado SSL para $DOMAIN..."
docker-compose run --rm certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email $EMAIL \
    --agree-tos \
    --no-eff-email \
    -d $DOMAIN

# Verificar se o certificado foi criado
if [ $? -eq 0 ]; then
    echo "✅ Certificado SSL criado com sucesso!"
    
    # Restaurar configuração SSL
    mv nginx/sites/quiz.conf.ssl nginx/sites/quiz.conf
    
    # Recarregar Nginx
    echo "🔄 Recarregando Nginx com SSL..."
    docker-compose exec nginx nginx -s reload
    
    # Subir toda a stack
    echo "🚢 Iniciando aplicação completa..."
    docker-compose up -d
    
    echo "🎉 HTTPS configurado com sucesso!"
    echo "📱 Acesse: https://$DOMAIN"
    echo "🔧 Admin: https://$DOMAIN/login?admin=1"
    
else
    echo "❌ Erro ao gerar certificado SSL"
    echo "Verifique se:"
    echo "1. O domínio $DOMAIN aponta para este servidor"
    echo "2. As portas 80 e 443 estão abertas"
    echo "3. Não há outros serviços usando essas portas"
    exit 1
fi

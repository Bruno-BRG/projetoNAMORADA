#!/bin/bash

# Script para desenvolvimento local com HTTPS (certificado auto-assinado)

echo "🔒 Configurando HTTPS local para desenvolvimento..."

# Criar diretório para certificados locais
mkdir -p nginx/ssl/local

# Gerar certificado auto-assinado para localhost
if [ ! -f "nginx/ssl/local/localhost.crt" ]; then
    echo "🔐 Gerando certificado auto-assinado..."
    
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout nginx/ssl/local/localhost.key \
        -out nginx/ssl/local/localhost.crt \
        -subj "/C=BR/ST=State/L=City/O=Organization/OU=OrgUnit/CN=localhost"
        
    echo "✅ Certificado gerado!"
fi

# Criar configuração local do Nginx
cat > nginx/sites/local.conf << 'EOF'
server {
    listen 80;
    server_name localhost;
    return 301 https://$host$request_uri;
}

server {
    listen 443 ssl http2;
    server_name localhost;
    
    # Certificado auto-assinado para desenvolvimento
    ssl_certificate /etc/nginx/ssl/local/localhost.crt;
    ssl_certificate_key /etc/nginx/ssl/local/localhost.key;
    
    # Configurações SSL básicas
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;
    
    # Headers de segurança
    add_header X-Frame-Options "DENY" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    location / {
        proxy_pass http://valentine-quiz:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $server_name;
    }
}
EOF

# Backup e substituir configuração
if [ -f "nginx/sites/quiz.conf" ]; then
    cp nginx/sites/quiz.conf nginx/sites/quiz.conf.production
fi
cp nginx/sites/local.conf nginx/sites/quiz.conf

# Criar docker-compose para desenvolvimento
cat > docker-compose.dev.yml << 'EOF'
version: '3.8'

services:
  valentine-quiz:
    build: .
    expose:
      - "8080"
    environment:
      - GIN_MODE=debug
      - PORT=8080
      - DB_PATH=/data/quiz.db
      - ADMIN_USERNAME=admin
      - ADMIN_PASSWORD=admin123
      - VISITOR_USERNAME=momo
      - VISITOR_PASSWORD=momo3006
      - JWT_SECRET=dev_secret_not_for_production
    volumes:
      - quiz_data:/data
    networks:
      - quiz_network

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/sites/:/etc/nginx/conf.d/:ro
      - ./nginx/ssl/local/:/etc/nginx/ssl/local/:ro
    depends_on:
      - valentine-quiz
    networks:
      - quiz_network

volumes:
  quiz_data:
    driver: local

networks:
  quiz_network:
    driver: bridge
EOF

echo "🚀 Iniciando ambiente de desenvolvimento com HTTPS..."
docker-compose -f docker-compose.dev.yml up -d

echo "⏳ Aguardando aplicação estar pronta..."
sleep 15

# Verificar se está funcionando
if curl -f -k https://localhost/ > /dev/null 2>&1; then
    echo "✅ Desenvolvimento rodando com HTTPS!"
    echo "📱 Acesse: https://localhost (aceite o certificado auto-assinado)"
    echo "🔧 Admin: https://localhost/login?admin=1"
    echo ""
    echo "⚠️  ATENÇÃO: Certificado auto-assinado para desenvolvimento!"
    echo "   O navegador mostrará aviso de segurança - é normal."
else
    echo "❌ Erro: Aplicação não está respondendo"
    docker-compose -f docker-compose.dev.yml logs
fi

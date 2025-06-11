#!/bin/bash

# Script de deploy para o Quiz do Amor
# Para usar com Cloudflare e Docker

set -e

echo "🚀 Iniciando deploy do Quiz do Amor..."

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar se Docker está rodando
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker não está rodando. Inicie o Docker e tente novamente.${NC}"
    exit 1
fi

# Verificar se .env existe
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}⚠️  Arquivo .env não encontrado. Criando arquivo de exemplo...${NC}"
    cat > .env << EOF
# Configurações de Produção
GIN_MODE=release
DATABASE_PATH=/app/data/quiz.db
ADMIN_PASSWORD=SuaSenhaSeguraAqui123!
JWT_SECRET=$(openssl rand -base64 32)
PORT=8080

# Configurações do Cloudflare (opcional)
CLOUDFLARE_API_TOKEN=seu_token_aqui
CLOUDFLARE_ZONE_ID=seu_zone_id_aqui
DOMAIN=seudominio.com
EOF
    echo -e "${GREEN}✅ Arquivo .env criado. Configure as variáveis antes de continuar.${NC}"
    exit 1
fi

# Carregar variáveis do .env
source .env

echo -e "${GREEN}📦 Construindo imagem Docker...${NC}"
docker build -t namorada-quiz:latest .

echo -e "${GREEN}🛑 Parando containers existentes...${NC}"
docker-compose down || true

echo -e "${GREEN}🚀 Subindo nova versão...${NC}"
docker-compose up -d

echo -e "${GREEN}⏳ Aguardando aplicação inicializar...${NC}"
sleep 10

# Health check
echo -e "${GREEN}🔍 Verificando saúde da aplicação...${NC}"
if curl -f http://localhost:8080/ > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Aplicação está rodando corretamente!${NC}"
else
    echo -e "${RED}❌ Aplicação não respondeu ao health check${NC}"
    echo -e "${YELLOW}📋 Logs do container:${NC}"
    docker-compose logs quiz-app
    exit 1
fi

echo -e "${GREEN}🎉 Deploy concluído com sucesso!${NC}"
echo -e "${GREEN}🌐 Aplicação disponível em: http://localhost:8080${NC}"

if [ ! -z "$DOMAIN" ] && [ "$DOMAIN" != "seudominio.com" ]; then
    echo -e "${GREEN}🌍 Domínio configurado: https://$DOMAIN${NC}"
fi

echo -e "${YELLOW}📝 Próximos passos:${NC}"
echo -e "1. Configure seu proxy reverso (Nginx/Cloudflare Tunnel)"
echo -e "2. Configure SSL/TLS"
echo -e "3. Teste a aplicação thoroughly"
echo -e "4. Configure backup do banco de dados"

echo -e "${GREEN}🔐 Credenciais padrão:${NC}"
echo -e "   Usuário: admin"
echo -e "   Senha: Verifique a variável ADMIN_PASSWORD no .env"

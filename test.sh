#!/bin/bash

# Script de teste para verificar se o sistema está funcionando

echo "🧪 Testando o Quiz do Dia dos Namorados..."

# Cores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Função para testar endpoints
test_endpoint() {
    local url=$1
    local expected_code=$2
    local description=$3
    
    response=$(curl -s -o /dev/null -w "%{http_code}" "$url")
    
    if [ "$response" -eq "$expected_code" ]; then
        echo -e "${GREEN}✅ $description${NC}"
    else
        echo -e "${RED}❌ $description (Expected $expected_code, got $response)${NC}"
    fi
}

# Testar se o servidor está rodando
echo -e "${YELLOW}📡 Testando conectividade...${NC}"
test_endpoint "http://localhost:8080" 200 "Homepage carregando"
test_endpoint "http://localhost:8080/login" 200 "Página de login acessível"
test_endpoint "http://localhost:8080/login?admin=1" 200 "Login admin acessível"

# Testar rotas protegidas (devem redirecionar)
echo -e "${YELLOW}🔒 Testando autenticação...${NC}"
test_endpoint "http://localhost:8080/quiz/" 302 "Quiz protegido (redireciona)"
test_endpoint "http://localhost:8080/admin/" 302 "Admin protegido (redireciona)"

# Verificar banco de dados
echo -e "${YELLOW}💾 Verificando banco de dados...${NC}"
if [ -f "quiz.db" ]; then
    question_count=$(sqlite3 quiz.db "SELECT COUNT(*) FROM questions;")
    echo -e "${GREEN}✅ Banco existe com $question_count perguntas${NC}"
else
    echo -e "${RED}❌ Banco de dados não encontrado${NC}"
fi

# Verificar arquivos importantes
echo -e "${YELLOW}📁 Verificando arquivos...${NC}"
files=("cmd/server/main.go" "internal/handlers/handlers.go" "web/templates/home.html" ".env" "Dockerfile")

for file in "${files[@]}"; do
    if [ -f "$file" ]; then
        echo -e "${GREEN}✅ $file existe${NC}"
    else
        echo -e "${RED}❌ $file não encontrado${NC}"
    fi
done

# Testar compilação
echo -e "${YELLOW}🔧 Testando compilação...${NC}"
if go build -o test-quiz ./cmd/server/main.go 2>/dev/null; then
    echo -e "${GREEN}✅ Projeto compila sem erros${NC}"
    rm -f test-quiz
else
    echo -e "${RED}❌ Erro de compilação${NC}"
fi

echo ""
echo -e "${YELLOW}📋 Resumo dos testes:${NC}"
echo "🏠 Homepage: http://localhost:8080"
echo "👨‍💼 Admin: http://localhost:8080/login?admin=1"
echo "💕 Visitante: http://localhost:8080/login"
echo ""
echo -e "${GREEN}🎉 Testes concluídos!${NC}"

# Verificar se perguntas estão no horário correto
echo -e "${YELLOW}⏰ Verificando horários das perguntas...${NC}"
if [ -f "quiz.db" ]; then
    sqlite3 quiz.db "SELECT title, datetime(scheduled_at, 'localtime') as horario FROM questions ORDER BY scheduled_at;" | while read line; do
        echo "📝 $line"
    done
fi

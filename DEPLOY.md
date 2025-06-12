# 🚀 Guia de Deploy - Quiz do Dia dos Namorados

## Resumo do que fizemos

✅ **Sistema completo funcionando**:
- Backend Go com Gin Framework
- Frontend com HTMX para interatividade
- Sistema de quiz programado por horário
- Painel administrativo
- Deploy ready com Docker

## 🏠 Desenvolvimento Local

O projeto já está funcionando! Para testar:

```bash
# 1. Execute o servidor (já está rodando)
go run cmd/server/main.go

# 2. Acesse as páginas:
# - Homepage: http://localhost:8080
# - Login visitante: http://localhost:8080/login
# - Login admin: http://localhost:8080/login?admin=1
```

**Credenciais de teste:**
- **Admin**: `admin` / `admin123`
- **Visitante**: `momo` / `momo3006`

## 🌐 Deploy em Produção

### Opção 1: VPS/Servidor próprio

1. **Upload dos arquivos**
```bash
# Copie todo o projeto para o servidor
scp -r . usuario@seu-servidor.com:/app/quiz
```

2. **Configure as variáveis de ambiente**
```bash
# No servidor, edite o .env
nano .env

# Mude as senhas para algo seguro:
ADMIN_PASSWORD=sua_senha_super_segura
VISITOR_PASSWORD=momo3006
JWT_SECRET=chave_jwt_bem_complexa_aqui
GIN_MODE=release
```

3. **Execute com Docker**
```bash
# No servidor
docker-compose up -d

# Verifique se está funcionando
curl http://localhost:8080
```

### Opção 2: Cloud Provider (Render, Railway, Heroku)

1. **Crie conta no Render.com** (recomendado - fácil e grátis)

2. **Conecte seu repositório GitHub**

3. **Configure o deploy:**
   - Build Command: `go build -o main cmd/server/main.go`
   - Start Command: `./main`
   - Environment Variables:
     ```     GIN_MODE=release
     ADMIN_PASSWORD=sua_senha
     VISITOR_PASSWORD=momo3006
     JWT_SECRET=chave_complexa
     ```

## ⚙️ Configuração do Cloudflare

1. **Adicione seu domínio no Cloudflare**

2. **Configure DNS:**
   ```
   Type: A
   Name: quiz (ou @)
   Value: IP_DO_SEU_SERVIDOR
   Proxy: ON (nuvem laranja)
   ```

3. **Configurações SSL/TLS:**
   - Vá em SSL/TLS > Overview
   - Escolha "Full (strict)"

4. **Page Rules (opcional):**
   ```
   URL: quiz.seudominio.com/*
   Settings: Always Use HTTPS
   ```

## 📱 Como usar no Dia dos Namorados

### Preparação (você):

1. **Acesse o painel admin**: `https://quiz.seudominio.com/login?admin=1`

2. **Crie as perguntas**: Vá em "Gerenciar Perguntas" > "Nova Pergunta"
   - Defina horários estratégicos (café, almoço, tarde, jantar)
   - Seja criativo nas recompensas!

3. **Teste tudo**: Use uma aba anônima para testar como visitante

### No dia especial (ela):

1. **Envie o link**: `https://quiz.seudominio.com`

2. **Ela faz login** com as credenciais que você definiu

3. **Quiz aparecem automaticamente** nos horários programados

4. **HTMX atualiza em tempo real** - ela vê countdown até próximo quiz

## 🔧 Comandos Úteis

```bash
# Ver logs em produção
docker-compose logs -f

# Backup do banco
cp quiz.db quiz_backup_$(date +%Y%m%d).db

# Restart da aplicação
docker-compose restart

# Verificar se está rodando
curl -I https://quiz.seudominio.com

# Adicionar mais perguntas via script
go run cmd/seed/main.go
```

## 🚨 Checklist Final

- [ ] Senhas alteradas no .env
- [ ] Domínio configurado no Cloudflare
- [ ] SSL funcionando (cadeado verde)
- [ ] Testado login admin e visitante
- [ ] Perguntas criadas com horários corretos
- [ ] Backup do banco feito

## 💡 Dicas Pro

1. **Horários estratégicos:**
   - 07:00 - "Bom dia, meu amor"
   - 12:00 - "Quiz do almoço"
   - 18:00 - "Chegando em casa"
   - 21:00 - "Quiz da noite"

2. **Recompensas criativas:**
   - Físicas: beijos, abraços, massagens
   - Experiências: jantar fora, cinema, viagem
   - Gestos: café na cama, flores, carta

3. **Perguntas interessantes:**
   - Memórias do relacionamento
   - Preferências dela que você observou
   - Planos futuros juntos
   - Coisas específicas que só vocês sabem

## 🆘 Troubleshooting

**Site não carrega:**
```bash
# Verifique se está rodando
docker ps
# Se não estiver, restart
docker-compose up -d
```

**Erro 502 Bad Gateway:**
- Verifique se a porta 8080 está aberta
- Confirme que a aplicação está rodando internamente

**Login não funciona:**
- Verifique as variáveis de ambiente
- Confirme que o .env foi carregado

---

**💕 Agora é só dar deploy e impressionar sua namorada!**

*P.S.: Se ela é dev, prepare-se para code review... 😄*

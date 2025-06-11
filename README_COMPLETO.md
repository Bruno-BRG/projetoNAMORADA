# 💕 Quiz do Amor - Dia dos Namorados

Um sistema de quiz interativo desenvolvido especialmente para o Dia dos Namorados, onde perguntas são liberadas ao longo do dia em horários específicos, testando o conhecimento sobre o relacionamento.

## 🎯 Características

- **Perguntas Programadas**: Sistema de agendamento para liberar perguntas em horários específicos
- **Interface Moderna**: Design responsivo com HTMX para interações dinâmicas
- **Sistema de Recompensas**: Recompensas personalizadas para respostas corretas
- **Painel Admin**: Interface completa para gerenciar perguntas e monitorar respostas
- **Deployment Fácil**: Docker + Cloudflare ready
- **Segurança**: JWT authentication e middleware de segurança

## 🚀 Tecnologias

- **Backend**: Go (Gin framework)
- **Frontend**: HTML5 + HTMX + CSS3
- **Database**: SQLite (fácil deployment)
- **Containerização**: Docker + Docker Compose
- **Proxy/CDN**: Cloudflare (opcional)

## 📁 Estrutura do Projeto

```
namorada-quiz/
├── cmd/server/          # Aplicação principal
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Middleware (auth, cors, etc)
│   ├── models/          # Modelos de dados
│   └── database/        # Operações de banco
├── templates/           # Templates HTML
│   └── partials/        # Componentes HTMX
├── static/             # Assets estáticos
├── docker/             # Configurações Docker
├── Dockerfile          # Build da aplicação
├── docker-compose.yml  # Orquestração
├── deploy.ps1          # Script de deploy (Windows)
├── deploy.sh           # Script de deploy (Linux/Mac)
└── start.ps1           # Script de desenvolvimento
```

## 🛠️ Configuração e Instalação

### Pré-requisitos

- Go 1.21+
- Docker Desktop
- PowerShell (Windows) ou Bash (Linux/Mac)

### Instalação Rápida

1. **Clone o repositório**
```bash
git clone <seu-repo>
cd projetoNAMORADA
```

2. **Iniciar desenvolvimento** (Windows)
```powershell
.\start.ps1 -Dev
```

3. **Acessar aplicação**
- URL: http://localhost:8080
- Admin: `admin` / `admin123`

### Configuração Manual

1. **Instalar dependências**
```bash
go mod download
```

2. **Configurar ambiente**
```bash
cp .env.example .env
# Edite o arquivo .env com suas configurações
```

3. **Executar aplicação**
```bash
go run ./cmd/server
```

## 🐳 Deploy com Docker

### Desenvolvimento
```powershell
.\start.ps1 -Build -Deploy
```

### Produção
```powershell
# Configure o .env para produção
.\deploy.ps1
```

Ou manualmente:
```bash
docker-compose up -d
```

## ⚙️ Configuração

### Variáveis de Ambiente (.env)

```env
# Modo da aplicação
GIN_MODE=release                    # debug | release

# Banco de dados
DATABASE_PATH=/app/data/quiz.db

# Autenticação
ADMIN_PASSWORD=SuaSenhaSegura123!
JWT_SECRET=seu-jwt-secret-super-seguro

# Servidor
PORT=8080

# Cloudflare (opcional)
CLOUDFLARE_API_TOKEN=seu_token
CLOUDFLARE_ZONE_ID=seu_zone_id
DOMAIN=seudominio.com
```

### Criando Perguntas

Acesse o painel admin em `/admin` e use a interface para criar perguntas:

```json
{
  "title": "Nossa primeira viagem juntos",
  "description": "Onde fomos em nossa primeira viagem romântica?",
  "options": ["Paris", "Roma", "Praia do Rosa", "Gramado"],
  "correct_answer": "Gramado",
  "reward": "💝 Uma massagem relaxante!",
  "scheduled_at": "2025-02-14T14:00:00Z"
}
```

## 📱 Uso da Aplicação

### Para Administradores

1. **Login**: Use credenciais de admin
2. **Criar Perguntas**: Defina título, opções, resposta correta, recompensa e horário
3. **Monitorar**: Veja estatísticas de uso e respostas
4. **Gerenciar**: Edite ou desative perguntas conforme necessário

### Para Visitantes (Namorada)

1. **Login**: Use credenciais fornecidas pelo admin
2. **Dashboard**: Veja perguntas disponíveis e estatísticas
3. **Responder**: Clique em perguntas liberadas para responder
4. **Recompensas**: Receba recompensas por respostas corretas

## 🔧 Scripts Disponíveis

### Windows PowerShell

```powershell
# Desenvolvimento
.\start.ps1 -Dev

# Build
.\start.ps1 -Build

# Deploy
.\start.ps1 -Deploy

# Limpeza
.\start.ps1 -Clean
```

### Comandos Docker

```bash
# Ver logs
docker-compose logs -f

# Reiniciar
docker-compose restart

# Parar
docker-compose down

# Backup do banco
docker-compose exec quiz-app cp /app/data/quiz.db /backup/
```

## 🌐 Deploy em Produção

### Com Cloudflare Tunnel

1. **Configure o tunnel**
```bash
cloudflared tunnel create quiz-amor
cloudflared tunnel route dns quiz-amor seudominio.com
```

2. **Configure o docker-compose.prod.yml**
```yaml
version: '3.8'
services:
  quiz-app:
    build: .
    environment:
      - GIN_MODE=release
    # ... outras configurações
```

3. **Deploy**
```bash
.\deploy.ps1
```

### Com Nginx (Opcional)

```nginx
server {
    listen 443 ssl;
    server_name seudominio.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## 📊 Monitoramento

### Health Check
```bash
curl http://localhost:8080/health
```

### Logs
```bash
# Logs da aplicação
docker-compose logs quiz-app

# Logs em tempo real
docker-compose logs -f
```

### Métricas
- Dashboard admin mostra estatísticas em tempo real
- API endpoint: `/api/admin/stats`

## 🔒 Segurança

- JWT tokens com expiração
- Passwords hasheados com bcrypt
- Headers de segurança
- CORS configurado
- Sanitização de inputs
- Rate limiting (implementar se necessário)

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/nova-feature`)
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`)
4. Push para a branch (`git push origin feature/nova-feature`)
5. Abra um Pull Request

## 📝 TODO

- [ ] Rate limiting
- [ ] Email notifications
- [ ] Backup automático
- [ ] Métricas avançadas
- [ ] Temas personalizáveis
- [ ] Internacionalização
- [ ] PWA support
- [ ] WebSocket para updates em tempo real

## 🐛 Troubleshooting

### Problemas Comuns

1. **Erro de permissão de banco**
```bash
chmod 666 quiz.db
```

2. **Docker não inicia**
```bash
docker system prune -f
docker-compose up --build
```

3. **JWT inválido**
- Verifique se o JWT_SECRET está correto
- Limpe cookies do browser

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes.

## 💝 Créditos

Desenvolvido com amor para o Dia dos Namorados 💕

---

**Nota**: Este README assume que você está usando este projeto para fins pessoais/românticos. Lembre-se de personalizar todas as configurações e senhas antes do deploy em produção!

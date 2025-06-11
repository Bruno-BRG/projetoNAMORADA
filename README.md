# Projeto Quiz Namorada - Dia dos Namorados 💕

Um site interativo com quiz temporizado para o Dia dos Namorados, construído com Go, HTMX e muito amor (e um pouco de sarcasmo técnico).

## Funcionalidades

- 🕐 **Quiz com horários programados**: Perguntas liberadas ao longo do dia
- 👩‍💼 **Painel Admin**: Criar perguntas, definir horários e recompensas
- 💖 **Interface para Namorada**: Responder quiz e receber recompensas
- 🔒 **Autenticação simples**: Acesso seguro para ambas as partes
- 🐳 **Docker**: Deploy fácil em qualquer cloud
- ☁️ **Cloudflare Ready**: Configurado para proxy reverso

## Estrutura do Projeto

```
.
├── cmd/server/           # Entry point da aplicação
├── internal/
│   ├── handlers/         # HTTP handlers
│   ├── models/          # Estruturas de dados
│   ├── database/        # Camada de dados
│   └── middleware/      # Middlewares customizados
├── templates/           # Templates HTML
├── static/             # CSS, JS, imagens
├── docker/             # Dockerfiles e configs
└── migrations/         # Scripts de banco
```

## Como Rodar

```bash
# Desenvolvimento
go run cmd/server/main.go

# Com Docker
docker-compose up --build

# Deploy (com suas configurações de domínio)
docker build -t namorada-quiz .
```

## Stack Técnica

- **Backend**: Go + Gin
- **Frontend**: HTMX + Vanilla CSS (sem frameworks desnecessários)
- **Banco**: SQLite (simples e eficiente para esse caso)
- **Deploy**: Docker + Cloudflare
- **Autenticação**: JWT com cookies httpOnly

## TODO

- [ ] Implementar timezone handling
- [ ] Sistema de notificações (opcional)
- [ ] Backup automático das respostas
- [ ] Logs detalhados para debug

---

*"O amor é como código: funciona melhor quando é simples, elegante e bem testado."* 😉

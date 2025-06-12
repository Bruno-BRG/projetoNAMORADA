# 💕 Quiz do Dia dos Namorados

Um site interativo de quiz romântico com horários programados, construído com **Go**, **HTMX** e muito amor! 

## ✨ Características

- 🕐 **Quiz com horários programados** - Perguntas são liberadas em horários específicos ao longo do dia
- ❤️ **Interface romântica** - Design bonito e responsivo com animações suaves
- ⚡ **HTMX para interatividade** - Atualizações em tempo real sem refresh da página
- 🎁 **Sistema de recompensas** - Cada resposta certa tem uma recompensa especial
- 👨‍💼 **Painel administrativo** - Para criar perguntas e gerenciar o quiz
- 🐳 **Deploy com Docker** - Pronto para produção com Cloudflare
- 🔒 **Autenticação simples** - Sistema de login para admin e visitante

## 🛠️ Tecnologias

- **Backend**: Go + Gin Framework
- **Frontend**: HTMX + TailwindCSS
- **Banco de dados**: SQLite
- **Deploy**: Docker + Docker Compose
- **Proxy**: Cloudflare (opcional)

- **Backend**: Go com Gin framework
- **Frontend**: HTMX + TailwindCSS (porque CSS vanilla é coisa do passado)
- **Banco**: SQLite (simples, mas confiável)
- **Deploy**: Docker + Cloudflare
- **Autenticação**: Session-based (simples e seguro)

## Funcionalidades

### Para o Admin (você)
- Login seguro
- Criar/editar perguntas
- Definir horários de liberação
- Visualizar respostas
- Dashboard com progresso

### Para a Visitante (sua namorada)
- Login simples
- Quiz liberado por horário
- Recompensas por acertos
- Interface responsiva e bonita

## Estrutura do Projeto

```
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── models/
│   ├── middleware/
│   └── database/
├── web/
│   ├── templates/
│   ├── static/
│   └── assets/
├── docker/
├── migrations/
└── configs/
```

## Getting Started

```bash
# Desenvolvimento
go mod tidy
go run cmd/server/main.go

# Produção
docker-compose up -d
```

## Variáveis de Ambiente

```
DB_PATH=./quiz.db
ADMIN_PASSWORD=sua_senha_aqui
SECRET_KEY=sua_chave_secreta
PORT=8080
```

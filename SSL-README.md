# 🔒 Configuração HTTPS/SSL

## Opções Disponíveis

### 🏠 Desenvolvimento Local
Para testar HTTPS em desenvolvimento:
```bash
./dev-https.sh
```
- Gera certificado auto-assinado
- Roda em https://localhost
- Aceite o aviso de segurança no navegador

### 🌐 Produção com Let's Encrypt
Para deploy em produção com SSL real:

1. **Configure as variáveis:**
```bash
export DOMAIN=quiz.seudominio.com
export EMAIL=seu@email.com
```

2. **Configure DNS:**
   - Aponte seu domínio para o IP do servidor
   - Aguarde propagação (pode levar até 24h)

3. **Configure SSL:**
```bash
./setup-ssl.sh
```

4. **Deploy completo:**
```bash
./deploy.sh
```

## 📋 Checklist SSL

- [ ] Domínio apontando para o servidor
- [ ] Portas 80 e 443 abertas no firewall
- [ ] Variáveis DOMAIN e EMAIL definidas
- [ ] DNS propagado (teste: `dig seu-dominio.com`)

## 🔧 Comandos Úteis

```bash
# Verificar status dos certificados
docker-compose exec certbot certbot certificates

# Renovar certificados manualmente
docker-compose exec certbot certbot renew

# Ver logs do Nginx
docker-compose logs nginx

# Verificar configuração do Nginx
docker-compose exec nginx nginx -t

# Recarregar Nginx
docker-compose exec nginx nginx -s reload
```

## 🚨 Troubleshooting

### "Certificado não confiável"
- **Local**: Normal com certificado auto-assinado
- **Produção**: Verifique se usou Let's Encrypt

### "Erro 502 Bad Gateway"
- Verifique se a aplicação Go está rodando
- Confira logs: `docker-compose logs`

### "Falha ao gerar certificado"
- Confirme que o domínio aponta para o servidor
- Verifique se as portas 80/443 estão abertas
- Teste: `curl -I http://seu-dominio.com`

### "ERR_SSL_PROTOCOL_ERROR"
- Verifique se o certificado está no lugar certo
- Confira configuração do Nginx: `nginx -t`

## 🔄 Renovação Automática

O Certbot renova automaticamente os certificados a cada 12h.
Para verificar:
```bash
docker-compose logs certbot
```

## 🛡️ Segurança

A configuração inclui:
- ✅ TLS 1.2/1.3 apenas
- ✅ HSTS habilitado
- ✅ Headers de segurança
- ✅ Rate limiting
- ✅ CSP (Content Security Policy)
- ✅ Redirect HTTP → HTTPS

## 📱 URLs Finais

- **Site**: `https://seu-dominio.com`
- **Admin**: `https://seu-dominio.com/login?admin=1`
- **Health**: `https://seu-dominio.com/health`

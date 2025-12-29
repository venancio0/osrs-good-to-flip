# Deploy Guide

## Frontend no Vercel

### Pré-requisitos
- Conta no Vercel (https://vercel.com)
- Projeto conectado ao GitHub

### Configuração via Dashboard

1. **Conectar repositório no Vercel:**
   - Acesse https://vercel.com
   - Clique em "Add New Project"
   - Importe o repositório `venancio0/osrs-good-to-flip`
   - Configure o projeto:
     - **Framework Preset:** Next.js (detectado automaticamente)
     - **Root Directory:** `frontend` ⚠️ **IMPORTANTE**
     - **Build Command:** `npm run build` (automático)
     - **Output Directory:** `.next` (automático)
     - **Install Command:** `npm install` (automático)

2. **Variáveis de Ambiente:**
   No painel do Vercel, vá em Settings → Environment Variables e adicione:
   ```
   NEXT_PUBLIC_API_URL = https://osrs-good-to-flip.onrender.com
   ```
   - Marque para Production, Preview e Development

3. **Deploy:**
   - O Vercel fará deploy automaticamente após o push
   - Ou clique em "Deploy" no painel
   - Aguarde o build completar

### Configuração Alternativa (via CLI)

```bash
# Instalar Vercel CLI
npm i -g vercel

# Login
vercel login

# Deploy
cd frontend
vercel

# Configurar variáveis de ambiente
vercel env add NEXT_PUBLIC_API_URL
# Digite: https://osrs-good-to-flip.onrender.com
```

## Backend no Render

O backend já está configurado para rodar no Render em:
- URL: https://osrs-good-to-flip.onrender.com
- Porta: Configurada via variável de ambiente `PORT`

### Variáveis de Ambiente no Render

Configure no painel do Render:
- `PORT` (opcional, padrão: 8080)
- `OSRS_WIKI_USER_AGENT` (obrigatório)
- `OSRS_WIKI_BASE_URL` (opcional)
- `OSRS_WIKI_TIMEOUT_MS` (opcional)
- `OSRS_WIKI_CACHE_TTL_SEC` (opcional)
- `PRICE_UPDATE_INTERVAL_MIN` (opcional, padrão: 5)

## CORS

O backend já está configurado para aceitar requisições do Vercel. Se necessário, adicione o domínio do Vercel no arquivo `routes.go`:

```go
AllowedOrigins: []string{
    "http://localhost:3000",
    "http://localhost:3001",
    "https://seu-projeto.vercel.app", // Adicione aqui
},
```

Ou use variável de ambiente para configurar dinamicamente.


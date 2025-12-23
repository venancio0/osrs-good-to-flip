# OSRS Good to Flip

MVP para acompanhar preços do Grand Exchange do Old School RuneScape.

## Estrutura do Projeto

Monorepo com backend Go e frontend Next.js:

```
osrs-good-to-flip/
├── backend/          # API Go com Clean Architecture
├── frontend/         # Next.js + Tailwind CSS
└── README.md
```

## Arquitetura

O projeto segue Clean Architecture com as seguintes camadas:

- **Domain**: Entidades e interfaces (Ports)
- **Application**: Use cases (lógica de negócio)
- **Infrastructure**: Implementações concretas (OSRS client, cache, repository)
- **Interfaces**: HTTP handlers e rotas

## Pré-requisitos

- Go 1.21 ou superior
- Node.js 18+ e npm/yarn
- Git

## Setup

### Backend

1. Navegue até a pasta do backend:
```bash
cd backend
```

2. Instale as dependências:
```bash
go mod download
```

3. Execute o servidor:
```bash
go run cmd/api/main.go
```

O servidor estará rodando em `http://localhost:8080`

### Frontend

1. Navegue até a pasta do frontend:
```bash
cd frontend
```

2. Instale as dependências:
```bash
npm install
# ou
yarn install
```

3. Configure a URL da API (opcional):
Crie um arquivo `.env.local` com:
```
NEXT_PUBLIC_API_URL=http://localhost:8080
```

4. Execute o servidor de desenvolvimento:
```bash
npm run dev
# ou
yarn dev
```

O frontend estará rodando em `http://localhost:3000`

## Endpoints da API

### GET /health
Health check da API.

**Resposta:**
```json
{
  "status": "ok",
  "service": "osrs-good-to-flip"
}
```

### GET /items
Lista todos os itens ou busca por nome.

**Query Parameters:**
- `q` (opcional): Termo de busca

**Resposta:**
```json
[
  {
    "item_id": 1,
    "name": "Rune Scimitar",
    "price": 15000,
    "avg_24h": 14800,
    "avg_7d": 15000,
    "trend": "UP",
    "updated_at": "2024-01-01T00:00:00Z"
  }
]
```

### GET /items/{id}
Retorna detalhes de um item específico.

**Resposta:**
```json
{
  "item_id": 1,
  "name": "Rune Scimitar",
  "price": 15000,
  "avg_24h": 14800,
  "avg_7d": 15000,
  "trend": "UP",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Funcionalidades do MVP

- ✅ Lista de itens do Grand Exchange
- ✅ Busca por nome de item
- ✅ Detalhes do item (preço atual, médias 24h/7d, tendência)
- ✅ Interface navegável e responsiva
- ✅ Dados mockados para desenvolvimento

## Tecnologias

### Backend
- Go 1.21+
- Chi Router
- Clean Architecture

### Frontend
- Next.js 14
- React 18
- TypeScript
- Tailwind CSS

## Próximos Passos

- [ ] Integração com API real do OSRS Wiki
- [ ] Cálculo de médias baseado em histórico
- [ ] Persistência com SQLite
- [ ] Testes unitários e de integração
- [ ] Histórico de preços com gráficos
- [ ] Alertas de preço
- [ ] Top gainers/losers

## Desenvolvimento

### Executar ambos os serviços

Em terminais separados:

**Terminal 1 - Backend:**
```bash
cd backend
go run cmd/api/main.go
```

**Terminal 2 - Frontend:**
```bash
cd frontend
npm run dev
```

Acesse `http://localhost:3000` no navegador.

## Licença

Este é um projeto MVP para fins educacionais.


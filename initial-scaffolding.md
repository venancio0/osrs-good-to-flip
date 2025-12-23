# 🎯 Visão do Produto (MVP)

### Nome provisório

**OSRS Good to Flip**

### Problema que resolve

Jogadores precisam acompanhar preços do **Grand Exchange** sem ficar abrindo wiki ou planilhas.

### Objetivo do MVP

* Exibir preços atualizados de itens populares
* Mostrar tendência simples (↑ ↓ →)
* Zero login
* Performance boa (cache/persistência leve)

---

# 📦 Funcionalidades do MVP (escopo fechado)

### Backend

* Buscar preços atuais do OSRS Wiki
* Expor API própria simples
* Cachear resultados
* Persistência opcional (para histórico mínimo)

### Frontend

* Lista de itens
* Busca por nome
* Detalhe do item:

  * Preço atual
  * Média 24h / 7d
  * Tendência

---

# 🌐 API Externa (fonte de dados)

**OSRS Wiki – Grand Exchange**

Endpoints principais:

```
GET /api/v1/osrs/latest
GET /api/v1/osrs/5m
GET /api/v1/osrs/1h
```

Dados:

* item_id
* high / low
* timestamp

➡️ MVP pode começar **apenas com `/latest`**

---

# 🧱 Arquitetura (Clean Architecture)

### Camadas

```
┌──────────────────────────┐
│        Frontend          │
└───────────▲──────────────┘
            │ REST
┌───────────┴──────────────┐
│     Interface Adapters   │  (HTTP handlers)
└───────────▲──────────────┘
            │
┌───────────┴──────────────┐
│     Application Layer    │  (Use Cases)
└───────────▲──────────────┘
            │
┌───────────┴──────────────┐
│        Domain             │  (Entities, Interfaces)
└───────────▲──────────────┘
            │
┌───────────┴──────────────┐
│     Infrastructure       │  (API externa, cache, DB)
└──────────────────────────┘
```

---

# 🧠 Domain (núcleo)

### Entidade principal

```go
type ItemPrice struct {
    ItemID    int
    Name      string
    Price     int
    Avg24h    int
    Avg7d     int
    Trend     TrendType
    UpdatedAt time.Time
}
```

```go
type TrendType string

const (
    TrendUp   TrendType = "UP"
    TrendDown TrendType = "DOWN"
    TrendFlat TrendType = "FLAT"
)
```

### Interfaces (Ports)

```go
type PriceProvider interface {
    FetchLatestPrices(ctx context.Context) (map[int]int, error)
}

type ItemRepository interface {
    SavePrices(ctx context.Context, prices []ItemPrice) error
    GetItemByID(ctx context.Context, id int) (ItemPrice, error)
    SearchItems(ctx context.Context, query string) ([]ItemPrice, error)
}
```

➡️ **SOLID**:

* Dependency Inversion
* Domain não conhece HTTP, DB nem OSRS Wiki

---

# ⚙️ Application Layer (Use Cases)

### Casos de uso

* `UpdatePricesUseCase`
* `GetItemUseCase`
* `SearchItemsUseCase`

Exemplo:

```go
type UpdatePricesUseCase struct {
    provider PriceProvider
    repo     ItemRepository
}

func (uc *UpdatePricesUseCase) Execute(ctx context.Context) error {
    prices, err := uc.provider.FetchLatestPrices(ctx)
    if err != nil {
        return err
    }

    // map -> entity
    // calculate trend (simples)
    return uc.repo.SavePrices(ctx, items)
}
```

---

# 🌍 Infrastructure

### OSRS Wiki Client

```go
type OsrsWikiClient struct {
    http *http.Client
}

func (c *OsrsWikiClient) FetchLatestPrices(ctx context.Context) (map[int]int, error) {
    // chama API externa
}
```

### Cache (MVP)

Opções:

* In-memory (map + mutex)
* Redis (se quiser já deixar prod-like)

Recomendação MVP:

> **In-memory + TTL**

---

### Persistência (opcional no MVP)

Opções:

* SQLite
* PostgreSQL
* Ou nenhuma (apenas cache)

➡️ MVP **funciona sem DB**, só com cache

---

# 🌐 Interface Adapters (HTTP)

### Endpoints REST

```
GET /items
GET /items/{id}
GET /health
```

Exemplo handler:

```go
func GetItemHandler(uc *GetItemUseCase) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        id := chi.URLParam(r, "id")
        item, err := uc.Execute(r.Context(), id)
        // response
    }
}
```

Framework sugerido:

* `chi` (leve, idiomático)

---

# 🎨 Frontend (Next.js + Tailwind)

### Páginas

```
/            -> lista de itens
/items/[id] -> detalhe
```

### Componentes

* `ItemCard`
* `PriceBadge`
* `TrendIcon`

### Fetch

* `fetch('/api/items')`
* ISR ou SSR leve

---

# 📁 Estrutura de Pastas (Backend)

```
/cmd/api
/internal
  /domain
    item.go
    repository.go
  /application
    update_prices.go
    get_item.go
  /infrastructure
    osrs
    cache
    repository
  /interfaces
    http
      handlers
      routes.go
```

---

# 🧪 Testes (desde o começo)

* Domain: puro (100% testável)
* Use Cases: mock de interfaces
* Infra: testes de integração simples

---

# 🚀 Roadmap pós-MVP

* Histórico de preços (gráfico)
* Alertas (sem login → query params)
* “Top gainers / losers”
* Export CSV
* Dark mode

---

# 🧠 Prompt ideal pra jogar no Cursor (exemplo)

> “Crie o scaffolding inicial de um backend em Go seguindo Clean Architecture para um projeto chamado OSRS Market Watch. Use chi como router, implemente as entidades de domínio, interfaces (PriceProvider, ItemRepository), um cliente para consumir a API do OSRS Wiki e um endpoint GET /items que retorna preços mockados inicialmente.”


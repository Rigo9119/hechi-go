# hechi-go

Personal finance backend built with Go, GraphQL, and PostgreSQL.

## Stack

| Layer | Technology |
|-------|-----------|
| HTTP | [Chi v5](https://github.com/go-chi/chi) |
| GraphQL | [gqlgen](https://gqlgen.com) |
| Database | PostgreSQL via [pgx v5](https://github.com/jackc/pgx) |
| Auth | JWT via [golang-jwt](https://github.com/golang-jwt/jwt) |
| Money | [shopspring/decimal](https://github.com/shopspring/decimal) |

## Features

- User registration and login with JWT authentication
- Bank accounts (checking, savings, credit, investment)
- Transactions (income, expense, transfer) with automatic balance updates
- GraphQL API with playground for development

## Project Structure

```
hechi-go/
├── main.go                     # Entry point, Chi router, GraphQL handler
├── graph/
│   ├── schema/schema.graphqls  # GraphQL schema
│   ├── resolver.go             # Root resolver with dependencies
│   ├── schema.resolvers.go     # Resolver implementations
│   ├── generated/              # gqlgen-generated execution engine (do not edit)
│   └── model/                  # gqlgen-generated Go types (do not edit)
├── internal/
│   ├── auth/                   # JWT service and Chi middleware
│   ├── config/                 # Environment variable configuration
│   ├── domain/                 # Core types: User, Account, Transaction
│   └── repository/             # PostgreSQL queries (pgx)
├── migrations/                 # SQL migration files
├── gqlgen.yml                  # gqlgen configuration
└── Makefile                    # Common tasks
```

## Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL 14+

### Setup

**1. Clone and install dependencies:**
```bash
git clone <repo-url>
cd hechi-go
go mod download
```

**2. Create the database:**
```bash
createdb hechi
```

**3. Configure environment:**
```bash
cp .env.example .env
```

Edit `.env` and set at minimum:
```env
DATABASE_URL=postgres://user:password@localhost:5432/hechi?sslmode=disable
JWT_SECRET=your-long-random-secret-here
```

**4. Run migrations:**
```bash
DATABASE_URL=<your-url> make migrate-up
```

**5. Start the server:**
```bash
go run .
```

The server starts on `http://localhost:8080` by default.

| Endpoint | Description |
|----------|-------------|
| `POST /graphql` | GraphQL API |
| `GET /playground` | GraphQL playground (development) |
| `GET /health` | Health check |

## GraphQL API

### Authentication

All queries and mutations except `register` and `login` require a `Bearer` token in the `Authorization` header.

### Example Operations

**Register:**
```graphql
mutation {
  register(input: {
    email: "user@example.com"
    password: "secret"
    name: "Jane Doe"
  }) {
    token
    user { id email name }
  }
}
```

**Login:**
```graphql
mutation {
  login(input: { email: "user@example.com", password: "secret" }) {
    token
  }
}
```

**Create an account:**
```graphql
mutation {
  createAccount(input: {
    name: "Main Checking"
    type: CHECKING
    currency: "USD"
    initialBalance: "1000.00"
  }) {
    id name balance currency
  }
}
```

**Create a transaction:**
```graphql
mutation {
  createTransaction(input: {
    accountId: "<account-id>"
    amount: "50.00"
    type: EXPENSE
    category: "Groceries"
    description: "Weekly shopping"
    date: "2026-04-30T12:00:00Z"
  }) {
    id amount type category date
  }
}
```

**Fetch all accounts:**
```graphql
query {
  accounts {
    id name type balance currency createdAt
  }
}
```

**Fetch transactions for an account:**
```graphql
query {
  transactions(accountId: "<account-id>", limit: 20, offset: 0) {
    id amount type category description date
  }
}
```

## Development

**Regenerate GraphQL code after schema changes:**
```bash
make generate
```

After running generate, implement any new resolver stubs that appear at the bottom of `graph/schema.resolvers.go`.

**Build:**
```bash
make build   # outputs to bin/hechi-go
```

**Roll back migrations:**
```bash
DATABASE_URL=<your-url> make migrate-down
```

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://postgres:postgres@localhost:5432/hechi?sslmode=disable` | PostgreSQL connection string |
| `JWT_SECRET` | `change-me-in-production` | Secret key for JWT signing |
| `JWT_EXPIRY_HOURS` | `24` | Token expiry in hours |
| `PORT` | `8080` | HTTP server port |

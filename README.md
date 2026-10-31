# kl-assignment

Backend application built with **Go + Gin + GORM + MySQL**, deployable via Docker.

## Prerequisites

- Docker Compose
- Go 1.27

## Project Structure

```
.
├── app/                  # Domain modules (handler / service / repository)
│   ├── account/
│   ├── banner/
│   ├── debit_card/
│   ├── health/
│   ├── transaction/
│   └── user/
├── cmd/api/              # Entry point
├── migrations/           # SQL schema + mock data migrations
├── mock/                 # Raw SQL files & UI template reference
├── pkg/                  # Shared packages (config, db, logger, middleware, response, validation)
├── router/               # Route registration
├── docker-compose.yml
├── Dockerfile
└── .env.example
```

## Setup

### 1. Configure environment

```bash
cp .env.example .env
```

Edit `.env` as needed:
> [!NOTE] if run in local change `DB_HOST` to localhost
| Variable             | Description               | Default              |
|----------------------|---------------------------|----------------------|
| `PORT`               | API server port           | `8080`               |
| `DB_HOST`            | MySQL host                | `kl-assignment-db`   |
| `DB_PORT`            | MySQL port                | `3306`               |
| `DB_NAME`            | Database name             | `kl-assignment-db`   |
| `DB_USER`            | Database user             | `appuser`            |
| `DB_PASSWORD`        | Database password         | `my-secret-password` |
| `MYSQL_ROOT_PASSWORD`| MySQL root password       | `my-root-secret`     |
| `DB_DEBUG`           | Enable GORM SQL query log | `false`              |

### 2. Run with Docker (recommended)

```bash
# Start DB + run migrations + start API
docker compose up -d

# Run migrations only (if needed separately)
docker compose up -d migrate
```

### 3. Run locally

```bash
# Start MySQL separately, then:
make run
# or
go run cmd/api/main.go
```

## API Documentation

Base URL: `http://localhost:8080/api`

### Health

#### `GET /health`

Check service health.

**Response 200**
```json
{ "message": "ok" }
```

---

### User

#### `GET /user/:user_id`

Get user profile with greeting message.

**Example**
```bash
curl http://localhost:8080/api/user/000093b8e1a211ef95a30242ac180002
```

**Response 200**
```json
{
  "message": "success",
  "data": {
    "user_id": "000093b8e1a211ef95a30242ac180002",
    "name": "Clare",
    "greeting_msg": "Have a nice day"
  }
}
```

---

### Account

#### `GET /account/:user_id`

Get all accounts for a user including balance, details, and flags.

**Example**
```bash
curl http://localhost:8080/api/account/000093b8e1a211ef95a30242ac180002
```

**Response 200**
```json
{
  "message": "success",
  "data": [
    {
      "account_id": "5a5ae9dae1a111ef80ca0242ac180002",
      "account_number": "568-2-81740-9",
      "amount": "62000.00",
      "currency": "THB",
      "user_id": "000093b8e1a211ef95a30242ac180002",
      "type": "Saving Account",
      "issuer": "TestLab",
      "color": "#ffffff",
      "is_main_account": true,
      "progress": 0,
      "flags": [
        { "flag_type": "lock", "flag_value": "Disbursement" }
      ]
    }
  ]
}
```

#### `PUT /account/main`

Set an account as the main account for a user. Automatically unsets the previous main account (atomic transaction).

**Request Body**
```json
{
  "user_id": "000093b8e1a211ef95a30242ac180002",
  "account_id": "5a5ae9dae1a111ef80ca0242ac180002"
}
```

**Response 200**
```json
{ "message": "success", "data": "" }
```

---

### Debit Card

#### `GET /debit-card/:user_id`

Get all debit cards for a user including design and status.

**Example**
```bash
curl http://localhost:8080/api/debit-card/000093b8e1a211ef95a30242ac180002
```

**Response 200**
```json
{
  "message": "success",
  "data": [
    {
      "card_id": "abc123",
      "name": "My Salary",
      "status": "Active",
      "detail": {
        "issuer": "TestLab",
        "number": "9440 78•• •••• 3115"
      },
      "design": {
        "color": "#00a1e2",
        "border_color": "#00a1e2"
      }
    }
  ]
}
```

#### `PUT /debit-card/status`

Update the status of a debit card.

**Request Body**
```json
{
  "user_id": "000093b8e1a211ef95a30242ac180002",
  "card_id": "abc123",
  "status": "Active"
}
```

`status` accepts: `Active` | `InActive`

**Response 200**
```json
{ "message": "success", "data": "" }
```

---

### Banner

#### `GET /banner/`

Get all banners.

**Response 200**
```json
{
  "message": "success",
  "data": [
    {
      "banner_id": "b1",
      "title": "Want some money?",
      "description": "You can start apply",
      "image": "https://..."
    }
  ]
}
```

#### `GET /banner/:user_id`

Get banners targeted to a specific user.

**Example**
```bash
curl http://localhost:8080/api/banner/000093b8e1a211ef95a30242ac180002
```

---

### Transaction

#### `GET /transaction/:user_id?page=1&limit=10`

Get paginated recent transactions for a user.

**Query Parameters**

| Param   | Type | Default | Description        |
|---------|------|---------|--------------------|
| `page`  | int  | `1`     | Page number        |
| `limit` | int  | `10`    | Items per page     |

**Example**
```bash
curl "http://localhost:8080/api/transaction/000093b8e1a211ef95a30242ac180002?page=1&limit=10"
```

**Response 200**
```json
{
  "message": "success",
  "data": [
    {
      "transaction_id": "txn001",
      "name": "Emily",
      "image": "https://...",
      "is_bank": false
    }
  ],
  "page": 1,
  "limit": 10,
  "total_rows": 50,
  "total_pages": 5
}
```

---

## Database Schema

Schema migrations are in `migrations/`. They run automatically on `docker compose up`.

| Table                | Description                        |
|----------------------|------------------------------------|
| `users`              | User basic info                    |
| `user_greetings`     | Personalized greeting messages     |
| `accounts`           | Account core info (type, currency) |
| `account_balances`   | Account balance                    |
| `account_details`    | Account color, main flag, progress |
| `account_flags`      | Account label flags (lock, overdue)|
| `debit_cards`        | Debit card basic info              |
| `debit_card_details` | Card number, issuer                |
| `debit_card_design`  | Card color, border color           |
| `debit_card_status`  | Card active/inactive status        |
| `banners`            | Promotional banner content         |
| `transactions`       | Recent transaction list            |

---

## Unit Testing

Run all unit tests across all domain modules:

```bash
go test -v ./...
```

Run test coverage:

```bash
go test -cover ./app/...
```

Test coverage results:
- `app/account`: 100.0%
- `app/banner`: 100.0%
- `app/debit_card`: 96.8%
- `app/health`: 100.0%
- `app/transaction`: 100.0%
- `app/user`: 97.8%

---

## Stress Testing (k6)

A stress test script is provided in `k6-test.js` simulating virtual users (VUs) ramping up to 50 concurrent users.

### Run with local k6:

```bash
k6 run k6-test.js
```

### Run with Docker k6:

```bash
docker run --rm -i --network=host grafana/k6 run - k6-test.js
```

Custom URL or User ID:
```bash
k6 run -e BASE_URL=http://localhost:8080 -e USER_ID=000018b0e1a211ef95a30242ac180002 k6-test.js
```


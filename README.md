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

### Schema Optimization Recommendations
```sql
-- Speed up account listing & main account atomic update
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_account_details_user_id ON account_details(user_id);
CREATE INDEX idx_account_flags_account_id ON account_flags(account_id);

-- Speed up debit card queries
CREATE INDEX idx_debit_cards_user_id ON debit_cards(user_id);

-- Speed up transaction pagination
CREATE INDEX idx_transactions_user_id ON transactions(user_id);

-- Speed up public & user targeted banners
CREATE INDEX idx_banners_user_id ON banners(user_id);
```
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

A stress test script is provided in `k6-test.js` covering **all endpoints** (GET & PUT) with 50 concurrent virtual users (VUs) ramp-up. It automatically produces an HTML report (`summary.html`) and terminal output upon completion.

### Tested Endpoints:
1. `GET /api/health`
2. `GET /api/user/:user_id`
3. `GET /api/banner/`
4. `GET /api/banner/:user_id`
5. `GET /api/account/:user_id`
6. `PUT /api/account/main`
7. `GET /api/debit-card/:user_id`
8. `PUT /api/debit-card/status`
9. `GET /api/transaction/:user_id?page=1&limit=10`

### Execution:

```bash
# Run with local k6:
k6 run k6-test.js

# build k6 docker
docker compose up -d --build k6 

# start k6
docker start k6

# Custom host / user:
k6 run -e BASE_URL=http://localhost:8080 -e USER_ID=000018b0e1a211ef95a30242ac180002 k6-test.js
```

### Stress Test Sample Report (50 VUs Ramp-up) 

[!NOTE] **local Test 1 min**

- **Test Duration:** 1m 20s (10s ramp-up $\rightarrow$ 60s peak 50 VUs $\rightarrow$ 10s ramp-down)
    account_get_requests...........: 2218    27.490962/s
    account_set_main_requests......: 2218    27.490962/s
    banner_public_requests.........: 2218    27.490962/s
    banner_user_requests...........: 2218    27.490962/s
    health_requests................: 2218    27.490962/s
    debit_card_status_requests.....: 2218    27.490962/s
    debit_card_get_requests........: 2218    27.490962/s
    transaction_requests...........: 2218    27.490962/s
    user_profile_requests..........: 2218    27.490962/s
    checks.........................: 100.00% ✓ 19962      ✗ 0    
    data_received..................: 6.5 MB  80 kB/s
    data_sent......................: 2.6 MB  33 kB/s
    group_duration.................: min=121.81µs max=566.48ms avg=17.59ms med=1.68ms  p(90)=34.06ms p(95)=65.56ms  p(99)=315.5ms 
    http_req_blocked...............: min=569ns    max=8.46ms   avg=5.12µs  med=2.35µs  p(90)=4.54µs  p(95)=6.7µs    p(99)=22.7µs  
    http_req_connecting............: min=0s       max=2.58ms   avg=769ns   med=0s      p(90)=0s      p(95)=0s       p(99)=0s      
    http_req_duration..............: min=90.59µs  max=493.33ms avg=11.59ms med=2.01ms  p(90)=20.42ms p(95)=36.61ms  p(99)=220.28ms
    { expected_response:true }...: min=90.59µs  max=493.33ms avg=11.59ms med=2.01ms  p(90)=20.42ms p(95)=36.61ms  p(99)=220.28ms
    http_req_failed................: 0.00%   ✓ 0          ✗ 19962
    http_req_receiving.............: min=6.53µs   max=11.28ms  avg=55.3µs  med=36.31µs p(90)=85.9µs  p(95)=118.89µs p(99)=275.78µs
    http_req_sending...............: min=2.11µs   max=5.05ms   avg=14.82µs med=7.98µs  p(90)=22.07µs p(95)=29.72µs  p(99)=87.82µs 
    http_req_tls_handshaking.......: min=0s       max=0s       avg=0s      med=0s      p(90)=0s      p(95)=0s       p(99)=0s      
    http_req_waiting...............: min=71.13µs  max=493.28ms avg=11.52ms med=1.91ms  p(90)=20.35ms p(95)=36.53ms  p(99)=220.19ms
    http_reqs......................: 19962   247.418657/s
    iteration_duration.............: min=1.02s    max=1.63s    avg=1.1s    med=1.05s   p(90)=1.28s   p(95)=1.38s    p(99)=1.55s   
    iterations.....................: 2218    27.490962/s
    vus............................: 3       min=2        max=50 
    vus_max........................: 50      min=50       max=50 

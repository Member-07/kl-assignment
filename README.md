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

> [!NOTE]
> If got error range of CPUs you can adjust it yourself
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
CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_account_details_user_id ON account_details(user_id);
CREATE INDEX idx_account_flags_account_id ON account_flags(account_id);
CREATE INDEX idx_debit_cards_user_id ON debit_cards(user_id);
CREATE INDEX idx_transactions_user_id ON transactions(user_id);
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

A stress test script is provided in `k6-test.js` covering **all endpoints** (GET & PUT) with 50 concurrent virtual users (VUs) ramp-up. It automatically produces an HTML report (`reports/summary.html`) and terminal output upon completion.

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

# start docker k6
docker start k6

# Custom host / user:
k6 run -e BASE_URL=http://localhost:8080 -e USER_ID=000018b0e1a211ef95a30242ac180002 k6-test.js
```

### Stress Test Sample Report (50 VUs Ramp-up) 

> [!NOTE]
> Local Test (5 min continuous peak load)

- **Test Duration:** 5m 20s (10s ramp-up -> 600s peak 50 VUs -> 10s ramp-down)
- **Docker Container Resource Monitoring (at peak):**
  - CPU Reservation: 1.5 Cores
  - CPU Limit: 2 Cores
  - Actual CPU Usage: 1.313 Cores (65.65% of Limit)
  - Memory Usage: 565.1 MiB / 1.5 GiB (36.79%)

- **k6 Execution Results:**

```text
  account_get_requests...........: 9522    29.682913/s
  account_set_main_requests......: 9522    29.682913/s
  banner_public_requests.........: 9522    29.682913/s
  banner_user_requests...........: 9522    29.682913/s
  debit_card_get_requests........: 9522    29.682913/s
  debit_card_status_requests.....: 9522    29.682913/s
  transaction_requests...........: 9522    29.682913/s
  checks.........................: 100.00% ✓ 76176      ✗ 0    
  data_received..................: 26 MB   81 kB/s
  data_sent......................: 11 MB   33 kB/s
  group_duration.................: min=331.84µs max=1.19s    avg=25.29ms med=2.98ms  p(90)=50.69ms p(95)=117.8ms  p(99)=413.18ms
  http_req_blocked...............: min=621ns    max=2.93ms   avg=3.87µs  med=2.29µs  p(90)=4.45µs  p(95)=6.22µs   p(99)=19.62µs 
  http_req_connecting............: min=0s       max=2.43ms   avg=162ns   med=0s      p(90)=0s      p(95)=0s       p(99)=0s      
  ✓ http_req_duration..............: min=51.74µs  max=786.63ms avg=15.67ms med=3.49ms  p(90)=27.85ms p(95)=58.12ms  p(99)=276.16ms
  { expected_response:true }...: min=51.74µs  max=786.63ms avg=15.67ms med=3.49ms  p(90)=27.85ms p(95)=58.12ms  p(99)=276.16ms
  ✓ http_req_failed................: 0.00%   ✓ 0          ✗ 76176
  http_req_receiving.............: min=6.33µs   max=11.82ms  avg=59.79µs med=36.63µs p(90)=87.98µs p(95)=126.04µs p(99)=347.26µs
  http_req_sending...............: min=2.05µs   max=9.52ms   avg=14.25µs med=7.76µs  p(90)=21.58µs p(95)=28.76µs  p(99)=78.54µs 
  http_req_tls_handshaking.......: min=0s       max=0s       avg=0s      med=0s      p(90)=0s      p(95)=0s       p(99)=0s      
  http_req_waiting...............: min=0s       max=786.51ms avg=15.6ms  med=3.43ms  p(90)=27.73ms p(95)=58ms     p(99)=276.12ms
  iteration_duration.............: min=1.02s    max=2.48s    avg=1.12s   med=1.06s   p(90)=1.3s    p(95)=1.45s    p(99)=1.69s   
  http_reqs......................: 76176   237.463303/s
  iterations.....................: 9522    29.682913/s
```

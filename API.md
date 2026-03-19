# API Reference

Base URL: `/api/v1`

## Conventions
- Success envelope:
```json
{
  "data": {}
}
```
- Paginated envelope:
```json
{
  "data": [],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 42
  }
}
```
- Error envelope:
```json
{
  "error": {
    "code": "ERR_CODE",
    "message": "Human readable message"
  }
}
```

## Health
### `GET /health`
Response:
```json
{ "status": "ok" }
```

## Auth
### `POST /auth/sync`
Request:
```json
{
  "firebase_token": "eyJhbGci...",
  "display_name": "Jane Doe",
  "email": "jane@example.com",
  "currency_code": "INR"
}
```
Response:
```json
{
  "data": {
    "user": {
      "id": "e6a8...",
      "firebase_uid": "firebase-uid",
      "email": "jane@example.com",
      "display_name": "Jane Doe",
      "currency_code": "INR",
      "created_at": "2026-03-15T00:00:00Z",
      "updated_at": "2026-03-15T00:00:00Z"
    }
  }
}
```

## Users (Bearer token required)
### `GET /users/me`
Response:
```json
{
  "data": {
    "id": "e6a8...",
    "firebase_uid": "firebase-uid",
    "email": "jane@example.com",
    "display_name": "Jane Doe",
    "currency_code": "INR"
  }
}
```

### `PATCH /users/me`
Request:
```json
{
  "display_name": "Jane D",
  "currency_code": "USD"
}
```
Response:
```json
{
  "data": {
    "id": "e6a8...",
    "display_name": "Jane D",
    "currency_code": "USD"
  }
}
```

## Accounts (Bearer token required)
### `GET /accounts`
Response:
```json
{
  "data": [
    {
      "id": "acc-1",
      "name": "Cash",
      "account_type": "general",
      "opening_balance": 1000,
      "current_balance": 850
    }
  ]
}
```

### `POST /accounts`
Request:
```json
{
  "name": "Savings",
  "opening_balance": 5000
}
```
Response:
```json
{
  "data": {
    "id": "acc-2",
    "name": "Savings",
    "opening_balance": 5000,
    "current_balance": 5000
  }
}
```

### `PATCH /accounts/:id`
Request:
```json
{
  "name": "Emergency Fund"
}
```
Response:
```json
{
  "data": {
    "id": "acc-2",
    "name": "Emergency Fund"
  }
}
```

### `DELETE /accounts/:id`
Response: `204 No Content`

## Categories and Subcategories (Bearer token required)
### `GET /categories`
Response:
```json
{
  "data": [
    {
      "id": "cat-1",
      "name": "Food",
      "subcategories": [
        { "id": "sub-1", "name": "Groceries" }
      ]
    }
  ]
}
```

### `POST /categories`
Request:
```json
{ "name": "Transport" }
```
Response:
```json
{ "data": { "id": "cat-2", "name": "Transport" } }
```

### `PATCH /categories/:id`
Request:
```json
{ "name": "Travel" }
```
Response:
```json
{ "data": { "id": "cat-2", "name": "Travel" } }
```

### `DELETE /categories/:id`
Response: `204 No Content`

### `POST /categories/:id/subcategories`
Request:
```json
{ "name": "Metro" }
```
Response:
```json
{ "data": { "id": "sub-2", "category_id": "cat-2", "name": "Metro" } }
```

### `PATCH /subcategories/:id`
Request:
```json
{ "name": "Bus" }
```
Response:
```json
{ "data": { "id": "sub-2", "name": "Bus" } }
```

### `DELETE /subcategories/:id`
Response: `204 No Content`

## Transactions (Bearer token required)
### `GET /transactions`
Query params:
- `account_id`, `category_id`, `subcategory_id`
- `date_from`, `date_to`
- `search`
- `type` = `income|expense|all`
- `sort_by` = `transaction_date|amount|name|created_at`
- `sort_dir` = `asc|desc`
- `page`, `per_page`
- `passbook_mode` = `true|false`

Response:
```json
{
  "data": [
    {
      "id": "txn-1",
      "account_id": "acc-1",
      "category_id": "cat-1",
      "subcategory_id": "sub-1",
      "name": "Groceries",
      "amount": -1200,
      "transaction_date": "2026-03-15",
      "notes": "Weekly purchase",
      "is_scheduled": false
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1 }
}
```

### `POST /transactions`
Request:
```json
{
  "account_id": "acc-1",
  "category_id": "cat-1",
  "subcategory_id": "sub-1",
  "name": "Groceries",
  "amount": -1200,
  "transaction_date": "2026-03-15",
  "notes": "Weekly purchase",
  "is_scheduled": false
}
```
Response:
```json
{ "data": { "id": "txn-2", "name": "Groceries", "amount": -1200 } }
```

### `GET /transactions/:id`
Response:
```json
{ "data": { "id": "txn-2", "name": "Groceries", "amount": -1200 } }
```

### `PATCH /transactions/:id`
Request:
```json
{ "amount": -1000, "notes": "Discount applied" }
```
Response:
```json
{ "data": { "id": "txn-2", "amount": -1000, "notes": "Discount applied" } }
```

### `DELETE /transactions/:id`
Response: `204 No Content`

## Quick Transactions (Bearer token required)
### `GET /quick-transactions`
Response:
```json
{
  "data": [
    {
      "id": "qt-1",
      "label": "Lunch",
      "account_id": "acc-1",
      "category_id": "cat-1",
      "subcategory_id": "sub-1",
      "amount": -150,
      "sort_order": 1
    }
  ]
}
```

### `POST /quick-transactions`
Request:
```json
{
  "label": "Fuel",
  "account_id": "acc-1",
  "category_id": "cat-2",
  "subcategory_id": "sub-3",
  "amount": -2000,
  "notes": "Car"
}
```
Response:
```json
{ "data": { "id": "qt-2", "label": "Fuel" } }
```

### `PATCH /quick-transactions/:id`
Request:
```json
{ "label": "Fuel Full Tank", "amount": -2400 }
```
Response:
```json
{ "data": { "id": "qt-2", "label": "Fuel Full Tank", "amount": -2400 } }
```

### `DELETE /quick-transactions/:id`
Response: `204 No Content`

### `POST /quick-transactions/:id/execute`
Request:
```json
{}
```
Response:
```json
{
  "data": {
    "transaction_id": "txn-99",
    "message": "Quick transaction executed"
  }
}
```

### `PATCH /quick-transactions/reorder`
Request:
```json
{ "ordered_ids": ["qt-2", "qt-1"] }
```
Response:
```json
{ "data": { "ordered_ids": ["qt-2", "qt-1"] } }
```

## Stats (Bearer token required)
### `GET /stats/summary`
Query params: `period`, `value`, optional `account_id`

Response:
```json
{
  "data": {
    "total_income": 12000,
    "total_expense": -9500,
    "net": 2500,
    "category_breakdown_expense": [
      { "category": "Food", "subcategory": "Groceries", "amount": -3200, "percentage": 33.68 }
    ],
    "category_breakdown_income": [
      { "category": "Salary", "subcategory": "Monthly", "amount": 12000, "percentage": 100 }
    ],
    "timeseries": [
      { "label": "2026-03-01", "income": 4000, "expense": -1200 }
    ]
  }
}
```

### `GET /stats/export/pdf`
Query params: `period`, `value`, optional `account_id`

Response: `200 OK` with `application/pdf` binary stream.

### `GET /stats/compare`
Query params:
- `period` = `month|year` (required)
- `value1` — first period value, e.g. `2026-03` (required)
- `value2` — second period value, e.g. `2026-02` (required)
- `account_id` — optional account filter

Response:
```json
{
  "data": {
    "period1": {
      "label": "2026-03",
      "total_income": 12000,
      "total_expense": -9500,
      "net": 2500,
      "top_expense_categories": [
        { "category": "Food", "subcategory": "Groceries", "amount": -3200, "percentage": 33.68 }
      ],
      "top_income_categories": [
        { "category": "Salary", "subcategory": "Monthly", "amount": 12000, "percentage": 100 }
      ]
    },
    "period2": {
      "label": "2026-02",
      "total_income": 10000,
      "total_expense": -8000,
      "net": 2000,
      "top_expense_categories": [],
      "top_income_categories": []
    },
    "income_change_pct": 20.0,
    "expense_change_pct": 18.75,
    "net_change_pct": 25.0
  }
}
```

## Chat (Bearer token required)
### `POST /chat`
Natural-language chatbot that can record transactions, check balances, and answer questions.
Uses pattern matching by default; optionally connects to a local Ollama LLM if `OLLAMA_URL` is set.

Request:
```json
{
  "message": "Spent 500 on groceries from Cash"
}
```

Response (with action):
```json
{
  "data": {
    "reply": "✅ Recorded: expense of 500.00 for groceries (Cash)",
    "action": {
      "type": "expense",
      "name": "groceries",
      "amount": 500,
      "account": "Cash",
      "category": "",
      "date": ""
    },
    "success": true
  }
}
```

Response (informational):
```json
{
  "data": {
    "reply": "Your current balances:\n• Cash: ₹850.00\n• Savings: ₹5,000.00",
    "action": null,
    "success": true
  }
}
```

Supported patterns:
- `Spent X on Y [from ACCOUNT]`
- `Paid X for Y [from ACCOUNT]`
- `Bought Y for X [from ACCOUNT]`
- `Got/Received/Earned X [from] Y [from ACCOUNT]`
- `balance` / `how much` — shows account balances
- `recent` / `last transactions` — shows recent transactions
- `help` / `what can you do` — lists available commands

## Audit
### `GET /audit`
Current behavior:
```json
{
  "error": {
    "code": "ERR_NOT_IMPLEMENTED",
    "message": "audit endpoint will be enabled after audit query handler wiring"
  }
}
```

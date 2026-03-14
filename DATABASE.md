# Database Design

ledgerA uses PostgreSQL with UUID primary keys, soft-delete columns on mutable business tables, and migration-driven schema changes.

## Schema Summary

### `users`
- `id` UUID PK
- `firebase_uid` unique text
- `email`, `display_name`, `currency_code`
- `created_at`, `updated_at`, `deleted_at`

Indexes:
- `idx_users_firebase_uid (firebase_uid)`

### `accounts`
- `id` UUID PK
- `user_id` FK -> `users(id)`
- `name`, `account_type`
- `opening_balance`, `current_balance` as `numeric(20,4)`
- timestamps + `deleted_at`

Indexes:
- `idx_accounts_user_id (user_id)` filtered by `deleted_at IS NULL`

### `categories`
- `id` UUID PK
- `user_id` FK
- `name`
- timestamps + `deleted_at`

Indexes:
- `idx_categories_user_name (user_id, name)` unique, filtered on non-deleted rows

### `subcategories`
- `id` UUID PK
- `category_id` FK -> `categories(id)`
- `user_id` FK -> `users(id)`
- `name`
- timestamps + `deleted_at`

Indexes:
- `idx_subcategories_category_name (category_id, name)` unique, filtered on non-deleted rows

### `transactions`
- `id` UUID PK
- `user_id`, `account_id`, `category_id`, `subcategory_id` FKs
- `name`, `amount`, `transaction_date`
- optional `notes`, `is_scheduled`
- timestamps + `deleted_at`

Indexes:
- `idx_txn_user_date (user_id, transaction_date DESC)` filtered on active rows
- `idx_txn_account (account_id)` filtered on active rows
- `idx_txn_category (category_id)` filtered on active rows
- `idx_txn_name_search (user_id, name text_pattern_ops)`

### `quick_transactions`
- `id` UUID PK
- `user_id` FK
- optional references: `account_id`, `category_id`, `subcategory_id`
- `label`, optional `amount`, optional `notes`
- `sort_order`
- timestamps + `deleted_at`

### `audit_log`
- `id` UUID PK
- `user_id` FK
- `entity_type`, `entity_id`, `action`
- `diff` JSONB
- `created_at`

Indexes:
- `idx_audit_entity (entity_id, entity_type)`

## Balance Maintenance Strategy
Account balances are maintained as a stored `current_balance` and updated transactionally in service-layer operations:
- Create transaction: `current_balance += amount`
- Update transaction: `current_balance += (new_amount - old_amount)`
- Delete transaction: `current_balance -= amount`

This denormalized balance strategy is chosen for fast account listing and passbook reads. The source of truth remains the transaction ledger; if required, a reconciliation job can recompute balances from transaction history.

## Soft Delete Behavior
Soft deletes are used for:
- users
- accounts
- categories
- subcategories
- transactions
- quick_transactions

Rows are marked with `deleted_at` and excluded by filtered indexes and repository queries.

`audit_log` is append-only and not soft-deleted.

## Migration Order
1. `000001_create_users`
2. `000002_create_accounts`
3. `000003_create_categories` (+ subcategories)
4. `000004_create_transactions`
5. `000005_create_quick_transactions`
6. `000006_create_audit_log`

## Operational Notes
- Keep all monetary fields in `numeric(20,4)` from DB through service boundary.
- Use explicit transactions for any write path that also mutates account balances and audit entries.
- Maintain foreign key integrity for category/subcategory/account references before write.

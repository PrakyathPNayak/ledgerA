CREATE TABLE IF NOT EXISTS recurring_transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id),
    account_id uuid NOT NULL REFERENCES accounts(id),
    category_id uuid NOT NULL REFERENCES categories(id),
    subcategory_id uuid NOT NULL REFERENCES subcategories(id),
    name text NOT NULL,
    amount numeric(20,4) NOT NULL,
    notes text,
    frequency text NOT NULL CHECK (frequency IN ('daily', 'weekly', 'monthly', 'yearly')),
    start_date date NOT NULL,
    next_due_date date NOT NULL,
    end_date date,
    is_active boolean NOT NULL DEFAULT true,
    last_executed_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_recurring_user ON recurring_transactions(user_id) WHERE deleted_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_recurring_due ON recurring_transactions(next_due_date) WHERE deleted_at IS NULL AND is_active = true;

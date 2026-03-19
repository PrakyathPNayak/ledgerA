CREATE TABLE IF NOT EXISTS budgets (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id),
    category_id uuid NOT NULL REFERENCES categories(id),
    amount numeric(20,4) NOT NULL,
    period text NOT NULL CHECK (period IN ('monthly', 'yearly')),
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_budget_user_category_period ON budgets(user_id, category_id, period) WHERE deleted_at IS NULL;

CREATE TABLE IF NOT EXISTS quick_transactions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL REFERENCES users(id),
    label text NOT NULL,
    account_id uuid REFERENCES accounts(id),
    category_id uuid REFERENCES categories(id),
    subcategory_id uuid REFERENCES subcategories(id),
    amount numeric(20,4),
    notes text,
    sort_order integer NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TABLE IF NOT EXISTS withdrawals
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    withdrawn_by uuid NOT NULL,
    status INTEGER DEFAULT 0,
    amount DECIMAL DEFAULT 0.0::DECIMAL,
    withdrawn_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    reference_id uuid NOT NULL UNIQUE
);

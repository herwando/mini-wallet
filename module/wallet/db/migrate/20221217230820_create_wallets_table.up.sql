CREATE TABLE IF NOT EXISTS wallets
(
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owned_by uuid NOT NULL,
    status INTEGER DEFAULT 0,
    balance DECIMAL DEFAULT 0.0::DECIMAL,
    enabled_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add credits to users
ALTER TABLE users ADD COLUMN IF NOT EXISTS credits INTEGER NOT NULL DEFAULT 0;

-- Update plans table
ALTER TABLE plans ADD COLUMN IF NOT EXISTS credits INTEGER NOT NULL DEFAULT 0;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS features TEXT;
ALTER TABLE plans ADD COLUMN IF NOT EXISTS popular BOOLEAN NOT NULL DEFAULT false;

-- Payment configs table (易支付)
CREATE TABLE IF NOT EXISTS payment_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL DEFAULT '易支付',
    merchant_id VARCHAR(255) NOT NULL DEFAULT '',
    merchant_key VARCHAR(255) NOT NULL DEFAULT '',
    api_url VARCHAR(255) NOT NULL DEFAULT '',
    notify_url VARCHAR(255) NOT NULL DEFAULT '',
    return_url VARCHAR(255) NOT NULL DEFAULT '',
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Insert default payment config
INSERT INTO payment_configs (name, merchant_id, merchant_key, api_url, notify_url, return_url, is_active)
VALUES ('易支付', '', '', '', '', '', false)
ON CONFLICT DO NOTHING;

-- CDK codes table
CREATE TABLE IF NOT EXISTS cdk_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    credits INTEGER NOT NULL DEFAULT 0,
    max_uses INTEGER NOT NULL DEFAULT 1,
    used_count INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CDK usage records
CREATE TABLE IF NOT EXISTS cdk_usages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cdk_id UUID NOT NULL REFERENCES cdk_codes(id),
    user_id UUID NOT NULL REFERENCES users(id),
    used_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Update orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS credits_amount INTEGER NOT NULL DEFAULT 0;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS cdk_id UUID REFERENCES cdk_codes(id);

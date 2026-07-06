CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE payment_status AS ENUM ('pending','succeeded','failed');
CREATE TABLE payments (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), order_id VARCHAR(255) NOT NULL, stripe_intent_id VARCHAR(255) NOT NULL, amount NUMERIC(12,2) NOT NULL, currency VARCHAR(10) NOT NULL DEFAULT 'usd', status payment_status NOT NULL DEFAULT 'pending', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE INDEX idx_payments_order ON payments(order_id);

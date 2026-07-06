CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE order_status AS ENUM ('pending','confirmed','cancelled','payment_failed');
CREATE TABLE orders (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), user_id UUID NOT NULL, user_email VARCHAR(255) NOT NULL, status order_status NOT NULL DEFAULT 'pending', total NUMERIC(12,2) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE order_items (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE, product_id UUID NOT NULL, quantity INTEGER NOT NULL, price NUMERIC(12,2) NOT NULL);
CREATE INDEX idx_orders_user ON orders(user_id);

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE categories (id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name VARCHAR(100) UNIQUE NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
CREATE TABLE products (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(), name VARCHAR(200) NOT NULL, description TEXT NOT NULL,
	price NUMERIC(12,2) NOT NULL CHECK(price>0), stock INTEGER NOT NULL DEFAULT 0 CHECK(stock>=0),
	category_id UUID REFERENCES categories(id), is_active BOOLEAN NOT NULL DEFAULT true,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_products_category ON products(category_id);

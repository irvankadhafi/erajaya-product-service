-- +migrate Up notransaction
CREATE TABLE IF NOT EXISTS "products" (
    "id" BIGINT PRIMARY KEY,
    "name" TEXT,
    "slug" TEXT,
    "price" DECIMAL(10,2),
    "description" TEXT,
    "quantity" INT,
    "created_at" TIMESTAMP NOT NULL DEFAULT 'NOW()',
    "updated_at" TIMESTAMP NOT NULL DEFAULT 'NOW()'
);
ALTER TABLE IF EXISTS "products" ADD CONSTRAINT product_slug_unique UNIQUE (slug);

-- +migrate Down
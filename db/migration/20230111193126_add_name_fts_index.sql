-- +migrate Up notransaction
CREATE INDEX products_name_fts_index ON products USING gin(to_tsvector('english', name));

-- +migrate Down
DROP INDEX products_name_fts_index;

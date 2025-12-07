-- +migrate Up
ALTER TABLE order_products
ADD COLUMN image_path VARCHAR(255) DEFAULT NULL;

-- +migrate Down
ALTER TABLE order_products
DROP COLUMN image_path;

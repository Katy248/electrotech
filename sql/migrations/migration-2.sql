-- +migrate Up
ALTER TABLE order_products
DROP COLUMN quantity;

ALTER TABLE order_products
ADD COLUMN quantity REAL NOT NULL;

-- +migrate Down
ALTER TABLE order_products
DROP COLUMN quantity;

ALTER TABLE order_products
ADD COLUMN quantity INT NOT NULL;

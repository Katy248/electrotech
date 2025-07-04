-- name: InsertOrder :one

INSERT INTO orders (user_id, total_price) 
VALUES (@userId, @totalPrice)
RETURNING *;

-- name: AddOrderProduct :exec

INSERT INTO order_products (order_id, product_name, quantity, product_price) 
VALUES (@orderId, @productName, @quantity, @price);

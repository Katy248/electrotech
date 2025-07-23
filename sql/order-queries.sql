-- name: InsertOrder :one

INSERT INTO orders (user_id, creation_date) 
VALUES (@user_id, now())
RETURNING *;

-- name: AddOrderProduct :exec

INSERT INTO order_products (order_id, product_name, quantity, product_price) 
VALUES (@order_id, @product_name, @quantity, @price);

-- name: GetUserOrders :many

SELECT *
FROM orders o
    JOIN order_products p ON o.id = p.order_id
WHERE o.user_id = @user_id;


-- name: GetOrders :many

SELECT *
FROM orders o
WHERE user_id = @user_id;


-- name: GetOrderProducts :many

SELECT *
FROM order_products
WHERE order_id = @id;

-- name: GetUserByEmail :one

SELECT
    *
FROM
    users u
WHERE
    u.email = @email ;
    
-- name: InsertNewUser :exec

INSERT INTO users (
    email,
    password_hash,
    first_name,
    surname,
    last_name
) VALUES (
    @email,
    @passwordHash,
    @firstName,
    @surname,
    @lastName
);

-- name: UpdateUserEmail :exec

UPDATE users
SET
    email = @email
WHERE
    id = @id ;

-- name: UpdateUserPassword :exec

UPDATE users
SET
    password_hash = @passwordHash
WHERE
    id = @id ;

-- name: GetUserOrders :many

SELECT * 
FROM orders o 
WHERE o.user_id = @userId
LIMIT 40 * @page, 40;

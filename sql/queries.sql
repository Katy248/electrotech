-- name: GetByEmail :one

SELECT
    *
FROM
    users u
WHERE
    u.email = @email;
    
-- name: InsertNew :exec

INSERT INTO users (
    email,
    password_hash,
    first_name,
    surname,
    last_name,
    phone_number
) VALUES (
    @email,
    @password_hash,
    @first_name,
    @surname,
    @last_name,
    @phone_number
);

-- name: UpdateEmail :exec

UPDATE users
SET
    email = @email
WHERE
    id = @id;

-- name: UpdatePassword :exec

UPDATE users
SET
    password_hash = @password_hash
WHERE
    id = @id;

-- name: UpdateData :exec

UPDATE users
SET
    first_name = @first_name,
    surname = @surname,
    last_name = @last_name,
    phone_number = @phone_number
WHERE
    id = @id;

-- name: UpdateCompanyData :exec

UPDATE users
SET
    company_name = company_name,
    company_inn = company_inn,
    company_address = company_address,
    position_in_company = position_in_company
WHERE
    id = @id;


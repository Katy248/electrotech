-- name: GetByEmail :one

SELECT
    *
FROM
    users u
WHERE
    u.email = @email;

-- name: GetById :one

SELECT
    *
FROM
    users u
WHERE
    u.id = @id;
    
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
    email = @email;

-- name: UpdatePhoneNumber :exec

UPDATE users
SET
    phone_number = @phone_number
WHERE
    email = @email;

-- name: UpdateData :exec

UPDATE users
SET
    first_name = @first_name,
    surname = @surname,
    last_name = @last_name
WHERE
    email = @email;

-- name: UpdateCompanyData :exec

UPDATE users
SET
    company_name = @company_name,
    company_inn = @company_inn,
    company_address = @company_address,
    company_okpo = @company_okpo,
    position_in_company = @position_in_company
WHERE
    email = @email;

-- name: GetCompanyData :one

SELECT
    company_name,
    company_inn,
    company_address,
    company_okpo,
    position_in_company
FROM
    users
WHERE
    email = @email;

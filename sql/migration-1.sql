CREATE TABLE IF NOT EXISTS users (
    id INT PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone_number VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    -- Not necessary for viewing data
    company_name VARCHAR(255),
    company_inn VARCHAR(255),
    company_address VARCHAR(255),
    position_in_company VARCHAR(255) 
);

CREATE TABLE IF NOT EXISTS orders (
    id INT PRIMARY KEY,
    user_id INT NOT NULL,
    creation_date DATE NOT NULL,
    total_price DECIMAL(10, 2) NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS order_products (
    id INT PRIMARY KEY,
    order_id INT,
    product_name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    product_price DECIMAL NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders(id)
);

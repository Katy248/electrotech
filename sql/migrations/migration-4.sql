-- +migrate Up
CREATE TABLE requests (
    id INT PRIMARY KEY AUTOINCREMENT,
    creation_date DATE NOT NULL,
    person_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    message VARCHAR(1025) NOT NULL
);

-- +migrate Down
DROP TABLE requests;

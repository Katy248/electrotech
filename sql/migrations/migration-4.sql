-- +migrate Up
CREATE TABLE user_questions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    creation_date DATE NOT NULL,
    person_name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    phone VARCHAR(20),
    message VARCHAR(1025) NOT NULL,
    client_ip VARCHAR(20) NOT NULL DEFAULT '127.0.0.1'
);

-- +migrate Down
DROP TABLE user_questions;

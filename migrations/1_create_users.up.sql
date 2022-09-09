CREATE TABLE IF NOT EXISTS users (
    ID SERIAL PRIMARY KEY,
    first_name VARCHAR(40),
    last_name VARCHAR(40),
    email VARCHAR(40),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
)
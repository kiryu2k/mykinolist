-- CREATE TABLE lists (
--     id SERIAL PRIMARY KEY,
--     title VARCHAR(50) NOT NULL
-- );

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    hashed_password VARCHAR(80) NOT NULL,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
    -- list_id SERIAL REFERENCES lists (id) ON DELETE CASCADE 
);
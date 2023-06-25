CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    hashed_password VARCHAR(80) NOT NULL,
    created_on TIMESTAMP NOT NULL,
    last_login TIMESTAMP
);

CREATE TABLE lists (
    id SERIAL PRIMARY KEY,
    owner_id SERIAL REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE list_titles (
    list_id SERIAL REFERENCES lists (id) ON DELETE CASCADE,
    title_id SERIAL UNIQUE NOT NULL,
    status_name VARCHAR(15) NOT NULL,
    score SERIAL CHECK (score BETWEEN 0 AND 10),
    is_favorite BOOLEAN NOT NULL
);

CREATE TABLE tokens (
    user_id SERIAL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token VARCHAR(200) UNIQUE NOT NULL
);
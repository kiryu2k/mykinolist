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

CREATE TABLE title_status (
    id SERIAL PRIMARY KEY,
    status_name VARCHAR(15) UNIQUE NOT NULL
);

CREATE TABLE list_titles (
    list_id SERIAL REFERENCES lists (id) ON DELETE CASCADE,
    title_id SERIAL NOT NULL,
    status_id SERIAL REFERENCES title_status (id) ON DELETE CASCADE,
    score SERIAL CHECK (score BETWEEN 0 AND 10),
    is_favorite BOOLEAN NOT NULL
);

CREATE TABLE tokens (
    user_id SERIAL REFERENCES users (id) ON DELETE CASCADE,
    refresh_token VARCHAR(200) UNIQUE NOT NULL
);

INSERT INTO title_status (status_name)
VALUES
    ('Watching'),
    ('Completed'),
    ('On-Hold'),
    ('Dropped'),
    ('Plan to Watch');
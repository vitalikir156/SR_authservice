CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uname TEXT NOT NULL,
    taskread boolean,
    taskwrite boolean,
    userread boolean,
    userwrite boolean,
    password TEXT NOT NULL
);
CREATE TABLE tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);


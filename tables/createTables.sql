CREATE TABLE IF NOT EXISTS task(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255),
    "text" TEXT,
    author VARCHAR(255),
    urgent BOOLEAN
);
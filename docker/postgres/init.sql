CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(7) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
    expires_at TIMESTAMP WITH TIME ZONE
);

INSERT INTO
    links (id, short_code, original_url)
VALUES
    (100000, '000Q0u', 'https://example.com');
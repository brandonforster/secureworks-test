CREATE TABLE ip
(
    id BLOB PRIMARY KEY,
    created_at DATETIME,
    updated_at DATETIME,
    response_code TEXT,
    ip_address TEXT
)

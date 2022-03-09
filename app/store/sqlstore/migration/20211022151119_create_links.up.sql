CREATE TABLE IF NOT EXISTS links (
    id bigserial NOT NULL PRIMARY KEY, 
    link text NOT NULL UNIQUE, 
    code text NOT NULL, 
    userid text NOT NULL,
    deleted_at timestamp DEFAULT NULL
);
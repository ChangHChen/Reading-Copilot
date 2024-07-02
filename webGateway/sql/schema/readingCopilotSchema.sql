CREATE DATABASE IF NOT EXISTS readingcopilot CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE readingcopilot;


CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions(expiry);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);
ALTER TABLE users ADD CONSTRAINT users_uc_username UNIQUE (username);

CREATE TABLE IF NOT EXISTS gutendex_cache (
    cache_key VARCHAR(255) PRIMARY KEY,
    cache_value JSON,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
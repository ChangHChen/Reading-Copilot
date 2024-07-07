CREATE DATABASE IF NOT EXISTS reading_copilot CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE reading_copilot;


CREATE TABLE IF NOT EXISTS sessions (
    token CHAR(43) PRIMARY KEY,
    data BLOB NOT NULL,
    expiry TIMESTAMP(6) NOT NULL,
    INDEX sessions_expiry_idx (expiry)
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS gutendex_cache (
    cache_key VARCHAR(255) PRIMARY KEY,
    cache_value JSON,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reading_progress (
    id INT AUTO_INCREMENT,
    user_id INT,
    book_id INT,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    page INT,
    PRIMARY KEY (id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE (user_id, book_id),
    INDEX user_book_idx (user_id, book_id),
    INDEX user_idx(user_id)
);
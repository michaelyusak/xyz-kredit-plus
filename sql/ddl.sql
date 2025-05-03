DROP TABLE IF EXISTS consuments;
DROP TABLE IF EXISTS accounts;

CREATE TABLE accounts (
    account_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT DEFAULT NULL,
    INDEX idx_account_email (email)
);

CREATE TABLE consumers (
    consumer_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    account_id BIGINT NOT NULL,
    identity_number VARCHAR(100) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    legal_name VARCHAR(255) NOT NULL,
    place_of_birth VARCHAR(100) NOT NULL,
    date_of_birth VARCHAR(50) NOT NULL,
    salary BIGINT NOT NULL,
    identity_card_photo_url VARCHAR(2083) NOT NULL,
    selfie_photo_url VARCHAR(2083) NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    deleted_at BIGINT DEFAULT NULL,
    INDEX idx_consumer_identity_number (identity_number),
    INDEX idx_consumer_full_name (full_name),
    INDEX idx_consumer_account_id (account_id)
);

CREATE TABLE refresh_tokens (
    refresh_token_id BIGINT PRIMARY KEY AUTO_INCREMENT,
    refresh_token VARCHAR(500) NOT NULL DEFAULT '',
    account_id BIGINT NOT NULL,
    expired_at BIGINT NOT NULL,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL,
    INDEX idx_refresh_token (refresh_token)
);

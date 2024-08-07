CREATE TABLE IF NOT EXISTS Users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS Verifications (
    email VARCHAR(255) PRIMARY KEY,
    verification_code VARCHAR(255) NOT NULL,
    verified BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO
    Users (id, username, email, password_hash)
VALUES
    (
        '23e3b6f5-6785-42c6-a7f5-d8cecf04a6b9',
        'test',
        'test@test.de',
        '$2a$12$mGYv8a1151X6gMXRnhldoeptpSWreQqZGM94NgGxNsYHbrm0HQbuK'
    );
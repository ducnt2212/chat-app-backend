CREATE DATABASE chat_app

USE chat_app

CREATE TABLE users
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    hashed_password VARCHAR(MAX) NOT NULL,
    create_at DATETIME DEFAULT GETUTCDATE(),
);

CREATE TABLE rooms
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    name VARCHAR(MAX) NOT NULL,
    is_private BIT NOT NULL DEFAULT 0,
    created_by INT NOT NULL REFERENCES users(id),
    created_at DATETIME NOT NULL DEFAULT GETUTCDATE()
);

CREATE TABLE messages
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    sender_id INT NOT NULL REFERENCES users(id),
    room_id INT NOT NULL REFERENCES rooms(id),
    content VARCHAR(MAX) NOT NULL,
    created_at DATETIME DEFAULT GETUTCDATE()
);
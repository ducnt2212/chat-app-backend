CREATE DATABASE chat_app

USE chat_app

CREATE TABLE users
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) UNIQUE NOT NULL,
    hashed_password VARCHAR(MAX) NOT NULL,
    create_at DATETIME DEFAULT GETUTCDATE(),
);

CREATE TABLE messages
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    sender_id INT NOT NULL REFERENCES users(id),
    receiver_id INT NOT NULL REFERENCES users(id),
    content VARCHAR(MAX) NOT NULL,
    created_at DATETIME DEFAULT GETUTCDATE()
);
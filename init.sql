CREATE DATABASE ticketing;
USE ticketing;

CREATE TABLE events (
    id INT PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name VARCHAR(191) NOT NULL UNIQUE,
    description TEXT,
    start_date DATETIME NOT NULL,
    end_date DATETIME NOT NULL,
    location VARCHAR(191) NOT NULL,
    capacity INT NOT NULL,
    price FLOAT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active'
);


CREATE TABLE users (
    id INT PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    name VARCHAR(191) NOT NULL,
    email VARCHAR(191) NOT NULL UNIQUE,
    password VARCHAR(191) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user'
);


CREATE TABLE tickets (
    id INT PRIMARY KEY AUTO_INCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME,
    event_id INT NOT NULL,
    user_id INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    price FLOAT NOT NULL,
    purchased_at DATETIME NOT NULL,
    FOREIGN KEY (event_id) REFERENCES events(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);
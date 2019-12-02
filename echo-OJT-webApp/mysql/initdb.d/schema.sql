CREATE TABLE sample_app.users
(
    id    INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name  VARCHAR(45)  NOT NULL,
    email VARCHAR(255) NOT NULL
);

CREATE DATABASE sample_app_test;

CREATE TABLE sample_app_test.users
(
    id    INT          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name  VARCHAR(45)  NOT NULL,
    email VARCHAR(255) NOT NULL
);

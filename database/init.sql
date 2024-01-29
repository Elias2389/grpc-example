DROP TABLE IF EXISTS students;

CREATE TABLE students
(
    id   VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    age  INTEGER      NOT NULL
);

DROP TABLE IF EXISTS tests;

CREATE TABLE tests
(
    id   VARCHAR(32) PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

DROP TABLE IF EXISTS tests;

CREATE TABLE questions
(
    id       varchar(32) PRIMARY KEY,
    test_id  varchar(32)  NOT NULL,
    question varchar(255) NOT NULL,
    answer   varchar(255) NOT NULL,
    FOREIGN KEY (test_id) REFERENCES tests (id)
);
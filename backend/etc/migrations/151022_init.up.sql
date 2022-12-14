CREATE TABLE users (
    id SERIAL NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL
);

CREATE TABLE reports (
    id SERIAL NOT NULL UNIQUE,
    header VARCHAR(255) NOT NULL,
    short_body VARCHAR(255),
    edited TIMESTAMP WITH TIME ZONE
);

CREATE TABLE reports_body (
    id INT REFERENCES reports(id) ON DELETE CASCADE NOT NULL,
    body TEXT
);

CREATE TABLE users_reports (
    id SERIAL NOT NULL UNIQUE,
    users_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    reports_id INT REFERENCES reports(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE labels (
    id SERIAL NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE labels_reports (
    id SERIAL NOT NULL UNIQUE,
    labels_id INT REFERENCES labels(id) ON DELETE CASCADE NOT NULL,
    reports_id INT REFERENCES reports(id) ON DELETE CASCADE NOT NULL
);

CREATE TABLE users_labels (
    id SERIAL NOT NULL UNIQUE,
    users_id INT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    labels_id INT REFERENCES labels(id) ON DELETE CASCADE NOT NULL
);
CREATE KEYSPACE taskdb WITH replication = {'class': 'SimpleStrategy', 'replication_factor' : 1};
USE taskdb;
CREATE TABLE users (
    id UUID,
    username varchar,
    firstname varchar,
    lastname varchar,
    createdat text,
    updatedat text,
    deletedat text,
    status int,
    PRIMARY KEY(username)
);
CREATE INDEX ON users(status);

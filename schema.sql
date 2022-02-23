CREATE TABLE transfers(
    transaction VARCHAR(66) PRIMARY KEY,
    sender VARCHAR(42) NOT NULL,
    receiver VARCHAR(42) NOT NULL,
    token_id INT NOT NULL,
    read BOOLEAN NOT NULL
);
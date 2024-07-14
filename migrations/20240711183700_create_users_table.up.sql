CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    public_identifier VARCHAR(255) NOT NULL,
    firebase_uid      VARCHAR(255) NOT NULL UNIQUE
);

--bun:split

CREATE INDEX users_firebase_uid ON users(firebase_uid);

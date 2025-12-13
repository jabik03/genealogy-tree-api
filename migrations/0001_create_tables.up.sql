CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    email         VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS trees
(
    id         SERIAL PRIMARY KEY,
    owner_id   INTEGER   NOT NULL,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS persons
(
    id         SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name  VARCHAR(100) NOT NULL,
    birth_date DATE,
    death_date DATE,
    is_male    BOOLEAN      NOT NULL,
    biography  TEXT,
    tree_id    INTEGER      NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW(),
--     photo_url  VARCHAR(255),
--     metadata   JSONB
    FOREIGN KEY (tree_id) REFERENCES trees (id) ON DELETE CASCADE

);

CREATE TABLE IF NOT EXISTS relationships
(
    id                SERIAL PRIMARY KEY,
    parent_id         INTEGER NOT NULL REFERENCES persons (id) ON DELETE CASCADE,
    child_id          INTEGER NOT NULL REFERENCES persons (id) ON DELETE CASCADE,
    relationship_type VARCHAR(20),
    UNIQUE (parent_id, child_id)
);

CREATE INDEX idx_parent ON relationships (parent_id);
CREATE INDEX idx_child ON relationships (child_id);

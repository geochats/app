BEGIN;
CREATE TABLE points (
    chat_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL DEFAULT '',
    text TEXT NOT NULL DEFAULT '',
    latitude NUMERIC(16,12)  NOT NULL DEFAULT 0,
    longitude NUMERIC(16,12)  NOT NULL DEFAULT 0,
    members_count INT  NOT NULL DEFAULT 0,
    is_published BOOL  NOT NULL DEFAULT false,
    is_single BOOL NOT NULL DEFAULT false,
    PRIMARY KEY (chat_id)
);
COMMIT;

CREATE TABLE IF NOT EXISTS tracers (
    id VARCHAR PRIMARY KEY,
    ip VARCHAR UNIQUE NOT NULL,
    total_blinks INTEGER,
    last_blink_at TIMESTAMP NULL,
    updated_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL
);

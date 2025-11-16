CREATE TABLE IF NOT EXISTS tracers (
    id VARCHAR PRIMARY KEY,
    nickname VARCHAR UNIQUE NOT NULL,
    total_blinks INTEGER,
    last_blink_at TIMESTAMPTZ NULL,
    updated_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL
);

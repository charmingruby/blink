CREATE TABLE IF NOT EXISTS blinks (
    id VARCHAR PRIMARY KEY,
    tracer_id VARCHAR UNIQUE NOT NULL,
    result VARCHAR NOT NULL,
    status VARCHAR NOT NULL,
    updated_at TIMESTAMP NULL,
    created_at TIMESTAMP NOT NULL
);

ALTER TABLE blinks ADD FOREIGN KEY (tracer_id) REFERENCES tracers(id);

CREATE INDEX idx_blinks_tracer_id ON blinks(tracer_id);

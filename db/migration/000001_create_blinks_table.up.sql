CREATE TABLE IF NOT EXISTS blinks (
    id VARCHAR PRIMARY KEY,
    tracer_id VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE blinks ADD FOREIGN KEY (tracer_id) REFERENCES tracers(id);

CREATE INDEX idx_blinks_tracer_id ON blinks(tracer_id);

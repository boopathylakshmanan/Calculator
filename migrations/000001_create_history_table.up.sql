CREATE TABLE history (
    id BIGSERIAL PRIMARY KEY,
    expression TEXT NOT NULL,
    result TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

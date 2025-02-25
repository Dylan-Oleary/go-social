BEGIN;

CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL,
    content text NOT NULL,
    title text NOT NULL,
    tags VARCHAR(100) [],
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

COMMIT;
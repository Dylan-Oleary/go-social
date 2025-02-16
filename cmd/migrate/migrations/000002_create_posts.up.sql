BEGIN;

CREATE TABLE IF NOT EXISTS posts (
    id bigserial PRIMARY KEY,
    user_id bigserial NOT NULL,
    content text NOT NULL,
    title text NOT NULL,
    tags text[] NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at timestamp(0) WITH TIME ZONE,

    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id)
);

COMMIT;
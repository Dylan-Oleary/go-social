BEGIN;

CREATE TABLE IF NOT EXISTS followers (
    user_id bigserial NOT NULL,
    follower_id bigserial NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, follower_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) references users(id) ON DELETE CASCADE,
    CONSTRAINT fk_follower_id FOREIGN KEY (follower_id) references users(id) ON DELETE CASCADE
);

COMMIT;
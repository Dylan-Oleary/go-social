BEGIN;

CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Comments
CREATE INDEX IF NOT EXISTS idx_comments_content ON comments USING gin (content gin_trgm_ops);
CREATE INDEX IF NOT EXISTS isx_comments_post_id ON comments (post_id);

-- Posts
CREATE INDEX IF NOT EXISTS idx_posts_title ON posts USING gin (title gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_posts_tags ON posts USING gin (tags);
CREATE INDEX IF NOT EXISTS idx_posts_user_id on posts (user_id);

-- Users
CREATE INDEX IF NOT EXISTS idx_users_username ON users (username);

COMMIT;
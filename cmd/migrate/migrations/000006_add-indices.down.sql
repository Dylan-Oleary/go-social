BEGIN;

-- Comments
DROP INDEX IF EXISTS idx_comments_content;
DROP INDEX IF EXISTS isx_comments_post_id;

-- Posts
DROP INDEX IF EXISTS idx_posts_title;
DROP INDEX IF EXISTS idx_posts_tags;
DROP INDEX IF EXISTS idx_posts_user_id;

-- Users
DROP INDEX IF EXISTS idx_users_username;

COMMIT;
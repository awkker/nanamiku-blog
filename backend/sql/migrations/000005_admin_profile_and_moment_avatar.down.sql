ALTER TABLE moments
    DROP COLUMN IF EXISTS author_avatar_url;

ALTER TABLE admin_users
    DROP COLUMN IF EXISTS avatar_url,
    DROP COLUMN IF EXISTS display_name;

ALTER TABLE admin_users
    ADD COLUMN display_name text NOT NULL DEFAULT '',
    ADD COLUMN avatar_url text NOT NULL DEFAULT '/picture/author.jpg';

UPDATE admin_users
SET display_name = username
WHERE btrim(display_name) = '';

UPDATE admin_users
SET avatar_url = '/picture/author.jpg'
WHERE btrim(avatar_url) = '';

ALTER TABLE moments
    ADD COLUMN author_avatar_url text NOT NULL DEFAULT '/picture/author.jpg';

UPDATE moments
SET author_avatar_url = '/picture/author.jpg'
WHERE btrim(author_avatar_url) = '';

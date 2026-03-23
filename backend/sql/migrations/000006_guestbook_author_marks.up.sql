CREATE TABLE guestbook_author_marks (
    message_id       uuid PRIMARY KEY REFERENCES guestbook_messages(id) ON DELETE CASCADE,
    admin_id         uuid NOT NULL REFERENCES admin_users(id),
    admin_avatar_url text NOT NULL DEFAULT '',
    created_at       timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX idx_guestbook_author_marks_admin ON guestbook_author_marks(admin_id);

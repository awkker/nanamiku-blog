CREATE TABLE friend_link_applications (
    id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    site_name       text NOT NULL,
    site_url        text NOT NULL,
    avatar_url      text NOT NULL DEFAULT '',
    description     text NOT NULL DEFAULT '',
    contact_email   text NOT NULL DEFAULT '',
    contact_note    text NOT NULL DEFAULT '',
    status          friend_link_status NOT NULL DEFAULT 'pending',
    created_at      timestamptz NOT NULL DEFAULT now(),
    reviewed_at     timestamptz,
    review_note     text NOT NULL DEFAULT '',
    reviewed_by     uuid REFERENCES admin_users(id) ON DELETE SET NULL
);

CREATE INDEX idx_friend_link_applications_status ON friend_link_applications(status, created_at DESC);

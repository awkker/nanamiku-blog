CREATE TABLE site_settings (
    key         text PRIMARY KEY,
    value       jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_at  timestamptz NOT NULL DEFAULT now(),
    updated_at  timestamptz NOT NULL DEFAULT now(),
    updated_by  uuid REFERENCES admin_users(id) ON DELETE SET NULL
);

-- name: GetSiteSetting :one
SELECT key, value, created_at, updated_at, updated_by
FROM site_settings
WHERE key = $1;

-- name: UpsertSiteSetting :one
INSERT INTO site_settings (key, value, updated_by)
VALUES ($1, $2, $3)
ON CONFLICT (key) DO UPDATE
SET value = EXCLUDED.value,
    updated_by = EXCLUDED.updated_by,
    updated_at = now()
RETURNING key, value, created_at, updated_at, updated_by;

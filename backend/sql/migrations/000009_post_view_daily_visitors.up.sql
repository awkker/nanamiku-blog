CREATE TABLE post_view_daily_visitors (
    post_id     uuid NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    day         date NOT NULL,
    visitor_id  uuid NOT NULL REFERENCES visitors(id) ON DELETE CASCADE,
    created_at  timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (post_id, day, visitor_id)
);

CREATE INDEX idx_post_view_daily_visitors_day ON post_view_daily_visitors(day DESC);

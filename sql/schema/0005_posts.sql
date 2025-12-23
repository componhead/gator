-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(), 
  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  published_at TIMESTAMP NOT NULL DEFAULT now(), 
  feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;

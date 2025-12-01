-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  updated_at TIMESTAMP NOT NULL DEFAULT now(), 
  feed_id UUID NOT NULL references feeds(id) ON DELETE CASCADE,
  user_id UUID NOT NULL references users(id) ON DELETE CASCADE,
  UNIQUE(user_id,feed_id)
);

-- +goose Down
DROP TABLE feed_follows; 


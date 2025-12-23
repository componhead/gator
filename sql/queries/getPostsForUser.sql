-- name: GetPostsForUser :many
SELECT *
FROM posts p
  INNER JOIN feeds f ON p.feed_id = f.id
WHERE
  f.user_id = $1
ORDER BY p.created_at DESC
LIMIT $2;

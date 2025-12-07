-- name: GetNextFeedToFetch :one
SELECT f.id, f.name as feedName, f.url as feedUrl
FROM feed_follows ff
INNER JOIN feeds f ON ff.feed_id = f.id
INNER JOIN users u ON ff.user_id = u.id
WHERE u.ID = $1
ORDER BY f.last_fetched_at ASC NULLS FIRST
LIMIT 1;

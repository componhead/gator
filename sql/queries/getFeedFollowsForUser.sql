-- name: GetFeedsFollowForUser :many
SELECT ff.id, f.name as feedName, u.name as userName
FROM feed_follows ff
INNER JOIN feeds f ON ff.feed_id = f.id
INNER JOIN users u ON ff.user_id = u.id
WHERE ff.user_id = $1;

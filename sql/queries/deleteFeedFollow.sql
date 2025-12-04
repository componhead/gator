-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows where user_id = $1 and feed_id = $2;

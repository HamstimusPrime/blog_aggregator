-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows ff
WHERE ff.user_id = $1 AND ff.feed_id = (SELECT id FROM feeds WHERE url = $2);
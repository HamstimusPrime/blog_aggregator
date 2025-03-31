-- name: GetFeeds :many
SELECT feeds.name AS feedName, feeds.url, users.name AS username
FROM users
JOIN feeds ON feeds.user_id = users.id;
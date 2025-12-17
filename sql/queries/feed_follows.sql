-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT inserted_feed_follow.*, users.name AS username, feeds.name AS feed_name
FROM inserted_feed_follow 
INNER JOIN users ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id;

-- name: FindFeedID :one
SELECT id FROM feeds WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*, users.name AS username, feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON feed_follows.user_id = users.id
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE users.name = $1;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows 
WHERE user_id IN (SELECT id FROM users WHERE users.name = $1) 
AND feed_id IN  (SELECT id FROM feeds WHERE feeds.url = $2);

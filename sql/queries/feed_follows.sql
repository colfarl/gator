-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows(id, created_at, updated_at, feed_id, user_id)
    VALUES (
	$1,
	$2,
	$3,
	$4,
	$5
    )
    RETURNING *
)

SELECT ff.*, u.name AS user_name, f.name AS feed_name
FROM inserted_feed_follow ff
    JOIN users u ON u.id = ff.user_id
    JOIN feeds f ON f.id = ff.feed_id;

-- name: GetFeedFollowsForUser :many
SELECT ff.*, u.name AS user_name, f.name AS feed_name
FROM feed_follows ff
    JOIN users u ON u.id = ff.user_id
    JOIN feeds f ON f.id = ff.feed_id
WHERE ff.user_id = $1;


-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

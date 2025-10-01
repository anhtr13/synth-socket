-- name: CreateFriendRequest :one
INSERT INTO
	friend_requests (sender_id, receiver_id, created_at)
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: DeleteFriendRequest :exec
DELETE FROM friend_requests
WHERE
	sender_id = $1
	AND receiver_id = $2;

-- name: AcceptFriendRequest :one
UPDATE friend_requests
SET
	accepted = TRUE
WHERE
	request_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: RejectFriendRequest :one
UPDATE friend_requests
SET
	accepted = FALSE
WHERE
	request_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: GetUserFriendRequests :many
SELECT
	fr.*,
	u.user_name sender_name,
	u.profile_image sender_image
FROM
	(
		SELECT
			*
		FROM
			friend_requests
		WHERE
			receiver_id = $1
		LIMIT
			$2
		OFFSET
			$3
	) fr
	LEFT JOIN users u ON fr.sender_id = u.user_id;

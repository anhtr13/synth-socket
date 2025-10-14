-- name: CreateMessage :one
INSERT INTO
	messages (room_id, sender_id, text, media_url)
VALUES
	($1, $2, $3, $4)
RETURNING
	*;

-- name: DeleteMessageById :one
DELETE FROM messages
WHERE
	message_id = $1
RETURNING
	*;

-- name: GetMesssagesByRoomId :many
SELECT
	*
FROM
	messages
WHERE
	room_id = $1
ORDER BY
	created_at DESC
LIMIT
	$2
OFFSET
	$3;

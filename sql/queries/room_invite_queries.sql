-- name: CreateRoomInvite :one
INSERT INTO
	room_invites (room_id, sender_id, receiver_id)
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: DeleteRoomInvite :exec
DELETE FROM room_invites
WHERE
	sender_id = $1
	AND receiver_id = $2;

-- name: AcceptRoomInvite :one
UPDATE room_invites
SET
	accepted = TRUE
WHERE
	invite_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: RejectRoomInvite :one
DELETE FROM room_invites
WHERE
	invite_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: GetRoomInvitesByRoomId :many
SELECT
	rr.*,
	u.user_name receiver_name,
	u.profile_image receiver_image
FROM
	(
		SELECT
			*
		FROM
			room_invites
		WHERE
			room_id = $1
		LIMIT
			$2
		OFFSET
			$3
	) rr
	LEFT JOIN users u ON rr.receiver_id = u.user_id;

-- name: GetRoomInvitesByReceiverId :many
SELECT
	rr.*,
	u.user_name sender_name,
	u.profile_image sender_image
FROM
	(
		SELECT
			*
		FROM
			room_invites
		WHERE
			receiver_id = $1
		LIMIT
			$2
		OFFSET
			$3
	) rr
	LEFT JOIN users u ON rr.sender_id = u.user_id;

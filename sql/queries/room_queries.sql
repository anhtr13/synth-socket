-- name: CreateRoom :one
INSERT INTO
	rooms (room_name, room_picture, created_by)
VALUES
	($1, $2, $3)
RETURNING
	*;

-- name: UpdateRoomLastMessage :exec
UPDATE rooms
SET
	last_message = $1,
	updated_at = $2
WHERE
	room_id = $3;

-- name: FindRoomById :one
SELECT
	*
FROM
	rooms
WHERE
	room_id = $1;

-- name: FindRoomByCreatorAndName :one
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
	AND room_name = $2;

-- name: GetRoomsByCreator :many
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
ORDER BY
	room_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsByCreatorAndName :many
SELECT
	*
FROM
	rooms
WHERE
	created_by = $1
	AND room_name LIKE ('%' || sqlc.arg (room_name) || '%')
ORDER BY
	room_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsDataByMemberIdAndRoomName :many
SELECT
	r.*,
	rm.joined_at
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			member_id = $1
	) rm
	LEFT JOIN rooms r ON rm.room_id = r.room_id
WHERE
	r.room_name LIKE ('%' || sqlc.arg (name) || '%')
ORDER BY
	r.room_name
LIMIT
	$2
OFFSET
	$3;

-- name: GetRoomsDataByMemberId :many
SELECT
	r.*,
	rm.joined_at
FROM
	(
		SELECT
			*
		FROM
			room_members
		WHERE
			member_id = $1
	) rm
	LEFT JOIN rooms r ON rm.room_id = r.room_id
ORDER BY
	r.updated_at DESC
LIMIT
	$2
OFFSET
	$3;

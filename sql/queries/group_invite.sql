-- name: CreateGroupInvite :one
INSERT INTO
	group_invites (group_id, sender_id, receiver_id, created_at)
VALUES
	($1, $2, $3, $4)
RETURNING
	*;

-- name: DeleteGroupInvite :exec
DELETE FROM group_invites
WHERE
	sender_id = $1
	AND receiver_id = $2;

-- name: AcceptGroupInvite :one
UPDATE group_invites
SET
	accepted = TRUE
WHERE
	invite_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: RejectGroupInvite :one
UPDATE group_invites
SET
	accepted = FALSE
WHERE
	invite_id = $1
	AND receiver_id = $2
RETURNING
	*;

-- name: GetGroupInvites :many
SELECT
	rr.*,
	u.user_name receiver_name,
	u.profile_image receiver_image
FROM
	(
		SELECT
			*
		FROM
			group_invites
		WHERE
			group_id = $1
		LIMIT
			$2
		OFFSET
			$3
	) rr
	LEFT JOIN users u ON rr.receiver_id = u.user_id;

-- name: GetUserGroupInvites :many
SELECT
	rr.*,
	u.user_name sender_name,
	u.profile_image sender_image
FROM
	(
		SELECT
			*
		FROM
			group_invites
		WHERE
			receiver_id = $1
		LIMIT
			$2
		OFFSET
			$3
	) rr
	LEFT JOIN users u ON rr.sender_id = u.user_id;

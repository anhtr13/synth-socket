-- name: CreateNotification :one
INSERT INTO
	notifications (message, user_id, type, id_ref)
VALUES
	($1, $2, $3, $4)
RETURNING
	*;

-- name: GetUnSeenNotificationsByUserId :many
SELECT
	*
FROM
	notifications
WHERE
	user_id = $1
	AND seen = FALSE
LIMIT
	$2
OFFSET
	$3;

-- name: DeleteNotification :exec
DELETE FROM notifications
WHERE
	notification_id = $1;

-- name: SeenNotification :exec
UPDATE notifications
SET
	seen = TRUE
WHERE
	notification_id = $1
	AND user_id = $2;

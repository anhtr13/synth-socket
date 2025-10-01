-- +goose Up
CREATE TYPE notification_type AS enum('friend_request', 'group_request');

CREATE TABLE notifications (
	notification_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL REFERENCES users (user_id),
	message varchar(1024) NOT NULL,
	type notification_type NOT NULL,
	id_ref uuid,
	seen bool DEFAULT FALSE,
	created_at timestamp
);

-- +goose Down
DROP TABLE notifications;

DROP TYPE notification_type;

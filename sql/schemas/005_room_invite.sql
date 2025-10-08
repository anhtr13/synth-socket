-- +goose Up
CREATE TABLE room_invites (
	invite_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	room_id uuid NOT NULL REFERENCES rooms (room_id),
	sender_id uuid NOT NULL REFERENCES users (user_id),
	receiver_id uuid NOT NULL REFERENCES users (user_id),
	accepted bool DEFAULT FALSE,
	created_at timestamp without time zone DEFAULT now(),
	UNIQUE (room_id, sender_id, receiver_id)
);

CREATE INDEX idx_room_invites_sender_id ON room_invites (sender_id);

CREATE INDEX idx_room_invites_receiver_id ON room_invites (receiver_id);

-- +goose Down
DROP INDEX idx_room_invites_sender_id;

DROP INDEX idx_room_invites_receiver_id;

DROP TABLE room_invites;

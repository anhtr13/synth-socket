-- +goose Up
CREATE TABLE room_members (
	room_id uuid NOT NULL REFERENCES rooms (room_id),
	member_id uuid NOT NULL REFERENCES users (user_id),
	joined_at timestamp without time zone DEFAULT now(),
	UNIQUE (room_id, member_id)
);

CREATE INDEX idx_room_members_member_id ON room_members (member_id);

-- +goose Down
DROP INDEX idx_room_members_member_id;

DROP TABLE room_members;

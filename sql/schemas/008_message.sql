-- +goose Up
CREATE TABLE messages (
	message_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	room_id uuid NOT NULL REFERENCES rooms (room_id),
	sender_id uuid NOT NULL REFERENCES users (user_id),
	text varchar(2048) NOT NULL,
	media_url varchar(255),
	created_at timestamp without time zone DEFAULT now()
);

CREATE INDEX idx_messages_created_at ON messages (created_at);

CREATE INDEX idx_messages_room_id ON messages (room_id);

-- +goose Down
DROP INDEX idx_messages_created_at;

DROP INDEX idx_messages_room_id;

DROP TABLE messages;

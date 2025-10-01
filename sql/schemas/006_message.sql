-- +goose Up
CREATE TYPE message_status AS enum('sent', 'delivered', 'read');

CREATE TABLE messages (
	"message_id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	"group_id" uuid NOT NULL REFERENCES groups (group_id),
	"sender_id" uuid NOT NULL REFERENCES users (user_id),
	"text" varchar(2048) NOT NULL,
	"media_url" varchar(255) NOT NULL,
	"status" message_status DEFAULT 'sent',
	"created_at" timestamp NOT NULL
);

CREATE INDEX idx_messages_created_at ON messages (created_at);

CREATE INDEX idx_messages_group_id ON messages (group_id);

-- +goose Down
DROP INDEX idx_messages_created_at;

DROP INDEX idx_messages_group_id;

DROP TABLE messages;

DROP TYPE message_status;

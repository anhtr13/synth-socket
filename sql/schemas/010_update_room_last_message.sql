-- +goose Up
ALTER TABLE rooms
ADD COLUMN last_message uuid REFERENCES messages (message_id);

-- +goose Down
ALTER TABLE rooms
DROP COLUMN last_message;

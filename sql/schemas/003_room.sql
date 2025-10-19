-- +goose Up
CREATE TABLE rooms (
	room_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	room_name varchar(255) NOT NULL,
	room_picture varchar(255),
	created_by uuid NOT NULL REFERENCES users (user_id),
	updated_at timestamp without time zone DEFAULT now(),
	created_at timestamp without time zone DEFAULT now(),
	UNIQUE (room_name, created_by)
);

CREATE INDEX idx_rooms_room_name ON rooms (room_name);

-- +goose Down
DROP INDEX idx_rooms_room_name;

DROP TABLE rooms;

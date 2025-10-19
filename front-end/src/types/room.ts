export type Room = {
	room_id: string;
	room_name: string;
	room_picture: string | null;
	created_by: string;
	created_at: string;
	joined_at: string;
	seen_last_message?: boolean;
};

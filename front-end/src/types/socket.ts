export type SEvent = "message" | "error" | "notification" | "room_io";

export type SMessage = {
	sender_id: string;
	receiver_id: string;
	text: string;
	media_url: string;
};

export type SNotification = {
	notification_id: string;
	user_id: string;
	message: string;
	type: "friend_request" | "room_invite";
	id_ref: string;
	seen: boolean;
	created_at: string;
};

export type RoomIo = {
	user_id: string;
	room_id: string;
	type: "room_in" | "room_out";
};

export type SError = {
	error: string;
};

export type SPayload = {
	event: SEvent;
	data: SMessage | SError | SNotification;
};

export type SEvent = "message" | "notification";

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

export type SPayload = {
	event: SEvent;
	data: any; // Message | Notification
};

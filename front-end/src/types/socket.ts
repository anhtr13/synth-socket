export type SEvent = "message" | "notification" | "room_io" | "friend_io";

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

export type SRoomIO = {
	user_id: string;
	room_id: string;
	type: "room_in" | "room_out";
};

export type SFriendIO = {
	user1_id: string;
	user2_id: string;
	type: "friend_in" | "friend_out";
};

export type SPayload = {
	event: SEvent;
	data: any; // Message | Notification | RoomIO
};

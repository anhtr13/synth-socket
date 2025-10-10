export type UserInfo = {
	user_id: string;
	user_name: string;
	profile_image: string | null;
	is_friend?: boolean;
};

export type UserNoti = {
	notification_id: string;
	user_id: string;
	message: string;
	type: "friend_request" | "room_invite";
	id_ref: string;
	seen: boolean;
	created_at: string;
};

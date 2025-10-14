export type UserInfo = {
	user_id: string;
	user_name: string;
	profile_image: string | null;
	is_friend?: boolean;
};

export type RoomMemberInfo = {
	user_id: string;
	user_name: string;
	profile_image: string | null;
	joined_at: string;
};

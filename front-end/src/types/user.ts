export type UserInfo = {
	user_id: string;
	user_name: string;
	profile_image: string | null;
	is_friend?: boolean;
	last_active?: string;
};

export type RoomMemberInfo = {
	user_id: string;
	user_name: string;
	profile_image: string | null;
	joined_at?: string;
};

export type FriendshipInfo = {
	friendship_id: string;
	user1_id: string;
	user2_id: string;
	created_at: string;
};

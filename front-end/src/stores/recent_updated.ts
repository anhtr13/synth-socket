import { defineStore } from "pinia";
import { ref } from "vue";
import type { Room } from "@/types/room";
import type { UserInfo } from "@/types/user";
import type { SOnlineStatus } from "@/types/socket";

export const useRecentUpdatedStore = defineStore("recentUpdatedStore", () => {
	// least recent -> most recent
	const roomSet = ref<Map<string, Room>>(new Map());
	const friendSet = ref<Map<string, UserInfo>>(new Map());

	const updateRoomSet = (list: Room[]) => {
		let newSet = new Map<string, Room>();
		for (let i = list.length - 1; i >= 0; i--) {
			let room = list[i]!;
			newSet.set(room.room_id, room);
		}
		roomSet.value = newSet;
	};
	const updateRoomNewMessage = (room_id: string) => {
		let room = roomSet.value.get(room_id);
		if (room) {
			roomSet.value.delete(room_id);
			room.seen_last_message = false;
			roomSet.value.set(room_id, room);
		}
	};
	const updateRoomSeenMessage = (room_id: string) => {
		let room = roomSet.value.get(room_id);
		if (room) {
			room.seen_last_message = true;
		}
	};

	const updateFriendSet = (list: UserInfo[]) => {
		let newSet = new Map<string, UserInfo>();
		for (let i = list.length - 1; i >= 0; i--) {
			let fr = list[i]!;
			newSet.set(fr.user_id, fr);
		}
		friendSet.value = newSet;
	};

	const updateFriendOnlineStatus = (status: SOnlineStatus) => {
		for (let fr_id in status) {
			let friend = friendSet.value.get(fr_id);
			if (friend) {
				friend.online_status = status[fr_id];
			}
		}
	};

	return {
		roomSet,
		friendSet,
		updateRoomSet,
		updateRoomSeenMessage,
		updateRoomNewMessage,
		updateFriendSet,
    updateFriendOnlineStatus
	};
});

import { defineStore } from "pinia";
import { ref } from "vue";
import type { RoomMemberInfo } from "@/types/user";
import type { SMessage } from "@/types/socket";

export const useRoomDataStore = defineStore("roomDataStore", () => {
	const current_id = ref<string>("");
	const members = ref<Map<string, RoomMemberInfo>>(new Map());
	const messages = ref<SMessage[]>([]);

	function updateMemberSet(data: RoomMemberInfo[]) {
		if (current_id.value === "") return;
		const map = new Map<string, RoomMemberInfo>();
		for (let mem of data) {
			map.set(mem.user_id, mem);
		}
		members.value = map;
	}
	function updateMessageList(data: SMessage[]) {
		if (current_id.value === "") return;
		let msgs: SMessage[] = [];
		for (let i = data.length - 1; i >= 0; i--) {
			msgs.push(data[i]!);
		}
		messages.value = msgs;
	}
	function addNewMessage(msg: SMessage) {
		if (current_id.value === "") return;
		messages.value.push(msg);
	}

	return { current_id, members, messages, updateMemberSet, updateMessageList, addNewMessage };
});

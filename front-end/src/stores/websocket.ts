import { defineStore } from "pinia";
import { ref } from "vue";
import { useRecentUpdatedStore } from "./recent_updated";
import { useNotificationStore } from "./notifications";
import { useRoomDataStore } from "./room_data";
import { usePersonalStore } from "./personal";
import type { SFriendIO, SMessage, SPayload, SRoomIO } from "@/types/socket";
import type { UserInfo } from "@/types/user";
import { _get } from "@/utils/fetch";
import type { Room } from "@/types/room";

export const useWebSocketStore = defineStore("websocketStore", () => {
	const ws = ref<WebSocket | null>(null);
	const isConnected = ref(false);
	const personalStore = usePersonalStore();
	const recentUpdatedStore = useRecentUpdatedStore();
	const notificationStore = useNotificationStore();
	const roomDataStore = useRoomDataStore();

	function connect(url: string) {
		ws.value = new WebSocket(url);

		ws.value.onopen = () => {
			isConnected.value = true;
			console.log("websocket connected");
		};
		ws.value.onclose = () => {
			isConnected.value = true;
			console.log("websocket disconnected");
			const cf = confirm("Websocket connection closed, reconnect?");
			if (cf) {
				connect(url);
			}
		};
		ws.value.onerror = (err) => {
			isConnected.value = true;
			alert("Something went wrong, try latter.");
			console.log("websocket error:", err);
		};
		ws.value.onmessage = async (event) => {
			const payload: SPayload = JSON.parse(event.data);
			console.log("New message", payload);
			switch (payload.event) {
				case "message":
					const msg: any = payload.data;
					recentUpdatedStore.updateRoomNewMessage(msg.receiver_id);
					if (msg.receiver_id === roomDataStore.current_id) {
						roomDataStore.addNewMessage(msg);
					}
					break;
				case "notification":
					const noti: any = payload.data;
					notificationStore.notifications.push(noti);
					break;
				case "room_io":
					const { room_id, user_id, type: r_io_type }: SRoomIO = payload.data;
					if (user_id === personalStore.info?.user_id) {
						switch (r_io_type) {
							case "room_in":
								const room: Room = await _get(`/api/v1/room/all/${room_id}`).catch((err) =>
									console.error(err),
								);
								room.seen_last_message = false;
								recentUpdatedStore.roomSet.set(room_id, room);
								break;
							case "room_out":
								recentUpdatedStore.roomSet.delete(room_id);
								break;
						}
						return;
					}
					if (room_id === roomDataStore.current_id) {
						switch (r_io_type) {
							case "room_in":
								const user: UserInfo = await _get(`/api/v1/user/${user_id}`).catch((err) =>
									console.error(err),
								);
								roomDataStore.members.set(user_id, user);
								roomDataStore.addNewMessage({
									sender_id: "server",
									text: `${user.user_name} has joined the room.`,
									media_url: "",
									receiver_id: room_id,
								});
								break;
							case "room_out":
								const member = roomDataStore.members.get(user_id);
								if (!member) return;
								roomDataStore.members.delete(user_id);
								roomDataStore.addNewMessage({
									sender_id: "server",
									text: `${member.user_name} has left the room.`,
									media_url: "",
									receiver_id: room_id,
								});
								break;
						}
					}
					recentUpdatedStore.updateRoomNewMessage(room_id);
					break;
				case "friend_io":
					const { user1_id, user2_id, type: f_io_type }: SFriendIO = payload.data;
					let friend_id = user1_id;
					if (user1_id === personalStore.info?.user_id) {
						friend_id = user2_id;
					}
					if (f_io_type === "friend_in") {
						const user: UserInfo = await _get(`/api/v1/user/${friend_id}`).catch((err) =>
							console.error(err),
						);
						recentUpdatedStore.friendSet.set(user.user_id, user);
					}
					if (f_io_type === "friend_out") {
						recentUpdatedStore.friendSet.delete(friend_id);
					}
					break;
				default:
			}
		};
	}

	function disconnect() {
		if (ws.value) {
			ws.value.close();
			ws.value = null;
		}
	}

	function sendMessage(msg: SMessage) {
		if (ws.value && isConnected.value) {
			ws.value.send(JSON.stringify(msg));
		} else {
			alert("WebSocket not connected, cannot send message.");
			console.warn("WebSocket not connected, cannot send message.");
		}
	}

	return { ws, isConnected, connect, disconnect, sendMessage };
});

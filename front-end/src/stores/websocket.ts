import { defineStore } from "pinia";
import { ref } from "vue";
import { useRecentUpdatedStore } from "./recent_updated";
import { useNotificationStore } from "./notifications";
import { useRoomDataStore } from "./room_data";
import type { SMessage, SPayload } from "@/types/socket";

export const useWebSocketStore = defineStore("websocketStore", () => {
	const ws = ref<WebSocket | null>(null);
	const isConnected = ref(false);
	const recentupdatedStore = useRecentUpdatedStore();
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
		ws.value.onmessage = (event) => {
			const payload: SPayload = JSON.parse(event.data);
			console.log("New message", payload);
			switch (payload.event) {
				case "message":
					const msg: any = payload.data;
					recentupdatedStore.updateRoomNewMessage(msg.receiver_id);
					if (msg.receiver_id === roomDataStore.current_id) {
						roomDataStore.addNewMessage(msg);
					}
					break;
				case "notification":
					const noti: any = payload.data;
					notificationStore.notifications.push(noti);
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

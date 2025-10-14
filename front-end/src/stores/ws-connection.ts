import { defineStore } from "pinia";
import { ref } from "vue";

export const useWsConnectionStore = defineStore("ws-connection", () => {
	const connection = ref<WebSocket | null>(null);
	return { connection };
});

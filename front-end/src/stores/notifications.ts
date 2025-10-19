import { defineStore } from "pinia";
import { ref } from "vue";
import type { SNotification } from "@/types/socket";

export const useNotificationStore = defineStore("notificationStore", () => {
	const notifications = ref<SNotification[]>([]);
	return { notifications };
});

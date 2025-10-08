import { defineStore } from "pinia";
import { ref } from "vue";

export type FocusingContact = "room" | "friend";
export const useContactBoxStore = defineStore("contactBoxStore", () => {
	const showMobileDropdown = ref(false);
	const focusing = ref<FocusingContact>("room");

	function toggleMobileDropdown() {
		showMobileDropdown.value = !showMobileDropdown.value;
	}

	return { focusing, showMobileDropdown, toggleMobileDropdown };
});

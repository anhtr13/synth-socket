import { defineStore } from "pinia";
import { ref } from "vue";

export type FocusingHeader = "room" | "friend";
export const useGlobalStateStore = defineStore("globalStateStore", () => {
	const showHeaderMobileDropdown = ref(false);
	const focusingHeader = ref<FocusingHeader>("room");

	function toggleHeaderMobileDropdown() {
		showHeaderMobileDropdown.value = !showHeaderMobileDropdown.value;
	}

	return {
		focusingHeader,
		showHeaderMobileDropdown,
		toggleHeaderMobileDropdown,
	};
});

import type { PersonalData } from "@/types/personal";
import { defineStore } from "pinia";
import { ref } from "vue";

export const usePersonalStore = defineStore("personalStore", () => {
	const info = ref<PersonalData | null>(null);
	return { info };
});

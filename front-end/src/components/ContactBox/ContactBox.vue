<script setup lang="ts">
import { useGlobalStateStore } from "@/stores/global_state";
import { _get } from "@/utils/fetch";
import IconClose from "@/components/icons/IconClose.vue";
import Rooms from "@/components/ContactBox/Rooms.vue";
import Friends from "@/components/ContactBox/Friends.vue";

const useGlobalStore = useGlobalStateStore();
</script>

<template>
	<div
		class="thin-scrollbar absolute top-0 left-0 h-full w-full shrink-0 overflow-auto bg-zinc-950 px-3 py-2 sm:relative sm:flex sm:w-80 sm:flex-col sm:border-r sm:border-neutral-700"
		:class="useGlobalStore.showHeaderMobileDropdown ? '' : 'hidden'">
		<Rooms v-if="useGlobalStore.focusingHeader === 'room'" />
		<Friends v-else-if="useGlobalStore.focusingHeader === 'friend'" />
	</div>
	<button
		@click="useGlobalStore.showHeaderMobileDropdown = false"
		v-if="useGlobalStore.showHeaderMobileDropdown"
		class="absolute top-3 right-3 flex h-5 w-5 items-center justify-center hover:bg-violet-500 sm:hidden">
		<IconClose class="size-4" />
	</button>
</template>

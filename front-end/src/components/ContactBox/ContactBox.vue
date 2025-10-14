<script setup lang="ts">
import { useContactBoxStore } from "@/stores/contact-box";
import RoomsBox from "@/components/ContactBox/RoomsBox.vue";
import FriendsBox from "./FriendsBox.vue";
import IconClose from "@/components/icons/IconClose.vue";

const contactBoxStore = useContactBoxStore();
</script>

<template>
	<div
		class="absolute top-0 left-0 h-full w-full shrink-0 bg-zinc-950 px-3 py-2 sm:relative sm:flex sm:w-80 sm:flex-col sm:border-r sm:border-neutral-700"
		:class="contactBoxStore.showMobileDropdown ? '' : 'hidden'">
		<RoomsBox v-if="contactBoxStore.focusing === 'room'" />
		<FriendsBox v-else-if="contactBoxStore.focusing === 'friend'" />
	</div>
	<button
		@click="contactBoxStore.showMobileDropdown = false"
		v-if="contactBoxStore.showMobileDropdown"
		class="absolute top-3 right-3 flex h-5 w-5 items-center justify-center hover:bg-violet-500 sm:hidden">
		<IconClose class="size-4" />
	</button>
</template>

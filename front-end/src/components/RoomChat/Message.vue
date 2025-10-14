<script setup lang="ts">
import type { RoomMemberInfo } from "@/types/user";
import type { SMessage } from "@/types/socket";
import IconUserAstronaut from "@/components/icons/IconUserAstronaut.vue";

const props = defineProps<{
	sender?: RoomMemberInfo;
	message: SMessage;
}>();
</script>

<template>
	<div class="mb-4 flex w-full flex-row">
		<div class="mr-3 shrink-0">
			<button v-if="!sender || !sender.profile_image">
				<IconUserAstronaut class="size-7" />
			</button>
			<button v-else>
				<img
					class="size-8 object-cover"
					:src="sender.profile_image" />
			</button>
		</div>
		<div class="flex grow flex-col">
			<span class="mb-1 text-sm font-bold text-neutral-400">
				{{ sender?.user_name || "Unknow member" }}
			</span>
			<div class="w-full max-w-3xl shrink-0 bg-neutral-800 px-3 py-2 break-all">
				{{ message.text }}
			</div>
		</div>
	</div>
</template>

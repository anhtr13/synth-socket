<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { useUserInfoStore } from "@/stores/user";
import { useWsConnectionStore } from "@/stores/ws-connection";
import { _get } from "@/utils/fetch";
import { toast } from "vue3-toastify";
import type { SMessage } from "@/types/socket";
import type { RoomMemberInfo } from "@/types/user";
import Message from "@/components/RoomChat/Message.vue";
import IconSend from "@/components/icons/IconSend.vue";

const route = useRoute();
const userInfoStore = useUserInfoStore();
const connectionStore = useWsConnectionStore();

const room_id = computed<any>(() => route.params.room_id);

const members_map = ref<{ [key: string]: RoomMemberInfo }>({});
const messages = ref<SMessage[]>([
	{
		receiver_id: "f9598d63-47c4-474a-8c46-7bb63d5361fd",
		sender_id: "4efdc06f-d80e-4f80-9278-23b9f2",
		media_url: "",
		text: "HHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldello world",
	},
	{
		receiver_id: "f9598d63-47c4-474a-8c46-7bb63d5361fd",
		sender_id: "4efdc06f-d80e-4f80-9278-23b9f2",
		media_url: "",
		text: "HHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldello world",
	},
	{
		receiver_id: "f9598d63-47c4-474a-8c46-7bb63d5361fd",
		sender_id: "4efdc06f-d80e-4f80-9278-23b9f2",
		media_url: "",
		text: "HHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldello world",
	},
	{
		receiver_id: "f9598d63-47c4-474a-8c46-7bb63d5361fd",
		sender_id: "4efdc06f-d80e-4f80-9278-23b9f2",
		media_url: "",
		text: "HHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldHello worldello world",
	},
]);
const typingMessage = ref("");

onMounted(() => {
	_get(`/api/v1/room/all/${room_id.value}/member`)
		.then((data: RoomMemberInfo[]) => {
			const obj: { [key: string]: RoomMemberInfo } = {};
			for (let mem of data) {
				obj[mem.user_id] = mem;
			}
			members_map.value = obj;
			console.log(members_map.value);
		})
		.catch((err) => {
			toast.error(err.error);
			console.error(err);
		});
});

watch(
	() => room_id.value,
	(newRoomId) => {
		console.log(newRoomId);
	},
);

function submitMessage() {
	const sMsg: SMessage = {
		receiver_id: room_id.value,
		sender_id: userInfoStore.info!.user_id,
		text: typingMessage.value,
		media_url: "",
	};
	connectionStore.connection?.send(JSON.stringify(sMsg));
	typingMessage.value = "";
}
</script>
<template>
	<div class="flex h-[calc(100vh-3.5rem)] grow flex-col p-3 sm:h-screen">
		<div class="thin-scrollbar flex w-full grow flex-col overflow-auto">
			<Message
				v-for="msg in messages"
				:sender="members_map[msg.sender_id]"
				:message="msg" />
		</div>
		<div class="w-full shrink-0 p-2">
			<form
				@submit.prevent="submitMessage"
				class="flex w-full shrink-0">
				<input
					v-model="typingMessage"
					placeholder="Write message..."
					class="h-10 grow px-2" />
				<button
					type="submit"
					class="ml-2 flex h-10 items-center justify-center gap-2 bg-violet-500 px-3 text-sm font-semibold hover:bg-violet-500/70">
					<span>Send</span>
					<IconSend class="size-4" />
				</button>
			</form>
		</div>
	</div>
</template>

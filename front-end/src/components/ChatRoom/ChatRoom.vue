<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { usePersonalStore } from "@/stores/personal";
import { useWebSocketStore } from "@/stores/websocket";
import { _get } from "@/utils/fetch";
import { toast } from "vue3-toastify";
import type { SMessage } from "@/types/socket";
import type { RoomMemberInfo } from "@/types/user";
import RoomMessage from "@/components/ChatRoom/Message.vue";
import IconSend from "@/components/icons/IconSend.vue";
import { useRoomDataStore } from "@/stores/room_data";

const personalStore = usePersonalStore();
const websocketStore = useWebSocketStore();
const roomDataStore = useRoomDataStore();

const route = useRoute();
const room_id = computed<any>(() => route.params.room_id);
roomDataStore.current_id = room_id.value;

const typingMessage = ref("");

onMounted(() => {
	_get(`/api/v1/room/all/${room_id.value}/member`)
		.then((data: RoomMemberInfo[]) => {
			roomDataStore.updateMemberSet(data);
		})
		.catch((err) => {
			toast.error(err.error);
			console.error(err);
		});
	_get(`/api/v1/room/all/${room_id.value}/message`)
		.then((data: SMessage[]) => {
			roomDataStore.updateMessageList(data);
		})
		.catch((err) => {
			toast.error(err.error);
			console.error(err);
		});
});

watch(
	() => room_id.value,
	(newRoomId) => {
		roomDataStore.current_id = newRoomId;
		_get(`/api/v1/room/all/${room_id.value}/member`)
			.then((data: RoomMemberInfo[]) => {
				roomDataStore.updateMemberSet(data);
			})
			.catch((err) => {
				toast.error(err.error);
				console.error(err);
			});
		_get(`/api/v1/room/all/${room_id.value}/message`)
			.then((data: SMessage[]) => {
				roomDataStore.updateMessageList(data);
			})
			.catch((err) => {
				toast.error(err.error);
				console.error(err);
			});
	},
);

function submitMessage() {
	const sMsg: SMessage = {
		receiver_id: room_id.value,
		sender_id: personalStore.info!.user_id,
		text: typingMessage.value,
		media_url: "",
	};
	websocketStore.sendMessage(sMsg);
	typingMessage.value = "";
}
</script>
<template>
	<div class="flex h-[calc(100vh-3.5rem)] grow flex-col p-3 sm:h-screen">
		<div class="thin-scrollbar flex w-full grow flex-col overflow-auto">
			<RoomMessage
				v-for="msg in roomDataStore.messages"
				:sender="roomDataStore.members.get(msg.sender_id)"
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

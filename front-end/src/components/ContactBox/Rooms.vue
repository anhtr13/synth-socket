<script setup lang="ts">
import { ref, watch } from "vue";
import { useRouter } from "vue-router";
import { useRecentUpdatedStore } from "@/stores/recent_updated";
import { useGlobalStateStore } from "@/stores/global_state";
import { _get, _post } from "@/utils/fetch";
import { toast } from "vue3-toastify";
import type { Room } from "@/types/room";
import type { UserInfo } from "@/types/user";
import IconUsersPlus from "@/components/icons/IconUsersPlus.vue";
import IconGridPlus from "@/components/icons/IconGridPlus.vue";
import IconGroup from "@/components/icons/IconGroup.vue";
import IconClose from "@/components/icons/IconClose.vue";
import IconUserAstronaut from "@/components/icons/IconUserAstronaut.vue";
import IconSend from "@/components/icons/IconSend.vue";

const router = useRouter();
const globalStateStore = useGlobalStateStore();
const recentUpdatedStore = useRecentUpdatedStore();

const mode = ref<"default" | "create_room" | "invite_friend">("default");

const allRoomInput = ref("");
const allRooms = ref<Room[]>([]);

const ownedRoomInput = ref("");
const ownedRooms = ref<Room[]>([]);

const searchFriendInput = ref("");
const searchFriends = ref<UserInfo[]>([]);

const newRoomName = ref("");
const inviteChosenRoom = ref<Room | null>(null);

watch(mode, (value) => {
	if (value !== "default") {
		_get("/api/v1/room/owned")
			.then((data) => {
				ownedRooms.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}
});

let timerId: any;

watch(allRoomInput, (value) => {
	if (value) {
		clearTimeout(timerId);
		timerId = setTimeout(() => {
			_get("/api/v1/room/all", {
				searchParams: { search: value },
			})
				.then((data) => {
					allRooms.value = data;
				})
				.catch((err) => {
					console.error(err);
				});
		}, 300);
	}
});

watch(ownedRoomInput, (value) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/room/owned", {
			searchParams: { search: value },
		})
			.then((data) => {
				ownedRooms.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

watch(searchFriendInput, (value) => {
	if (value) {
		clearTimeout(timerId);
		timerId = setTimeout(() => {
			_get("/api/v1/friend", {
				searchParams: { search: value },
			})
				.then((data) => {
					searchFriends.value = data;
				})
				.catch((err) => {
					console.error(err);
				});
		}, 300);
	}
});

function handleCreateRoom() {
	if (!newRoomName.value) {
		toast.warn("Enter room name");
		return;
	}
	_post("/api/v1/room", {
		body: {
			room_name: newRoomName.value,
		},
	})
		.then((data) => {
			console.log(data);
			toast.success("Success!");
			newRoomName.value = "";
			mode.value = "default";
		})
		.catch((err) => {
			console.error(err);
			toast.error(err.error);
		});
}

function createRoomInvite(target_id: string) {
	if (!inviteChosenRoom.value) {
		toast.warn("Chose a room");
		return;
	}
	_post(`/api/v1/room/owned/${inviteChosenRoom.value.room_id}/invite/${target_id}`)
		.then((data) => {
			console.log(data);
			toast.success("Success!");
			inviteChosenRoom.value = null;
			ownedRoomInput.value = "";
			searchFriendInput.value = "";
			mode.value = "default";
		})
		.catch((err) => {
			console.error(err);
			toast.error(err.error);
		});
}

function handleClickRoom(room_id: string) {
	globalStateStore.showHeaderMobileDropdown = false;
	const room = recentUpdatedStore.roomSet.get(room_id);
	if (room) room.seen_last_message = true;
	router.push(`/room/${room_id}`);
}
</script>
<template>
	<!-- default mode -->
	<template v-if="mode === 'default'">
		<h3 class="text-lg font-semibold">Rooms</h3>
		<div class="flex flex-col">
			<div class="mt-2 flex w-full items-center gap-x-1">
				<button
					@click="mode = 'create_room'"
					class="flex w-36 items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
					<IconGridPlus class="size-4" />
					<span class="ml-2 text-sm"> Create room </span>
				</button>
				<button
					@click="mode = 'invite_friend'"
					class="flex w-36 items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
					<IconUsersPlus class="size-4" />
					<span class="ml-2 text-sm"> Invite friend </span>
				</button>
			</div>
			<span class="mt-1 text-sm text-neutral-400"> Create or invite friends to a room. </span>
		</div>
		<div class="relative mt-4 w-full">
			<input
				v-model="allRoomInput"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Find your room..." />
			<button
				v-show="allRoomInput"
				@click="allRoomInput = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<!-- recent rooms -->
		<div
			v-if="allRoomInput === ''"
			class="mt-2 flex flex-col">
			<span
				v-if="recentUpdatedStore.roomSet.size === 0"
				class="ml-1 text-sm text-neutral-400">
				You haven't joined any rooms yet.
			</span>
			<div
				v-else
				class="flex w-full flex-col-reverse">
				<button
					v-for="room in recentUpdatedStore.roomSet"
					@click="() => handleClickRoom(room[1].room_id)"
					class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30">
					<IconGroup
						v-if="!room[1].room_picture"
						class="mr-2 size-7 shrink-0 border border-neutral-400 p-1" />
					<img
						v-if="room[1].room_picture"
						:src="room[1].room_picture"
						class="mr-2 size-7 object-cover" />
					<div class="flex w-[calc(100%-2rem)] grow items-center gap-1 truncate">
						<div
							v-if="!room[1].seen_last_message"
							class="size-1 bg-green-500"></div>
						<span
							class="w-full truncate text-start text-sm"
							:class="room[1].seen_last_message ? 'font-semibold' : 'font-bold'">
							{{ room[1].room_name }}
						</span>
					</div>
				</button>
			</div>
		</div>
		<!-- searched rooms -->
		<div
			v-else
			class="mt-2 flex flex-col">
			<span
				v-if="allRooms.length === 0"
				class="ml-1 text-sm text-neutral-400">
				No room found.
			</span>
			<button
				v-else
				v-for="room in allRooms"
				@click="() => handleClickRoom(room.room_id)"
				class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30">
				<IconGroup
					v-if="!room.room_picture"
					class="mr-3 size-6" />
				<img
					v-if="room.room_picture"
					:src="room.room_picture"
					class="mr-3 size-6 object-cover" />
				<span class="w-[calc(100%-2rem)] truncate text-start text-sm font-bold">
					{{ room.room_name }}
				</span>
			</button>
		</div>
	</template>

	<!-- create mode -->
	<template v-else-if="mode === 'create_room'">
		<h3 class="text-lg font-semibold">Create rooms</h3>
		<div class="relative my-2 w-full">
			<input
				v-model="newRoomName"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Enter room name..." />
			<button
				v-show="newRoomName"
				@click="newRoomName = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<button
			@click="handleCreateRoom"
			class="mt-2 flex w-full items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
			<IconGridPlus class="size-4" />
			<span class="ml-2 text-sm"> Create </span>
		</button>
		<button
			@click="mode = 'default'"
			class="mt-2 flex w-full items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
			<IconClose class="size-4" />
			<span class="ml-2 text-sm"> Cancel </span>
		</button>
		<div class="mt-4 flex flex-col">
			<h3 class="mb-1 font-semibold">Created rooms:</h3>
			<span
				v-if="ownedRooms.length === 0"
				class="ml-1 text-sm text-neutral-400">
				You haven't created any room.
			</span>
			<button
				v-else
				v-for="room in ownedRooms"
				@click="() => handleClickRoom(room.room_id)"
				class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30">
				<IconGroup
					v-if="!room.room_picture"
					class="mr-3 size-6" />
				<img
					v-if="room.room_picture"
					:src="room.room_picture"
					class="mr-3 size-6 object-cover" />
				<span class="w-[calc(100%-2rem)] truncate text-start text-sm">
					{{ room.room_name }}
				</span>
			</button>
		</div>
	</template>

	<!-- invite mode -->
	<template v-else-if="mode === 'invite_friend'">
		<h3 class="text-lg font-semibold">Invite to:</h3>
		<div v-if="inviteChosenRoom">
			<div class="my-1 flex w-full items-center bg-violet-500/30 px-3 py-2">
				<IconGroup
					v-if="!inviteChosenRoom.room_picture"
					class="mr-3 size-7" />
				<img
					v-if="inviteChosenRoom?.room_picture"
					:src="inviteChosenRoom.room_picture"
					class="mr-3 size-7 object-cover" />
				<span class="w-[calc(100%-2rem)] truncate text-sm">
					{{ inviteChosenRoom.room_name }}
				</span>
			</div>
			<div class="relative mt-4 w-full">
				<input
					v-model="searchFriendInput"
					class="h-10 w-full pr-8 pl-2"
					placeholder="Find friend name..." />
				<button
					v-show="searchFriendInput"
					@click="searchFriendInput = ''"
					class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
					<IconClose class="size-4" />
				</button>
			</div>
			<span class="mt-1 ml-1 text-sm text-neutral-400">
				Invite friends to your selected room.
			</span>
			<div
				v-if="searchFriendInput === ''"
				class="mt-4 flex w-full flex-col-reverse">
				<div
					v-for="friend in recentUpdatedStore.friendSet"
					class="flex items-center justify-between py-2 pl-3">
					<div class="flex w-full max-w-[calc(100%-2rem)] items-center">
						<div class="relative shrink-0">
							<img
								v-if="friend[1].profile_image"
								:src="friend[1].profile_image"
								class="size-8 object-cover" />
							<IconUserAstronaut
								v-else
								class="size-8 object-cover" />
							<div
								v-show="friend[1].online_status === 'online'"
								class="absolute right-0 bottom-0 size-2 bg-green-500"></div>
						</div>
						<div class="ml-2 flex w-full flex-col">
							<h3 class="max-w-[calc(100%-1.5rem)] truncate font-medium">
								{{ friend[1].user_name }}
							</h3>
							<span
								v-if="friend[1].online_status === 'online'"
								class="text-xs text-neutral-400">
								online now
							</span>
							<span
								v-else-if="friend[1].online_status"
								class="text-xs text-neutral-400">
								{{ new Date(friend[1].online_status).toLocaleString() }}
							</span>
						</div>
					</div>
					<button
						@click="() => createRoomInvite(friend[1].user_id)"
						title="Invite"
						class="bg-violet-500/30 px-2 py-1 hover:bg-violet-500">
						<IconSend class="size-4" />
					</button>
				</div>
				<span
					v-show="recentUpdatedStore.friendSet.size === 0"
					class="text-sm text-neutral-400">
					You have no friends yet.
				</span>
			</div>
			<div
				v-else
				class="mt-4 flex w-full flex-col">
				<div
					v-for="friend in searchFriends"
					class="flex items-center justify-between py-2 pl-3">
					<div class="flex max-w-[calc(100%-2rem)] items-center">
						<img
							v-if="friend.profile_image"
							:src="friend.profile_image"
							class="size-8 object-cover" />
						<IconUserAstronaut
							v-else
							class="size-8 object-cover" />
						<h3 class="ml-2 max-w-[calc(100%-1.5rem)] truncate font-medium">
							{{ friend.user_name }}
						</h3>
					</div>
					<button
						@click="() => createRoomInvite(friend.user_id)"
						title="Invite"
						class="bg-violet-500/30 px-2 py-1 hover:bg-violet-500">
						<IconSend class="size-4" />
					</button>
				</div>
				<span
					v-show="searchFriends.length === 0"
					class="text-sm text-neutral-400">
					No friends found.
				</span>
			</div>
			<button
				@click="inviteChosenRoom = null"
				class="mt-4 flex w-full items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
				<IconClose class="size-4" />
				<span class="ml-2 text-sm"> Cancel </span>
			</button>
		</div>
		<div
			v-else
			class="w-full">
			<div class="relative mt-2 w-full">
				<input
					v-model="ownedRoomInput"
					class="h-10 w-full pr-8 pl-2"
					placeholder="Find your room..." />
				<button
					v-show="ownedRoomInput"
					@click="ownedRoomInput = ''"
					class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
					<IconClose class="size-4" />
				</button>
			</div>
			<span class="mt-1 ml-1 text-sm text-neutral-400">
				Select your owned room to invite friends.
			</span>
			<div class="mt-2 flex flex-col">
				<span
					v-if="ownedRooms.length === 0"
					class="ml-1 text-sm text-neutral-400">
					No room found.
				</span>
				<div
					v-else
					v-for="room in ownedRooms"
					class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30"
					@click="inviteChosenRoom = room">
					<IconGroup
						v-if="!room.room_picture"
						class="mr-3 ml-1 size-6" />
					<img
						v-if="room.room_picture"
						:src="room.room_picture"
						class="mr-3 size-6 object-cover" />
					<span class="w-[calc(100%-2rem)] truncate text-sm">{{ room.room_name }}</span>
				</div>
			</div>
			<button
				@click="mode = 'default'"
				class="mt-4 flex w-full items-center justify-center bg-violet-500 py-2 text-start hover:bg-violet-500/70">
				<IconClose class="size-4" />
				<span class="ml-2 text-sm"> Cancel </span>
			</button>
		</div>
	</template>
</template>

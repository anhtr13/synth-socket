<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { _get, _post } from "@/utils/fetch";
import type { Room } from "@/types/room";
import type { UserInfo } from "@/types/user";
import IconUsersPlus from "@/components/icons/IconUsersPlus.vue";
import IconGridPlus from "@/components/icons/IconGridPlus.vue";
import IconGroup from "@/components/icons/IconGroup.vue";
import IconClose from "@/components/icons/IconClose.vue";
import IconUserAstronaut from "@/components/icons/IconUserAstronaut.vue";
import IconSend from "@/components/icons/IconSend.vue";

const mode = ref<"default" | "create_room" | "invite_friend">("default");
const rooms = ref<Room[]>([]);
const ownedRoom = ref<Room[]>([]);

const searchRoom = ref("");
const newRoomName = ref("");
const searchOwnedRoom = ref("");
const searchFriend = ref("");
const inviteChosenRoom = ref<Room | null>(null);
const friends = ref<UserInfo[]>([]);

let timerId: any;

onMounted(() => {
	_get("/api/v1/room/all")
		.then((data) => {
			rooms.value = data;
		})
		.catch((err) => {
			console.error(err);
		});
});

watch(mode, (value) => {
	if (value !== "default") {
		_get("/api/v1/room/owned")
			.then((data) => {
				console.log(data);
				ownedRoom.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}
});

watch(inviteChosenRoom, (value) => {
	if (mode.value === "invite_friend" && value) {
		_get("/api/v1/friend")
			.then((data) => {
				friends.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}
});

watch(searchRoom, (value) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/room/all", {
			searchParams: { search: value },
		})
			.then((data) => {
				rooms.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

watch(searchOwnedRoom, (value) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/room/owned", {
			searchParams: { search: value },
		})
			.then((data) => {
				ownedRoom.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

watch(searchFriend, (value) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/friend", {
			searchParams: { search: value },
		})
			.then((data) => {
				friends.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

function createRoom() {
	if (!newRoomName.value) {
		return;
	}
	_post("/api/v1/room", {
		body: {
			room_name: newRoomName.value,
		},
	})
		.then((data) => {
			console.log(data);
		})
		.catch((err) => console.error(err));
}

function createRoomInvite(target_id: string) {
	if (!inviteChosenRoom.value) {
		return;
	}
	_post(`/api/v1/room/owned/${inviteChosenRoom.value.room_id}/invite/${target_id}`)
		.then((data) => {
			console.log(data);
		})
		.catch((err) => console.error(err));
}
</script>
<template>
	<template v-if="mode === 'default'">
		<h3 class="text-lg font-semibold">Chat rooms</h3>
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
				v-model="searchRoom"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Find your room..." />
			<button
				v-show="searchRoom"
				@click="searchRoom = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<div class="mt-2 flex flex-col">
			<span
				v-if="rooms.length === 0"
				class="ml-1 text-sm text-neutral-400">
				No room founded.
			</span>
			<div
				v-else
				v-for="room in rooms"
				class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30">
				<IconGroup
					v-if="!room.room_picture"
					class="mr-3 size-5" />
				<img
					v-if="room.room_picture"
					:src="room.room_picture"
					class="mr-3 size-5" />
				<span class="w-[calc(100%-2rem)] truncate text-sm">{{ room.room_name }}</span>
			</div>
		</div>
	</template>
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
			@click="createRoom"
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
			<h3 class="font-semibold">Created rooms:</h3>
			<span
				v-if="ownedRoom.length === 0"
				class="ml-1 text-sm text-neutral-400">
				No room founded.
			</span>
			<div
				v-else
				v-for="room in ownedRoom"
				class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30">
				<IconGroup
					v-if="!room.room_picture"
					class="mr-3 size-5" />
				<img
					v-if="room.room_picture"
					:src="room.room_picture"
					class="mr-3 size-5" />
				<span class="w-[calc(100%-2rem)] truncate text-sm">{{ room.room_name }}</span>
			</div>
		</div>
	</template>
	<template v-else-if="mode === 'invite_friend'">
		<h3 class="text-lg font-semibold">Invite to:</h3>
		<div v-if="inviteChosenRoom">
			<div class="my-1 flex w-full items-center bg-violet-500/30 px-3 py-2">
				<IconGroup
					v-if="!inviteChosenRoom.room_picture"
					class="mr-3 size-5" />
				<img
					v-if="inviteChosenRoom?.room_picture"
					:src="inviteChosenRoom.room_picture"
					class="mr-3 size-5" />
				<span class="w-[calc(100%-2rem)] truncate text-sm">{{
					inviteChosenRoom.room_name
				}}</span>
			</div>
			<div class="relative mt-4 w-full">
				<input
					v-model="searchFriend"
					class="h-10 w-full pr-8 pl-2"
					placeholder="Find friend name..." />
				<button
					v-show="searchFriend"
					@click="searchFriend = ''"
					class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
					<IconClose class="size-4" />
				</button>
			</div>
			<span class="mt-1 ml-1 text-sm text-neutral-400">
				Invite friends to your selected room.
			</span>
			<div class="mt-4 flex w-full flex-col">
				<div
					v-for="friend in friends"
					class="flex items-center justify-between py-2 pl-3">
					<div class="flex max-w-[calc(100%-2rem)] items-center">
						<img
							v-if="friend.profile_image"
							:src="friend.profile_image"
							class="size-5 object-cover" />
						<IconUserAstronaut
							v-else
							class="size-5 object-cover" />
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
					v-show="friends.length === 0"
					class="text-sm text-neutral-400">
					No user found.
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
					v-model="searchOwnedRoom"
					class="h-10 w-full pr-8 pl-2"
					placeholder="Find your room..." />
				<button
					v-show="searchOwnedRoom"
					@click="searchOwnedRoom = ''"
					class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
					<IconClose class="size-4" />
				</button>
			</div>
			<span class="mt-1 ml-1 text-sm text-neutral-400">
				Select your owned room to invite friends.
			</span>
			<div class="mt-2 flex flex-col">
				<span
					v-if="ownedRoom.length === 0"
					class="ml-1 text-sm text-neutral-400">
					No room founded.
				</span>
				<div
					v-else
					v-for="room in ownedRoom"
					class="flex w-full cursor-pointer items-center px-3 py-2 hover:bg-violet-500/30"
					@click="inviteChosenRoom = room">
					<IconGroup
						v-if="!room.room_picture"
						class="mr-3 ml-1 size-5" />
					<img
						v-if="room.room_picture"
						:src="room.room_picture"
						class="mr-3 size-5" />
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

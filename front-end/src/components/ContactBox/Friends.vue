<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import { _get, _post } from "@/utils/fetch";
import { usePersonalStore } from "@/stores/personal";
import { useRecentUpdatedStore } from "@/stores/recent_updated";
import { toast } from "vue3-toastify";
import type { UserInfo } from "@/types/user";
import IconUserPlus from "@/components/icons/IconUserPlus.vue";
import IconUserChecked from "@/components/icons/IconUserChecked.vue";
import IconArrowLeft from "@/components/icons/IconArrowLeft.vue";
import IconUserAstronaut from "@/components/icons/IconUserAstronaut.vue";
import IconClose from "@/components/icons/IconClose.vue";

const personalStore = usePersonalStore();
const recentUpdatedStore = useRecentUpdatedStore();
const title = ref<"Friends" | "Users">("Friends");

const searchFriendInput = ref("");
const searchFriends = ref<UserInfo[]>([]);

const searchUserInput = ref("");
const searchUsers = ref<UserInfo[]>([]);

onMounted(() => {
	_get("/api/v1/user")
		.then((data) => {
			searchUsers.value = data;
		})
		.catch((err) => {
			console.error(err);
		});
});

let timerId: any;
watch(searchFriendInput, (newValue) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/friend", { searchParams: { search: newValue } })
			.then((data) => {
				searchFriends.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});
watch(searchUserInput, (newValue) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/user", { searchParams: { search: newValue } })
			.then((data) => {
				searchUsers.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

function requestFriend(target_id: string) {
	_post(`api/v1/user/${target_id}/friend_request`)
		.then((data) => {
			console.log(data);
			toast.success("Success");
		})
		.catch((err) => {
			toast.error(err.error);
			console.error(err);
		});
}
</script>

<template>
	<!-- Friends -->
	<div v-show="title === 'Friends'">
		<div class="flex w-full items-center justify-between">
			<h3 class="text-lg font-semibold">{{ title }}</h3>
		</div>
		<div class="relative mt-2 h-auto w-full">
			<input
				v-model="searchFriendInput"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Search friend..." />
			<button
				v-show="searchFriendInput"
				@click="searchFriendInput = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<div
			v-if="searchFriendInput === ''"
			class="mt-4 flex w-full flex-col-reverse">
			<div
				v-for="friend in recentUpdatedStore.friendSet"
				class="flex w-full items-center px-3 py-2 hover:bg-violet-500/30">
				<div class="relative">
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
			<span
				v-show="recentUpdatedStore.friendSet.size === 0"
				class="text-sm text-neutral-400">
				You don't have any friend.
			</span>
		</div>
		<div
			v-else
			class="mt-4 flex w-full flex-col">
			<div
				v-for="friend in searchFriends"
				class="flex w-full items-center px-3 py-2 hover:bg-violet-500/30">
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
			<span
				v-show="searchFriends.length === 0"
				class="text-sm text-neutral-400">
				No one found.
			</span>
		</div>
		<button
			class="mt-4 flex h-9 w-full items-center justify-center bg-violet-500 hover:bg-violet-500/70"
			@click="title = 'Users'">
			<IconUserPlus class="size-4" />
			<span class="ml-2 text-sm"> Add new friend? </span>
		</button>
	</div>

	<!-- Users -->
	<div v-show="title === 'Users'">
		<div class="flex w-full items-center justify-between">
			<h3 class="text-lg font-semibold">{{ title }}</h3>
		</div>
		<div class="relative mt-2 h-auto w-full">
			<input
				v-model="searchUserInput"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Find user..." />
			<button
				v-show="searchUserInput"
				@click="searchUserInput = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<div class="mt-4 flex w-full flex-col">
			<div
				v-for="user in searchUsers"
				class="flex items-center justify-between py-2 pl-3 hover:bg-violet-500/30">
				<div class="flex max-w-[calc(100%-2rem)] items-center">
					<img
						v-if="user.profile_image"
						:src="user.profile_image"
						class="size-8 object-cover" />
					<IconUserAstronaut
						v-else
						class="size-8 object-cover" />
					<h3 class="ml-2 max-w-[calc(100%-1.5rem)] truncate font-medium">
						{{ user.user_name }}
					</h3>
				</div>
				<div
					v-if="user.is_friend"
					class="px-2">
					<IconUserChecked class="size-4" />
				</div>
				<button
					v-else-if="user.user_id !== personalStore.info?.user_id"
					@click="() => requestFriend(user.user_id)"
					title="Add"
					class="bg-violet-500/30 px-2 py-1 hover:bg-violet-500">
					<IconUserPlus class="size-4" />
				</button>
			</div>
			<span
				v-show="searchUsers.length === 0"
				class="text-sm text-neutral-400">
				No user found.
			</span>
		</div>
		<button
			@click="title = 'Friends'"
			class="mt-4 flex w-full items-center justify-center bg-violet-500 py-2 hover:bg-violet-500/70">
			<IconArrowLeft class="size-4" />
			<span class="ml-2 text-sm font-semibold">Back</span>
		</button>
	</div>
</template>

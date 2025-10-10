<script setup lang="ts">
import { ref, watch, onMounted } from "vue";
import { _get, _post } from "@/utils/fetch";
import { useUserInfoStore } from "@/stores/user";
import type { UserInfo } from "@/types/user";
import IconUserPlus from "@/components/icons/IconUserPlus.vue";
import IconUserChecked from "@/components/icons/IconUserChecked.vue";
import IconArrowLeft from "@/components/icons/IconArrowLeft.vue";
import IconUserAstronaut from "@/components/icons/IconUserAstronaut.vue";
import IconClose from "@/components/icons/IconClose.vue";
import { toast } from "vue3-toastify";

const userInfoStore = useUserInfoStore();
const title = ref<"Friends" | "Users">("Friends");
const searchFriend = ref("");
const searchUser = ref("");
const friends = ref<UserInfo[]>([]);
const users = ref<UserInfo[]>([]);

onMounted(() => {
	_get("/api/v1/friend")
		.then((data) => {
			friends.value = data;
		})
		.catch((err) => {
			console.error(err);
		});
});

let timerId: any;

watch(searchFriend, (newValue) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/friend", { searchParams: { search: newValue } })
			.then((data) => {
				friends.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

watch(searchUser, (newValue) => {
	clearTimeout(timerId);
	timerId = setTimeout(() => {
		_get("/api/v1/user", { searchParams: { search: newValue } })
			.then((data) => {
				users.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}, 300);
});

watch(title, (newValue) => {
	if (newValue === "Friends") {
		_get("/api/v1/friend")
			.then((data) => {
				friends.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	} else if (newValue === "Users") {
		_get("/api/v1/user")
			.then((data) => {
				users.value = data;
			})
			.catch((err) => {
				console.error(err);
			});
	}
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
				v-model="searchFriend"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Search friend..." />
			<button
				v-show="searchFriend"
				@click="searchFriend = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<div class="mt-4 flex w-full flex-col">
			<div
				v-for="friend in friends"
				class="flex w-full items-center px-3 py-2 hover:bg-violet-500/30">
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
			<span
				v-show="friends.length === 0"
				class="text-sm text-neutral-400">
				No friend found.
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
				v-model="searchUser"
				class="h-10 w-full pr-8 pl-2"
				placeholder="Find user..." />
			<button
				v-show="searchUser"
				@click="searchUser = ''"
				class="absolute top-1/2 right-2 z-10 flex size-5 -translate-y-1/2 items-center justify-center bg-violet-500/30 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
		</div>
		<div class="mt-4 flex w-full flex-col">
			<div
				v-for="user in users"
				class="flex items-center justify-between py-2 pl-3 hover:bg-violet-500/30">
				<div class="flex max-w-[calc(100%-2rem)] items-center">
					<img
						v-if="user.profile_image"
						:src="user.profile_image"
						class="size-5 object-cover" />
					<IconUserAstronaut
						v-else
						class="size-5 object-cover" />
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
					v-else-if="user.user_id !== userInfoStore.info?.user_id"
					@click="() => requestFriend(user.user_id)"
					title="Add"
					class="bg-violet-500/30 px-2 py-1 hover:bg-violet-500">
					<IconUserPlus class="size-4" />
				</button>
			</div>
			<span
				v-show="users.length === 0"
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

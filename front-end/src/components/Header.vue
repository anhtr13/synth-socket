<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { useContactBoxStore, type FocusingContact } from "@/stores/contact-box";
import { useUserInfoStore } from "@/stores/user";
import type { UserNoti } from "@/types/user";
import { _get, _post } from "@/utils/fetch";
import IconApp from "./icons/IconApp.vue";
import IconFriend from "./icons/IconFriend.vue";
import IconUserAstronaut from "./icons/IconUserAstronaut.vue";
import IconEdit from "./icons/IconEdit.vue";
import IconLogout from "./icons/IconLogout.vue";
import IconBell from "./icons/IconBell.vue";

const router = useRouter();
const userInfoStore = useUserInfoStore();
const contactBoxStore = useContactBoxStore();

const notifications = ref<UserNoti[]>([]);
const showSettingDropdown = ref(false);
const showNotificationDropdown = ref(false);

function handleHeaderClick(element: FocusingContact) {
	if (contactBoxStore.focusing === element) {
		contactBoxStore.toggleMobileDropdown();
	} else {
		contactBoxStore.focusing = element;
		contactBoxStore.showMobileDropdown = true;
	}
	notifications.value.push({
		message: "nnotificationsnotificationsnotificationsnotificationsnotificationsotifications",
		created_at: "qwdqdqd",
		id_ref: "791a7cc9-9cb2-4132-b315-888760f72c45",
		notification_id: "f51fe755-6578-4bec-82c0-cb53941e8a8e",
		seen: false,
		type: "friend_request",
		user_id: "e7d964ae-c09c-41b2-9842-dc8e7213118d",
	});
}
function handleLogout() {
	_post("/api/v1/auth/logout")
		.then(() => {
			userInfoStore.updateInfo(null);
			localStorage.removeItem("access_token");
			localStorage.removeItem("refresh_token");
			router.push("/auth/login");
		})
		.catch((err) => {
			console.error(err);
		});
}
</script>

<template>
	<header
		class="flex h-14 w-full items-center justify-between border-b border-neutral-700 bg-black sm:h-full sm:w-14 sm:flex-col sm:border-r sm:border-b-0">
		<nav class="flex h-auto w-auto items-center sm:flex-col">
			<button
				@click="handleHeaderClick('room')"
				class="ml-3 flex size-8 items-center justify-center hover:text-violet-500 sm:mt-4 sm:ml-0"
				:class="contactBoxStore.focusing === 'room' ? 'text-violet-500' : ''">
				<IconApp class="size-6" />
			</button>
			<button
				@click="handleHeaderClick('friend')"
				class="ml-2 flex size-8 items-center justify-center hover:text-violet-500 sm:mt-2 sm:ml-0"
				:class="contactBoxStore.focusing === 'friend' ? 'text-violet-500' : ''">
				<IconFriend class="size-5" />
			</button>
		</nav>
		<nav class="relative flex h-auto w-auto items-center sm:flex-col">
			<div v-clickoutside="() => (showNotificationDropdown = false)">
				<button
					@click="() => (showNotificationDropdown = !showNotificationDropdown)"
					class="relative mr-1 flex size-8 items-center justify-center hover:text-violet-500 sm:mr-0 sm:mb-1">
					<IconBell class="size-5" />
					<span
						v-if="notifications.length > 0"
						class="absolute top-0 right-0 size-4 border border-black bg-red-500 text-center text-white"
						style="font-size: 10px">
						{{ notifications.length }}
					</span>
				</button>
				<div
					v-if="showNotificationDropdown && notifications.length > 0"
					class="thin-scrollbar absolute top-12 right-2 z-10 flex max-h-96 w-72 flex-col overflow-y-auto border border-neutral-700 bg-black p-3 sm:top-auto sm:right-auto sm:bottom-2 sm:left-12">
					<button
						class="w-full truncate px-3 py-2 text-start hover:bg-neutral-900"
						:title="noti.message"
						v-for="noti in notifications">
						{{ noti.message }}
					</button>
				</div>
			</div>
			<div
				v-clickoutside="() => (showSettingDropdown = false)"
				class="mr-3 flex size-8 items-center justify-center sm:mr-0 sm:mb-4">
				<button
					@click="showSettingDropdown = !showSettingDropdown"
					v-if="!userInfoStore.info || !userInfoStore.info.profile_image"
					class="hover:text-violet-500">
					<IconUserAstronaut class="size-5" />
				</button>
				<button
					@click="showSettingDropdown = !showSettingDropdown"
					v-if="userInfoStore.info?.profile_image">
					<img
						class="size-6 object-cover"
						:src="userInfoStore.info.profile_image" />
				</button>
				<div
					v-if="showSettingDropdown"
					class="absolute top-12 right-2 z-10 flex w-64 flex-col border border-neutral-700 bg-black p-4 sm:top-auto sm:right-auto sm:bottom-2 sm:left-12">
					<div class="flex items-center gap-3 bg-neutral-900 px-4 py-4">
						<img
							v-if="userInfoStore.info?.profile_image"
							class="size-12 object-cover"
							:src="userInfoStore.info.profile_image" />
						<IconUserAstronaut
							v-if="!userInfoStore.info?.profile_image"
							class="size-12" />
						<div class="flex flex-col gap-0">
							<h3 class="text-xl font-bold">{{ userInfoStore.info?.user_name }}</h3>
							<span class="text-xs text-neutral-300">online</span>
						</div>
					</div>
					<div class="mt-3 bg-neutral-900 p-2">
						<button class="flex w-full items-center px-3 py-2 text-start hover:bg-black/50">
							<IconEdit class="mr-2 size-4" />
							<span class="text-sm">Edit profile</span>
						</button>
					</div>
					<div class="mt-3 bg-neutral-900 p-2">
						<button
							@click="handleLogout"
							class="flex w-full items-center px-3 py-2 text-start hover:bg-black/50">
							<IconLogout class="mr-2 size-4" />
							<span class="text-sm">Log out</span>
						</button>
					</div>
				</div>
			</div>
		</nav>
	</header>
</template>

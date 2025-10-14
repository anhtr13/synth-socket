<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { useContactBoxStore, type FocusingContact } from "@/stores/contact-box";
import { useUserInfoStore } from "@/stores/user";
import { _delete, _get, _patch, _post } from "@/utils/fetch";
import { toast } from "vue3-toastify";
import type { SNotification } from "@/types/socket";
import IconApp from "./icons/IconApp.vue";
import IconFriend from "./icons/IconFriend.vue";
import IconUserAstronaut from "./icons/IconUserAstronaut.vue";
import IconEdit from "./icons/IconEdit.vue";
import IconLogout from "./icons/IconLogout.vue";
import IconBell from "./icons/IconBell.vue";
import IconClose from "./icons/IconClose.vue";

const router = useRouter();
const userInfoStore = useUserInfoStore();
const contactBoxStore = useContactBoxStore();

const notifications = ref<SNotification[]>([]);
const showSettingDropdown = ref(false);
const showNotificationDropdown = ref(false);
const chosenProfileImage = ref(userInfoStore.info?.profile_image || null);
const showProfileDialog = ref(false);

const notiCount = computed(() => {
	let count = 0;
	notifications.value.forEach((noti) => {
		if (!noti.seen) {
			count++;
		}
	});
	return count;
});

onMounted(() => {
	_get("/api/v1/notification")
		.then((data) => {
			console.log("notifications:", data);
			notifications.value = data;
		})
		.catch((err) => {
			toast.error(err.error);
		});
});

function handleHeaderClick(element: FocusingContact) {
	if (contactBoxStore.focusing === element) {
		contactBoxStore.toggleMobileDropdown();
	} else {
		contactBoxStore.focusing = element;
		contactBoxStore.showMobileDropdown = true;
	}
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

function handleNotificationAction(noti: SNotification, action: "accept" | "reject") {
	let path = `/api/v1/${noti.type}/${noti.id_ref}`;
	if (action === "accept") {
		_post(path)
			.then(() => {
				toast.success("Success!");
			})
			.catch((err) => {
				console.error(err);
				toast.error(err.error);
			});
	} else {
		_delete(path)
			.then(() => {
				toast.success("Success!");
			})
			.catch((err) => {
				console.error(err);
				toast.error(err.error);
			});
	}
	_post(`/api/v1/notification/${noti.notification_id}`).then(() => {
		noti.seen = true;
	});
}

function handleSaveNewAvatar() {
	_patch("/api/v1/me/info", {
		body: {
			profile_image: chosenProfileImage.value,
		},
	})
		.then(() => {
			userInfoStore.info!.profile_image = chosenProfileImage.value;
			showProfileDialog.value = false;
			toast.success("Success!");
		})
		.catch((err) => {
			console.log(err);
			toast.error(err.error);
		});
}
</script>

<template>
	<header
		class="flex h-14 w-full shrink-0 items-center justify-between border-b border-neutral-700 bg-black sm:h-full sm:w-14 sm:flex-col sm:border-r sm:border-b-0">
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
						v-if="notiCount > 0"
						class="absolute top-0 right-0 size-4 border border-black bg-red-500 text-center text-white"
						style="font-size: 10px">
						{{ notiCount }}
					</span>
				</button>
				<div
					v-if="showNotificationDropdown"
					class="thin-scrollbar absolute top-12 right-2 z-10 flex max-h-96 w-72 flex-col overflow-y-auto border border-neutral-700 bg-black p-3 sm:top-auto sm:right-auto sm:bottom-2 sm:left-12">
					<span
						v-if="notifications.length === 0"
						class="text-sm text-neutral-400">
						Notification is empty
					</span>
					<div
						class="w-full px-3 py-2 text-start hover:bg-neutral-900"
						:title="noti.message"
						v-for="noti in notifications">
						<span class="mb-1 line-clamp-2 w-full font-semibold break-all">
							{{ noti.message }}
						</span>
						<div
							v-if="!noti.seen"
							class="flex gap-x-2">
							<button
								@click="() => handleNotificationAction(noti, 'accept')"
								class="h-8 w-20 bg-violet-500 text-sm hover:bg-violet-500/70">
								Accept
							</button>
							<button
								@click="() => handleNotificationAction(noti, 'reject')"
								class="h-8 w-20 bg-violet-500 text-sm hover:bg-violet-500/70">
								Reject
							</button>
						</div>
					</div>
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
						<button
							@click="showProfileDialog = true"
							class="flex w-full items-center px-3 py-2 text-start hover:bg-black/50">
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
	<dialog
		v-if="showProfileDialog"
		class="z-20 flex h-screen w-screen items-center justify-center bg-transparent text-white backdrop-blur-lg">
		<div class="relative flex w-full max-w-[26rem] flex-col items-center border bg-black p-6">
			<button
				@click="
					() => {
						chosenProfileImage = userInfoStore.info?.profile_image || null;
						showProfileDialog = false;
					}
				"
				class="absolute top-2 right-2 p-1 hover:bg-violet-500">
				<IconClose class="size-4" />
			</button>
			<img
				v-if="chosenProfileImage"
				class="size-12 object-cover"
				:src="chosenProfileImage" />
			<IconUserAstronaut
				v-else
				class="size-12" />
			<h3 class="mt-1 text-xl font-bold">{{ userInfoStore.info?.user_name }}</h3>
			<p class="mt-4 text-sm text-neutral-400">Change your avatar</p>
			<div class="mt-2 flex w-full flex-wrap items-center justify-center gap-2 px-6">
				<button
					v-for="i in 8"
					@click="chosenProfileImage = `/images/avatar${i}.jpg`">
					<img
						class="size-12 object-cover"
						:src="`/images/avatar${i}.jpg`" />
				</button>
			</div>
			<button
				@click="handleSaveNewAvatar"
				class="mt-6 w-32 bg-violet-500 py-2 text-sm hover:bg-violet-500/70">
				Save avatar
			</button>
		</div>
	</dialog>
</template>

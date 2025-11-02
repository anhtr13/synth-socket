<script setup lang="ts">
import { onBeforeUnmount, onMounted } from "vue";
import { useRouter } from "vue-router";
import { usePersonalStore } from "@/stores/personal";
import { useWebSocketStore } from "@/stores/websocket";
import { useRecentUpdatedStore } from "@/stores/recent_updated";
import { _get } from "@/utils/fetch";
import Header from "@/components/Header.vue";
import ContactBox from "@/components/ContactBox/ContactBox.vue";

const HOST = import.meta.env.MODE === "development" ? "localhost:3000" : window.location.host;

const router = useRouter();
const personalStore = usePersonalStore();
const webSocketStore = useWebSocketStore();
const recentUpdatedStore = useRecentUpdatedStore();

onMounted(async () => {
	await _get("/api/v1/me/info")
		.then((data) => {
			personalStore.info = data;
		})
		.catch((err) => {
			personalStore.info = null;
			console.error("here", err);
			router.push("/auth/login");
		});
	await Promise.all([
		_get("/api/v1/room/all")
			.then((data) => {
				console.log(data);
				recentUpdatedStore.updateRoomSet(data);
			})
			.catch((err) => {
				console.error(err);
			}),
		_get("/api/v1/friend")
			.then((data) => {
				recentUpdatedStore.updateFriendSet(data);
			})
			.catch((err) => {
				console.error(err);
			}),
	]);
	if (!window["WebSocket"]) {
		alert("Your browser does not support WebSockets.");
		return;
	}
	let access_token = window.localStorage.getItem("access_token");
	if (!access_token) {
		alert("Cannot connect to WebSocket server: access_token not found");
		return;
	}
	let protocol = "ws://";
	if (window.location.protocol === "https:") {
		protocol = "wss://";
	}
	const url = protocol + HOST + "/ws?access_token=" + access_token;
	webSocketStore.connect(url);
});

onBeforeUnmount(() => {
	webSocketStore.disconnect();
});
</script>

<template>
	<template v-if="personalStore.info">
		<Header />
		<main class="relative flex h-auto w-auto grow flex-col sm:flex-row">
			<ContactBox />
			<router-view />
		</main>
	</template>
</template>

<style scoped></style>

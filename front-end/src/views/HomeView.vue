<script setup lang="ts">
import { onBeforeMount, onBeforeUnmount, onMounted, watch } from "vue";
import { useRouter } from "vue-router";
import { usePersonalStore } from "@/stores/personal";
import { useWebSocketStore } from "@/stores/websocket";
import { _get } from "@/utils/fetch";
import Header from "@/components/Header.vue";
import ContactBox from "@/components/ContactBox/ContactBox.vue";

const HOST = import.meta.env.MODE === "development" ? "localhost:3000" : window.location.host;

const router = useRouter();
const personalStore = usePersonalStore();
const webSocketStore = useWebSocketStore();

onBeforeMount(async () => {
	if (!personalStore.info) {
		_get("/api/v1/me/info")
			.then((data) => {
				personalStore.info = data;
			})
			.catch((err) => {
				personalStore.info = null;
				console.error("here", err);
				router.push("/auth/login");
			});
	}
});

watch(
	() => personalStore.info,
	(info) => {
		if (info) {
			if (!window["WebSocket"]) {
				alert("Your browser does not support WebSockets.");
				return null;
			}
			var access_token = window.localStorage.getItem("access_token");
			if (!access_token) {
				alert("Cannot connect to WebSocket server: access_token not found");
				return null;
			}
			var protocol = "ws://";
			if (window.location.protocol === "https:") {
				protocol = "wss://";
			}
			const url = protocol + HOST + "/ws?access_token=" + access_token;
			webSocketStore.connect(url);
		}
	},
);

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

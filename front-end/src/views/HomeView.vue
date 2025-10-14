<script setup lang="ts">
import { onBeforeMount } from "vue";
import { useRouter } from "vue-router";
import { useUserInfoStore } from "@/stores/user";
import { useWsConnectionStore } from "@/stores/ws-connection";
import { _get } from "@/utils/fetch";
import { createWs, initWsConnection } from "@/utils/websocket";
import Header from "@/components/Header.vue";
import ContactBox from "@/components/ContactBox/ContactBox.vue";

const router = useRouter();
const userStore = useUserInfoStore();
const connectionStore = useWsConnectionStore();

onBeforeMount(() => {
	if (!userStore.info) {
		_get("/api/v1/me/info")
			.then((data) => {
				userStore.updateInfo(data);
				const conn = createWs();
				if (conn === null) {
					alert("Failed to connect WebSockets.");
					return;
				}
				initWsConnection(conn);
				connectionStore.connection = conn;
			})
			.catch((err) => {
				userStore.updateInfo(null);
				console.error(err);
				router.push("/auth/login");
			});
	}
});
</script>

<template>
	<template v-if="userStore.info">
		<Header />
		<main class="relative flex h-auto w-auto grow flex-col sm:flex-row">
			<ContactBox />
			<router-view />
		</main>
	</template>
</template>

<style scoped></style>

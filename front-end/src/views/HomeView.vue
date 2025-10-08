<script setup lang="ts">
import { useRouter } from "vue-router";
import { useUserInfoStore } from "@/stores/user";
import { onBeforeMount } from "vue";
import { _get } from "@/utils/fetch";
import Header from "@/components/Header.vue";
import ContactBox from "@/components/ContactBox/ContactBox.vue";

const router = useRouter();
const userStore = useUserInfoStore();

onBeforeMount(() => {
	if (userStore.info === null) {
		_get("/api/v1/me/info")
			.then((data) => {
				userStore.updateInfo(data);
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
		<main class="flex h-auto w-auto grow flex-row">
			<ContactBox />
			<router-view />
		</main>
	</template>
</template>

<style scoped></style>

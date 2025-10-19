<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { _get, ORIGIN } from "@/utils/fetch";
import { usePersonalStore } from "@/stores/personal";
import { LoginSchema, type LoginPayload } from "@/types/auth";
import IconApp from "@/components/icons/IconApp.vue";

const router = useRouter();
const personalStore = usePersonalStore();

const data = ref<LoginPayload>({
	user_email: "",
	password: "",
});
const validationErrors = ref<{ [key: string]: string }>({
	user_email: "",
	password: "",
});
const errMsg = ref("");

const handleSubmit = () => {
	errMsg.value = "";
	validationErrors.value = {
		user_email: "",
		password: "",
	};
	const result = LoginSchema.safeParse(data.value);
	if (result.error) {
		const err = JSON.parse(result.error.message);
		console.log(err);
		for (let e of err) {
			validationErrors.value[e.path[0]] = e.message;
		}
		return;
	}
	fetch(`${ORIGIN}/api/v1/auth/login`, {
		method: "POST",
		body: JSON.stringify(result.data),
	})
		.then(async (res) => {
			const data = await res.json();
			if (!res.ok) throw data;
			localStorage.setItem("access_token", data.access_token);
			localStorage.setItem("refresh_token", data.refresh_token);
		})
		.then(() => {
			_get("/api/v1/me/info").then((data) => {
				personalStore.info = data;
			});
		})
		.then(() => {
			router.push("/");
		})
		.catch((err) => {
			errMsg.value = err.error;
			console.error(err);
			return;
		});
};
</script>

<template>
	<form
		class="flex h-full w-full flex-col items-center bg-neutral-900 p-8 pt-10 sm:h-auto"
		@submit.prevent="handleSubmit">
		<router-link
			to="/"
			class="flex items-center justify-center">
			<IconApp class="h-auto w-10" />
			<h1 class="ml-2 text-3xl font-bold">Synth-Socket</h1>
		</router-link>
		<span class="mt-2 text-neutral-400">Welcome back!</span>
		<div class="mt-6 flex w-full flex-col">
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="validationErrors.user_email === '' ? '' : 'text-red-500'"
					for="email">
					Email*
				</label>
				<input
					class="h-9 px-2"
					:class="validationErrors.user_email === '' ? '' : 'border-red-500'"
					v-model="data.user_email"
					name="email" />
				<small
					v-if="validationErrors.user_email"
					class="text-xs text-red-500">
					{{ validationErrors.user_email }}
				</small>
			</div>
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="validationErrors.password === '' ? '' : 'text-red-500'"
					for="password">
					Password*
				</label>
				<input
					class="h-9 px-2"
					:class="validationErrors.password === '' ? '' : 'border-red-500'"
					v-model="data.password"
					name="password"
					type="password" />
				<small
					v-if="validationErrors.password"
					class="text-xs text-red-500">
					{{ validationErrors.password }}
				</small>
			</div>
			<div
				v-if="errMsg"
				class="flex w-full items-center justify-center bg-red-500/30 px-3 py-2">
				<span class="max-w-full break-all text-red-500">
					{{ errMsg }}
				</span>
			</div>
			<button
				class="mt-4 h-10 w-full bg-violet-500 hover:bg-violet-500/70"
				type="submit">
				Login
			</button>
		</div>
		<router-link
			to="/auth/signup"
			class="mt-8 text-sm text-neutral-400 hover:underline"
			>Don't have an account? Register here!</router-link
		>
	</form>
</template>

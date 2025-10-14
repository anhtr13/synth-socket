<script setup lang="ts">
import { useRouter } from "vue-router";
import { ref } from "vue";
import { useUserInfoStore } from "@/stores/user";
import { _get, ORIGIN } from "@/utils/fetch";
import { SignupSchema, type SignupPayload } from "@/types/auth";
import IconApp from "@/components/icons/IconApp.vue";

const router = useRouter();
const userStore = useUserInfoStore();

const data = ref<SignupPayload>({
	user_email: "",
	user_name: "",
	password: "",
});
const validationError = ref<{ [key: string]: string }>({
	user_email: "",
	user_name: "",
	password: "",
});
const errMsg = ref("");

const handleSubmit = () => {
	errMsg.value = "";
	validationError.value = {
		user_email: "",
		user_name: "",
		password: "",
	};
	const result = SignupSchema.safeParse(data.value);
	if (result.error) {
		const err = JSON.parse(result.error.message);
		for (let e of err) {
			validationError.value[e.path[0]] = e.message;
		}
		return;
	}
	fetch(`${ORIGIN}/api/v1/auth/signup`, {
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
				userStore.updateInfo(data);
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
		<span class="mt-2 text-neutral-400">Create an account</span>
		<div class="mt-6 flex w-full flex-col">
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="validationError.user_email === '' ? '' : 'text-red-500'"
					for="email">
					Email*
				</label>
				<input
					class="h-9 px-2"
					:class="validationError.user_email === '' ? '' : 'border-red-500'"
					v-model="data.user_email"
					name="email" />
				<small
					v-if="validationError.user_email"
					class="text-xs text-red-500">
					{{ validationError.user_email }}
				</small>
			</div>
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="validationError.user_name === '' ? '' : 'text-red-500'"
					for="user_name">
					User name*
				</label>
				<input
					class="h-9 px-2"
					:class="validationError.user_name === '' ? '' : 'border-red-500'"
					v-model="data.user_name"
					name="user_name" />
				<small
					v-if="validationError.user_name"
					class="text-xs text-red-500">
					{{ validationError.user_name }}
				</small>
			</div>
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="validationError.password === '' ? '' : 'text-red-500'"
					for="password">
					Password*
				</label>
				<input
					class="h-9 px-2"
					:class="validationError.password === '' ? '' : 'border-red-500'"
					v-model="data.password"
					name="password"
					type="password" />
				<small
					v-if="validationError.password"
					class="text-xs text-red-500">
					{{ validationError.password }}
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
				Sign up
			</button>
		</div>
		<router-link
			to="/auth/login"
			class="mt-8 text-sm text-neutral-400 hover:underline">
			Already have an account? Login here!
		</router-link>
	</form>
</template>

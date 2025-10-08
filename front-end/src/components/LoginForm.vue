<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { LoginSchema, type LoginPayload } from "@/types/auth";
import { BACKEND_URL } from "@/utils/fetch";
import IconApp from "@/components/icons/IconApp.vue";

const router = useRouter();

const data = ref<LoginPayload>({
	user_email: "",
	password: "",
});

const error_msg = ref<{ [key: string]: string }>({
	user_email: "",
	password: "",
});

const handleSubmit = () => {
	error_msg.value = {
		user_email: "",
		password: "",
	};
	const result = LoginSchema.safeParse(data.value);
	if (result.error) {
		const err = JSON.parse(result.error.message);
		console.log(err);
		for (let e of err) {
			error_msg.value[e.path[0]] = e.message;
		}
	}
	if (result.data) {
		fetch(`${BACKEND_URL}/api/v1/auth/login`, {
			method: "POST",
			body: JSON.stringify(result.data),
		})
			.then(async (res) => {
				const data = await res.json();
				if (!res.ok) throw data;
				console.log(data);
				localStorage.setItem("access_token", data.access_token);
				localStorage.setItem("refresh_token", data.refresh_token);
				router.push("/");
			})
			.catch((err) => {
				console.error(err);
			});
	}
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
					:class="error_msg.user_email === '' ? '' : 'text-red-500'"
					for="email">
					Email*
				</label>
				<input
					class="h-9 px-2"
					:class="error_msg.user_email === '' ? '' : 'border-red-500'"
					v-model="data.user_email"
					name="email" />
				<small
					v-if="error_msg.user_email"
					class="text-xs text-red-500">
					{{ error_msg.user_email }}
				</small>
			</div>
			<div class="flex flex-col justify-start gap-y-1 pb-4">
				<label
					:class="error_msg.password === '' ? '' : 'text-red-500'"
					for="password">
					Password*
				</label>
				<input
					class="h-9 px-2"
					:class="error_msg.password === '' ? '' : 'border-red-500'"
					v-model="data.password"
					name="password"
					type="password" />
				<small
					v-if="error_msg.password"
					class="text-xs text-red-500">
					{{ error_msg.password }}
				</small>
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

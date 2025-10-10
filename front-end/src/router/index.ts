import { createRouter, createWebHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import { useUserInfoStore } from "@/stores/user";
import { _get } from "@/utils/fetch";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/",
			name: "home",
			component: HomeView,
			meta: { requireAuth: true },
			children: [
				{ path: "room", component: () => import("@/views/RoomView.vue") },
				{ path: "user", name: "user", component: () => import("@/views/UserView.vue") },
			],
		},
		{
			path: "/auth",
			name: "auth",
			component: () => import("@/views/AuthView.vue"),
			children: [
				{
					path: "login",
					name: "login",
					component: () => import("@/components/LoginForm.vue"),
				},
				{
					path: "signup",
					name: "signup",
					component: () => import("@/components/SignupForm.vue"),
				},
			],
		},
	],
});

router.beforeEach(async (to, from, next) => {
	const userStore = useUserInfoStore();

	if (to.meta.requiresAuth && userStore.info === null) {
		// If the route requires authentication and the user is not authenticated, redirect to login
		next({ name: "login" });
	} else {
		// Otherwise, allow navigation to proceed
		next();
	}
});

export default router;

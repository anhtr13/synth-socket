import { createRouter, createWebHistory } from "vue-router";
import HomeView from "@/views/HomeView.vue";
import Room from "@/components/ChatRoom/ChatRoom.vue";

const router = createRouter({
	history: createWebHistory(import.meta.env.BASE_URL),
	routes: [
		{
			path: "/",
			name: "home",
			component: HomeView,
			meta: { requireAuth: true },
			children: [
				{
					path: "room/:room_id",
					name: "room_id",
					component: Room,
				},
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
	const refresh_token = localStorage.getItem("refresh_token");

	if (to.meta.requireAuth && !refresh_token) {
		next({ name: "login" });
		return;
	}

	next();
});

export default router;

import "@/index.css";

import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";
import Vue3Toastify, { type ToastContainerOptions } from "vue3-toastify";
import "vue3-toastify/dist/index.css";

const app = createApp(App);

app.directive("clickoutside", {
	mounted: function (element, binding) {
		element.__handleClickOutside__ = (event: MouseEvent) => {
			if (!(element === event.target || element.contains(event.target))) {
				binding.value();
			}
		};
		document.addEventListener("click", element.__handleClickOutside__);
	},
	unmounted: function (element) {
		document.removeEventListener("click", element.__handleClickOutside__);
	},
});

app.use(createPinia());
app.use(router);
app.use(Vue3Toastify, {
	autoClose: 2000,
	theme: "dark",
} as ToastContainerOptions);

app.mount("#app");

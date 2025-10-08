import "@/index.css";

import { createApp } from "vue";
import { createPinia } from "pinia";

import App from "./App.vue";
import router from "./router";

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

app.mount("#app");

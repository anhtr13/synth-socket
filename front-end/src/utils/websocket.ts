export const HOST =
	import.meta.env.MODE === "development" ? "localhost:3000" : window.location.host;

export function createWs(): WebSocket | null {
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
	return new WebSocket(protocol + HOST + "/ws?access_token=" + access_token);
}

export function initWsConnection(conn: WebSocket) {
	conn.onopen = () => {
		console.log("websocket connected");
	};
	conn.onclose = function (event) {
		alert("Websocket closed");
		console.log("Connection closed:", event);
	};
	conn.onerror = function (event) {
		console.log("Connection error:", event);
	};
	conn.onmessage = function (event) {
		const payload = JSON.parse(event.data);
		console.log(payload);
	};
}

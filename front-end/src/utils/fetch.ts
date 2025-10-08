export const BACKEND_URL = "http://localhost:3000";

type Methods = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
type Options = {
	body?: any;
	searchParams?: { [key: string]: string };
};

/*<---request configuration--->*/
const request = (path: string, method: Methods, authorization: string, options?: Options) => {
	const url = new URL(path, BACKEND_URL);
	if (options?.searchParams) {
		for (let key of Object.keys(options?.searchParams))
			url.searchParams.set(key, options?.searchParams[key] || "");
	}
	return fetch(url, {
		method,
		// credentials: "include", // allows cross-origin requests to include cookies and authentication headers.
		headers: {
			accept: "application/json",
			"content-type": "application/json",
			authorization: authorization,
		},
		...(method !== "GET" &&
			method !== "DELETE" &&
			options?.body && { body: JSON.stringify(options.body) }),
	});
};

const getAccessToken = async (refresh_token: string) => {
	const res = await fetch(
		`${BACKEND_URL}/api/v1/auth/access_token?refresh_token=${refresh_token}`,
		{ method: "POST" },
	);
	const data: { access_token: string } = await res.json();
	if (!res.ok || !data.access_token) {
		throw data;
	}
	return data.access_token;
};

/*<---handle token rotation--->*/
const handleTokenRotation = async (path: string, method: Methods, options?: Options) => {
	let access_token = localStorage.getItem("access_token");
	let refresh_token = localStorage.getItem("refresh_token");
	if (!access_token) {
		if (!refresh_token) {
			throw new Error("No access_token or refresh_token founded");
		}
		access_token = await getAccessToken(refresh_token);
		localStorage.setItem("access_token", access_token);
	}

	let res = await request(path, method, access_token, options);
	let data = await res.json();
	if (res.ok) {
		return data;
	}

	if (res.status !== 401) throw data;
	if (!refresh_token) throw data;

	access_token = await getAccessToken(refresh_token);
	localStorage.setItem("access_token", access_token);

	res = await request(path, method, access_token, options);
	data = await res.json();

	if (!res.ok) {
		throw data;
	}
	return data;
};

/*<---client fetch funcion--->*/
export const _get = async (path: string, options?: Options) =>
	handleTokenRotation(path, "GET", options);
export const _post = (path: string, options?: Options) =>
	handleTokenRotation(path, "POST", options);
export const _patch = (path: string, options?: Options) =>
	handleTokenRotation(path, "PATCH", options);
export const _put = (path: string, options?: Options) =>
	handleTokenRotation(path, "PUT", options);
export const _delete = (path: string, options?: Options) =>
	handleTokenRotation(path, "DELETE", options);

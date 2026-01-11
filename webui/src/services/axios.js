import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

export function setAuthToken(token) {
	if (token) {
		instance.defaults.headers.common['Authorization'] = `Bearer ${token}`;
	} else {
		delete instance.defaults.headers.common['Authorization'];
	}
}


try {
	const token = localStorage.getItem('token');
	if (token) setAuthToken(token);
} catch (e) {
	// ignore (e.g. server-side rendering or unavailable localStorage)
}

export default instance;

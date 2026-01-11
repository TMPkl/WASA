import axios, { setAuthToken } from './axios.js';

export async function login(username) {
  const res = await axios.post('/login', { username });
  const token = res.data && res.data.token;
  if (!token) throw new Error('No token returned from server');
  try { localStorage.setItem('token', token); } catch (e) {}
  setAuthToken(token);
  return token;
}

export function logout() {
  try { localStorage.removeItem('token'); } catch (e) {}
  setAuthToken(null);
}

export function isLoggedIn() {
  try { return !!localStorage.getItem('token'); } catch (e) { return false; }
}

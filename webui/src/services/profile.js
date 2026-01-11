import axios, { setAuthToken } from './axios.js'

function decodeToken(token) {
  try {
    const payload = token.split('.')[1]
    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch (e) {
    return null
  }
}
                                         ///trying my best not to implement get profile, using token to extract username
export async function getProfile() {
  try {
    const token = localStorage.getItem('token')
    if (!token) return null
    const claims = decodeToken(token)
    return { username: claims && claims.sub ? claims.sub : null }
  } catch (e) {
    return null
  }
}

export async function updateUsername(newUsername) {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('Not authenticated')
  const claims = decodeToken(token)
  const current = claims && claims.sub ? claims.sub : ''
  const payload = { username: current, 'new-username': newUsername }
  const res = await axios.patch('/me/username', payload)
  if (res.data && res.data.token) {
    const newToken = res.data.token
    try { localStorage.setItem('token', newToken) } catch (e) {}
    setAuthToken(newToken)
  }
  return res.data
}

export async function uploadPhoto(file) {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('Not authenticated')
  const claims = decodeToken(token)
  const current = claims && claims.sub ? claims.sub : ''
  const fd = new FormData()
  fd.append('username', current)
  fd.append('photo', file)
  const res = await axios.post('/me/photo', fd, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return res.data
}
export async function getUserPhoto(username) {
  const res = await axios.get(`/users/${encodeURIComponent(username)}/photo`, {
    responseType: 'blob'
  })
  return res.data
}
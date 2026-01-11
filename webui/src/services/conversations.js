import axios from './axios.js'

function decodeToken(token) {
  try {
    const payload = token.split('.')[1]
    const json = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(json)
  } catch (e) {
    return null
  }
}

export async function getConversations() {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('Not authenticated')
  const claims = decodeToken(token)
  const username = claims && claims.sub ? claims.sub : ''
  if (!username) throw new Error('No username in token')
  
  const res = await axios.get(`/conversations/${username}`)
  return res.data
}

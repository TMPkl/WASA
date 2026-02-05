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

export async function getMessageAttachments(messageId) {
  const res = await axios.get(`/messages/${messageId}/attachments`)
  return res.data
}

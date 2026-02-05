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

export async function setGroupName(groupId, newName) {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('Not authenticated')
  const claims = decodeToken(token)
  const username = claims && claims.sub ? claims.sub : ''
  
  const res = await axios.patch(`/groups/${groupId}/name`, {
    username: username,
    new_name: newName
  })
  return res.data
}

export async function setGroupPhoto(groupId, photoFile) {
  const token = localStorage.getItem('token')
  if (!token) throw new Error('Not authenticated')
  const claims = decodeToken(token)
  const username = claims && claims.sub ? claims.sub : ''
  
  const formData = new FormData()
  formData.append('username', username)
  formData.append('photo', photoFile)
  
  const res = await axios.post(`/groups/${groupId}/photo`, formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return res.data
}

export async function getGroupPhoto(groupId) {
  const res = await axios.get(`/groups/${groupId}/photo`, {
    responseType: 'blob'
  })
  return res.data
}

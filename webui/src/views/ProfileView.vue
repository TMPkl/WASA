<script>
import { getProfile, updateUsername, uploadPhoto } from '../services/profile.js'
import { logout } from '../services/auth.js'

export default {
  data() {
    return {
      username: '',
      newUsername: '',
      errormsg: null,
      loading: false,
      photoFile: null,
      success: null
    }
  },
  methods: {
    async load() {
      this.loading = true
      this.errormsg = null
      try {
        const p = await getProfile()
        this.username = p && p.username ? p.username : ''
      } catch (e) {
        this.errormsg = String(e)
      }
      this.loading = false
    },
    async changeUsername() {
      if (!this.newUsername) return
      this.loading = true
      this.errormsg = null
      try {
        await updateUsername(this.newUsername)
        this.success = 'Username updated'
        await this.load()
      } catch (e) {
        this.errormsg = (e.response && e.response.data) || e.message || String(e)
      }
      this.loading = false
    },
    handleFile(e) {
      this.photoFile = e.target.files[0]
    },
    async submitPhoto() {
      if (!this.photoFile) return
      this.loading = true
      this.errormsg = null
      try {
        await uploadPhoto(this.photoFile)
        this.success = 'Photo uploaded'
      } catch (e) {
        this.errormsg = (e.response && e.response.data) || e.message || String(e)
      }
      this.loading = false
    },
    doLogout() {
      logout()
      this.$router.replace('/login')
    }
  },
  mounted() { this.load() }
}
</script>

<template>
  <div class="profile-view mt-4">
    <div class="card p-3" style="max-width:700px">
      <h3>My Profile</h3>
      <div v-if="errormsg" class="alert alert-danger">{{ errormsg }}</div>
      <div v-if="success" class="alert alert-success">{{ success }}</div>

      <div class="mb-3">
        <label class="form-label">Current username</label>
        <input class="form-control" :value="username" disabled />
      </div>

      <div class="mb-3">
        <label class="form-label">Change username</label>
        <div class="d-flex gap-2">
          <input v-model="newUsername" class="form-control" placeholder="new username" />
          <button class="btn btn-primary" :disabled="!newUsername || loading" @click="changeUsername">Change</button>
        </div>
      </div>

      <div class="mb-3">
        <label class="form-label">Profile photo</label>
        <input type="file" accept="image/*" @change="handleFile" />
        <div class="mt-2">
          <button class="btn btn-secondary" :disabled="!photoFile || loading" @click="submitPhoto">Upload photo</button>
        </div>
      </div>

      <div class="mt-3">
        <button class="btn btn-outline-danger" @click="doLogout">Logout</button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.profile-view .card { font-size: 1rem }
</style>

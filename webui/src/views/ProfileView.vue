<script>
import { getProfile, updateUsername, uploadPhoto, getUserPhoto } from '../services/profile.js'
import { logout } from '../services/auth.js'
import defaultPhoto from '@/assets/account_icon_400x400.png'
import UserPhoto from '@/components/UserPhoto.vue'


export default {
    components: { UserPhoto },
  data() {
    return {
      username: '',
      userPhotoUrl: null,
      defaultPhoto,
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
        // load user photo blob asynchronously
        await this.loadUserPhoto()
      } catch (e) {
        this.errormsg = String(e)
      }
      this.loading = false
    },
    async loadUserPhoto() {
      // clear previous
      if (this.userPhotoUrl) {
        try { URL.revokeObjectURL(this.userPhotoUrl) } catch (e) {}
        this.userPhotoUrl = null
      }
      if (!this.username) return
      try {
        const blob = await getUserPhoto(this.username)
        if (blob) this.userPhotoUrl = URL.createObjectURL(blob)
      } catch (e) {
        // ignore missing photo
      }
    },
    async changeUsername() {
      if (!this.newUsername) return
      this.loading = true
      this.errormsg = null
      try {
        await updateUsername(this.newUsername)
        localStorage.setItem('username', this.newUsername)
        this.success = 'Username updated'
        this.newUsername = ''
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
        await this.loadUserPhoto()
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
          <button class="btn btn-primary" :disabled="!newUsername || loading" @click="changeUsername">Change1</button>
        </div>
      </div>

      <div class="mb-3">
        <div class="mb-3">
  <label class="form-label">Profile photo</label>

  <div class="input-group">
    <input
      type="file"
      class="form-control"
      accept="image/*"
      @change="handleFile"
    />
    <button
      class="btn btn-secondary"
      :disabled="!photoFile || loading"
      @click="submitPhoto"
    >
      Upload
    </button>
  </div>

  <small v-if="photoFile" class="text-muted">
    Selected: {{ photoFile.name }}
  </small>
</div>

      </div>
      <div class="photo-label ">
        <label class="">My current photo</label>
        <div class="photo-wrapper">
          <UserPhoto :username="username" :size="200" />
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
.photo-label { display: block; font-size: 2rem;
            font-weight: bold;
            justify-content: center;
            
            }

.photo-wrapper {
  display: flex;
  justify-content: center;
  align-items: center;
}

.photo-wrapper img,
.photo-wrapper > * {
  max-width: 200px;
  max-height: 200px;
  width: 100%;
  height: auto;
}
</style>

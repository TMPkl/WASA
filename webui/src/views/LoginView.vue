<script>
import { login } from '../services/auth.js'
export default {
  data() {
    return {
      username: '',
      errormsg: null,
      loading: false
    }
  },
  methods: {
    async submit() {
      this.errormsg = null;
      this.loading = true;
      try {
        await login(this.username);
        this.$router.push('/profile');
      } catch (e) {
        this.errormsg = (e.response && e.response.data) || e.message || String(e);
      }
      this.loading = false;
    }
  }
}
</script>

<template>
  <div class="container mt-4 mx-auto" style="max-width:420px">
    <div class="card p-3">
      <h3 class="mb-3">Log In</h3>
      <ErrorMsg v-if="errormsg" :msg="errormsg" />
      <div class="mb-3">
        <label class="form-label">Username</label>
        <input v-model="username" class="form-control" placeholder="username" />
      </div>
      <div class="d-grid">
        <button class="btn btn-primary" :disabled="loading || !username" @click="submit">
          <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
          Log In
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
</style>

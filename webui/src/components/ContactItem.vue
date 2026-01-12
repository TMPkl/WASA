<script>
import axios from 'axios';

export default {
  name: 'ContactItem',
  props: {
    title: String,
    lastMessage: String,
    time: String,
    status: String,
    isPrivate: Boolean,
    otherUsername: String
  },
  data() {
    return {
      userPhoto: null
    };
  },
  async mounted() {
    if (this.isPrivate && this.otherUsername) {
      await this.loadUserPhoto();
    }
  },
  unmounted() {
    if (this.userPhoto) {
      URL.revokeObjectURL(this.userPhoto);
    }
  },
  watch: {
    otherUsername(newVal) {
      if (this.userPhoto) {
        URL.revokeObjectURL(this.userPhoto);
        this.userPhoto = null;
      }
      if (this.isPrivate && newVal) {
        this.loadUserPhoto();
      }
    }
  },
  methods: {
    async loadUserPhoto() {
      try {
        const token = localStorage.getItem('token');
        const res = await axios({
          method: 'get',
          url: `${__API_URL__}/users/${this.otherUsername}/photo?t=${Date.now()}`,
          headers: {
            Authorization: `Bearer ${token}`
          },
          responseType: 'blob'
        });
        const url = URL.createObjectURL(res.data);
        this.userPhoto = url;
      } catch (e) {
        console.error('Failed to load user photo:', e);
      }
    }
  }
}
</script>

<template>
  <div class="d-flex p-3 text-muted border-bottom align-items-center justify-content-between">
    <div class="photo">
      <img :src="userPhoto || '/account_icon_400x400.png'" alt="Avatar" class="rounded-circle" width="100" height="100" style="object-fit: cover;" />
    </div>
    <div class="p-3 bg-light flex-grow-1 rounded-3">
      <div class="fw-bold">{{ title }}</div>
      <div class="text-dark">{{ lastMessage }}</div>
    </div>
    <div class="d-flex flex-column bg-light rounded-3">
      <div class="p-2 text-end">
        <div class="text-muted">{{ time }}</div>
      </div>
      <div class="p-2">
        <div class="fw-bold">{{ status }}</div>
      </div>
    </div>
  </div>
</template>

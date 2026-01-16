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
    otherUsername: String,
    groupId: Number
  },
  data() {
    return {
      photoUrl: null
    };
  },
  async mounted() {
    await this.loadPhoto();
  },
  unmounted() {
    if (this.photoUrl) {
      URL.revokeObjectURL(this.photoUrl);
    }
  },
  watch: {
    otherUsername() {
      this.reloadPhoto();
    },
    groupId() {
      this.reloadPhoto();
    }
  },
  methods: {
    async reloadPhoto() {
      if (this.photoUrl) {
        URL.revokeObjectURL(this.photoUrl);
        this.photoUrl = null;
      }
      await this.loadPhoto();
    },
    async loadPhoto() {
      if (this.isPrivate && this.otherUsername) {
        await this.loadUserPhoto();
      } else if (!this.isPrivate && this.groupId) {
        await this.loadGroupPhoto();
      }
    },
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
        this.photoUrl = URL.createObjectURL(res.data);
      } catch (e) {
        console.error('Failed to load user photo:', e);
      }
    },
    async loadGroupPhoto() {
      try {
        const token = localStorage.getItem('token');
        const res = await axios({
          method: 'get',
          url: `${__API_URL__}/groups/${this.groupId}/photo?t=${Date.now()}`,
          headers: {
            Authorization: `Bearer ${token}`
          },
          responseType: 'blob'
        });
        this.photoUrl = URL.createObjectURL(res.data);
      } catch (e) {
        console.error('Failed to load group photo:', e);
      }
    }
  }
}
</script>

<template>
  <div class="d-flex p-3 text-muted border-bottom align-items-center justify-content-between">
    <div class="photo">
      <img :src="photoUrl || '/account_icon_400x400.png'" alt="Avatar" class="rounded-circle" width="100" height="100" style="object-fit: cover;" />
    </div>
    <div class="p-3 bg-light flex-grow-1 rounded-3">
      <div class="text-muted small">{{ isPrivate ? 'Private chat with:' : 'Group:' }}</div>
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

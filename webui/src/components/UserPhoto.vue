<script>
import defaultPhoto from '@/assets/account_icon_400x400.png'
import { getUserPhoto } from '@/services/profile.js'

export default {
  name: 'UserPhoto',
  props: {
    username: {
      type: String,
      required: true
    },
    size: {
      type: Number,
      default: 200
    }
  },
  data() {
    return {
      photoUrl: null,
      defaultPhoto
    }
  },
  watch: {
    username: {
      immediate: true,
      handler() {
        this.loadPhoto()
      }
    }
  },
  methods: {
    async loadPhoto() {
      this.clearUrl()
      if (!this.username) return

      try {
        const blob = await getUserPhoto(this.username)
        if (blob) {
          this.photoUrl = URL.createObjectURL(blob)
        }
      } catch (_) {
        // brak zdjęcia = fallback
      }
    },
    clearUrl() {
      if (this.photoUrl) {
        try { URL.revokeObjectURL(this.photoUrl) } catch (_) {}
        this.photoUrl = null
      }
    }
  },
  unmounted() {
    this.clearUrl()
  }
}
</script>

<template>
  <img
    :src="photoUrl || defaultPhoto"
    :alt="`Photo of ${username}`"
    :style="{ maxWidth: size + 'px', maxHeight: size + 'px' }"
  />
</template>

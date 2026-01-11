<script>
import defaultPhoto from '../assets/account_icon_400x400.png';
import UserPhoto from './UserPhoto.vue';
export default {
  name: 'ContactItem',
  props: {
    title: String,
    lastMessage: String,
    time: String,
    status: String,
  }
}
</script>

<template>
  <div class="d-flex p-3  text-muted border-bottom align-items-center justify-content-between">
    <div class="photo">
      <img src="/account_icon_400x400.png" alt="Avatar" class="rounded-circle" width="100" height="100" />
    </div>
    <div class="p-3 bg-light  flex-grow-1 rounded-3">
    <div class="fw-bold">{{ title }}</div>
    <div class="text-dark">{{ lastMessage }}</div>
  </div>
    <div class="d-flex flex-column bg-light rounded-3">
      <div class="p-2 text-end">
        <div class="text-muted">{{ time }}</div>
        
      </div>
      <div class="p-2 ">
        <div class="fw-bold">{{ status }}</div>
      </div>
    </div>
</div>
</template>

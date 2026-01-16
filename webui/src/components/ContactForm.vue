<script>
import ContactItem from './ContactItem.vue';
export default {
  name: 'ContactForm',
  components: {
    ContactItem
  },
  props: {
    conversations: {
      type: Array,
      default: () => []
    }
  },
  methods: {
    goToConversation(id) {
      this.$router.push({ path: `/conversation/${id}` });
    }
  }
}
</script>

<template>
  <div class="list-group">
    <ContactItem
      v-for="conversation in conversations"
      :key="conversation.id"
      :title="conversation.title"
      :lastMessage="conversation.lastMessage"
      :time="conversation.time"
      :status="conversation.status"
      :is-private="conversation.isPrivate"
      :other-username="conversation.otherUsername"
      :group-id="conversation.groupId"
      @click="goToConversation(conversation.id)"
      style="cursor:pointer;"
    />
  </div>
</template>
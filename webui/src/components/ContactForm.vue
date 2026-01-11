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
    />
  </div>
</template>
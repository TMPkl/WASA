<script>
import MessageView from './MessageView.vue';

export default {
  name: 'ConversationMessagesView',
  components: {
    MessageView
  },
  data() {
    return {
    };
  }
}
</script>

<template>
  <div class="conversation-messages-view p-2">
    <h2>TUTAJ PISZEMY NAZWE GRUPY ALBO USERNAME</h2>
    <MessageView :messages="messages" />
  </div>
</template>

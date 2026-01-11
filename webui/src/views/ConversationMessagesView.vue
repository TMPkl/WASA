<script>
import MessageView from './MessageView.vue';
import axios from 'axios';

export default {
  name: 'ConversationMessagesView',
  components: {
    MessageView
  },
  props: ['id'],
  data() {
    return {
      messages: [],
      loading: false,
      error: null,
      conversationName: '',
      isGroup: false,
    };
  },
  async mounted() {
    this.loading = true;
    try {
      const token = localStorage.getItem('token');
      if (!this.id) {
        this.error = 'No conversation ID ';
        return;
      }
      const username = localStorage.getItem('username');
      const res = await axios({
        method: 'post',
        url: `${__API_URL__}/conversation/${this.id}`,
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        data: {
          username: username,
          message_quantity: 50
        }
      });
      console.log('API response:', res.data);
      //{ conversation_id, is_group, participants, messages: [ { id, sender_username, content, timestamp, has_attachment } ] }
      this.isGroup = res.data.is_group || false;
      this.messages = res.data.messages.map(msg => ({
        id: msg.id,
        sender: msg.sender_username,
        content: msg.content,
        timestamp: msg.timestamp,
        attachment: msg.has_attachment
        
      }));
      if (this.isGroup) {
        this.conversationName = res.data.participants?.join(', ') || `Group ${this.id}`;
      } else {
        //for priv show the other 
        const currentUser = localStorage.getItem('username');
        const otherUser = res.data.participants?.find(p => p !== currentUser);
        this.conversationName = otherUser || res.data.participants?.[0] || `Conversation ${this.id}`;
      }
    } catch (e) {
      this.error = e?.response?.data?.error || e.message || 'Error loading messages';
    } finally {
      this.loading = false;
    }
  },
}
</script>

<template>
  <div class="conversation-messages-view d-flex flex-column p-2">
    
    <div class="p-3 mb-2 bg-secondary text-white rounded-4 flex-shrink-0">
      <h2 class="m-0">{{ conversationName }}</h2>
    </div>

    <div class="messages flex-grow-1 mb-2">
      <div v-if="loading">Ładowanie wiadomości...</div>
      <div v-else-if="error" class="text-danger">{{ error }}</div>
      <MessageView v-else :messages="messages" />
    </div>

    <div class="flex-shrink-0">
      <form class="d-flex">
        <input
          type="text"
          class="form-control me-2"
          placeholder="Write a message..."
        />
        <button type="submit" class="btn btn-primary">Send</button>
        <button type="button" class="btn btn-secondary ms-2">Attach</button>
        <div v-if="isGroup" class="align-self-center d-flex">
            <button type="button" class="btn btn-info ms-2">Add User</button>
            <button type="button" class="btn btn-danger ms-2">Leave Group</button>
        </div>
      </form>
    </div>

  </div>
</template>
<style scoped>
    .conversation-messages-view {
  height: 90vh;      
  overflow: hidden;   
}

.messages {
  overflow-y: auto;    
  overflow-x: hidden;
}

</style>

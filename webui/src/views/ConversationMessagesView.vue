<script>
import MessageView from './MessageView.vue';
import MessageInputForm from '../components/MessageInputForm.vue';
import axios from 'axios';

export default {
  name: 'ConversationMessagesView',
  components: {
    MessageView,
    MessageInputForm
  },
  props: ['id'],
  data() {
    return {
      messages: [],
      loading: false,
      error: null,
      conversationName: '',
      isGroup: false,
      sending: false,
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
  methods: {
    async handleSendMessage({ content, files }) {
      this.sending = true;
      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        
        const formData = new FormData();
        formData.append('senderUsername', username);
        formData.append('content', content);
        
        // For group conversations, we need to handle differently
        // For now, we'll use the first other participant
        const currentUser = localStorage.getItem('username');
        const otherUser = this.isGroup 
          ? '' // Groups might need different handling
          : (await this.getOtherParticipant(currentUser));
        
        if (!this.isGroup) {
          formData.append('receiverUsername', otherUser);
        }
        
        // Add attachments
        for (const file of files) {
          formData.append('attachments', file);
        }

        await axios({
          method: 'post',
          url: `${__API_URL__}/messages`,
          headers: {
            Authorization: `Bearer ${token}`,
          },
          data: formData
        });

        // Clear form via ref
        this.$refs.messageForm.clearForm();
        
        // Reload messages
        await this.loadMessages();
      } catch (e) {
        this.error = e?.response?.data?.error || e.message || 'Failed to send message';
      } finally {
        this.sending = false;
      }
    },
    async getOtherParticipant(currentUser) {
      // This should be stored from the initial load
      const token = localStorage.getItem('token');
      const res = await axios({
        method: 'post',
        url: `${__API_URL__}/conversation/${this.id}`,
        headers: {
          Authorization: `Bearer ${token}`,
          'Content-Type': 'application/json'
        },
        data: {
          username: currentUser,
          message_quantity: 1
        }
      });
      return res.data.participants?.find(p => p !== currentUser) || '';
    },
    async loadMessages() {
      const token = localStorage.getItem('token');
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
      
      this.messages = res.data.messages.map(msg => ({
        id: msg.id,
        sender: msg.sender_username,
        content: msg.content,
        timestamp: msg.timestamp,
        attachment: msg.has_attachment
      }));
    }
  }
}
</script>

<template>
  <div class="conversation-messages-view d-flex flex-column p-2">
    
    <div class="p-3 mb-2 bg-secondary text-white rounded-4 flex-shrink-0">
      <h2 class="m-0">{{ conversationName }}</h2>
    </div>

    <div class="messages flex-grow-1 mb-2">
      <div v-if="loading">Loading...</div>
      <div v-else-if="error" class="text-danger">{{ error }}</div>
      <MessageView v-else :messages="messages" />
    </div>

    <MessageInputForm
      ref="messageForm"
      :is-group="isGroup"
      :sending="sending"
      @send-message="handleSendMessage"
    />

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

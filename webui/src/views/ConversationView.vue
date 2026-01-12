<script>
import ContactForm from '../components/ContactForm.vue';
import { getConversations } from '../services/conversations.js';

export default {
  name: 'ConversationView',
  components: {
    ContactForm
  },
  data() {
    return {
      conversations: [],
      loading: false,
      error: null
    }
  },
  async mounted() {
    await this.loadConversations();
  },
  methods: {
    async loadConversations() {
      this.loading = true;
      this.error = null;
      try {
        const data = await getConversations();
        this.conversations = data
          .map(conv => ({
            id: conv.ConversationID,
            title: conv.ConversationType === 'private' ? conv.OtherUsername : conv.GroupName,
            lastMessage: conv.LastMessage || '',
            time: this.formatTime(conv.LastMessageTime),
            timestamp: conv.LastMessageTime,
            status: conv.Status || ''
          }))
          .sort((a, b) => {
            if (!a.timestamp) return 1;
            if (!b.timestamp) return -1;
            return new Date(b.timestamp) - new Date(a.timestamp);
          });
      } catch (err) {
        this.error = err.message || 'Error loading conversations';
        console.error('Error loading conversations:', err);
      } finally {
        this.loading = false;
      }
    },
    formatTime(timestamp) {
      if (!timestamp) return '';
      const date = new Date(timestamp);
      const now = new Date();
      const diff = now - date;
      const days = Math.floor(diff / (1000 * 60 * 60 * 24));
      
      if (days === 0) {
        return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
      } else if (days === 1) {
        return 'yesterday';
      } else if (days < 7) {
        return `${days} days ago`;
      } else {
        return date.toLocaleDateString('en-US');
      }
    }
  }
}
</script>

<template>
  <div style ="height: 40px;"> 
    <!-- i KNOW it is bad idea to use inline styles for spaing but this is a quik fix -->
  </div>
  <div class="conversation-view">
    <div class="d-flex align-items-center mb-3">
      <h2>Conversations</h2>
      <button class="btn btn-primary ms-auto" @click="loadConversations" :disabled="loading">
        Create New Conversation
      </button>
    </div>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <ContactForm v-else :conversations="conversations" />
  </div>
</template>

<script>
import MessageView from './MessageView.vue';
import MessageInputForm from '../components/MessageInputForm.vue';
import UserList from '../components/UserList.vue';
import axios from 'axios';

export default {
  name: 'ConversationMessagesView',
  components: {
    MessageView,
    MessageInputForm
    ,UserList
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
      groupId: null,
      refreshInterval: null,
      otherUsername: null,
      otherUserPhoto: null,
      users: [],
      selectedUser: null,
      replyingTo: null,
      participants: [],
    };
  },
  async mounted() {
    this.loading = true;
    try {
      await this.loadConversationData();
    } catch (e) {
      this.error = e?.response?.data?.error || e.message || 'Error loading messages';
    } finally {
      this.loading = false;
    }

    this.scrollToBottom();
  },
  unmounted() {
    this.stopAutoRefresh();
  },
  methods: {
    async loadConversationData() {
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
      this.groupId = res.data.group_id || null;
      this.participants = res.data.participants || [];
      this.messages = (res.data.messages || []).map(msg => ({
        id: msg.id,
        sender: msg.sender_username,
        content: msg.content,
        timestamp: msg.timestamp,
        attachment: msg.has_attachment,
        reactions: msg.reactions || [],
        status: msg.status || '',
        replyingToId: msg.replying_to_id || null,
        replyingToSender: msg.replying_to_sender || null,
        replyingToContent: msg.replying_to_content || null
      }));
      if (this.isGroup) {
        this.conversationName = res.data.participants?.join(', ') || `Group ${this.id}`;
      } else {
        const currentUser = localStorage.getItem('username');
        const otherUser = res.data.participants?.find(p => p !== currentUser);
        this.otherUsername = otherUser;
        this.conversationName = otherUser || res.data.participants?.[0] || `Conversation ${this.id}`;
        if (otherUser) {
          this.loadUserPhoto(otherUser);
        }
      }

      // Update status to "received" for messages sent by others
      await this.updateReceivedMessagesStatus(username);
    },
    async handleSendMessage({ content, files, replyingToId }) {
      this.sending = true;
      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        
        const formData = new FormData();
        formData.append('senderUsername', username);
        formData.append('content', content);
        formData.append('conversationId', this.id);
        
        if (replyingToId) {
          formData.append('replyingToId', replyingToId);
        }
        
        const currentUser = localStorage.getItem('username');
        const otherUser = this.isGroup 
          ? ''
          : (await this.getOtherParticipant(currentUser));
        
        if (!this.isGroup) {
          formData.append('receiverUsername', otherUser);
        }
        
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

        this.$refs.messageForm.clearForm();
        this.replyingTo = null;
        
        await this.loadMessages();
      } catch (e) {
        this.error = e?.response?.data?.error || e.message || 'Failed to send message';
      } finally {
        this.sending = false;
      }
    },
    async getOtherParticipant(currentUser) {
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
      
      this.messages = (res.data.messages || []).map(msg => ({
        id: msg.id,
        sender: msg.sender_username,
        content: msg.content,
        timestamp: msg.timestamp,
        attachment: msg.has_attachment,
        reactions: msg.reactions || [],
        status: msg.status || '',
        replyingToId: msg.replying_to_id || null,
        replyingToSender: msg.replying_to_sender || null,
        replyingToContent: msg.replying_to_content || null
      }));
      this.scrollToBottom();
    },
    handleMessageDeleted(messageId) {
      this.messages = this.messages.filter(msg => msg.id !== messageId);
    },
    handleReplyMessage(messageId) {
      const message = this.messages.find(msg => msg.id === messageId);
      if (message) {
        this.replyingTo = {
          id: messageId,
          sender: message.sender,
          content: message.content
        };
        this.$nextTick(() => {
          if (this.$refs.messageForm && this.$refs.messageForm.$el) {
            this.$refs.messageForm.$el.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
          }
        });
      }
    },
    handleCancelReply() {
      this.replyingTo = null;
    }
    ,
    openAddUserModal() {
      if (!this.isGroup) return;
      this.$nextTick(() => {
        if (this.$refs.userList && this.$refs.userList.open) this.$refs.userList.open(false, this.participants);
      });
    },
    handleUserSelect(username) {
      console.log('Selected username:', username);
      // parent may refresh participants or UI here if desired
    }
    ,
    async handleLeaveGroup() {
      if (!this.isGroup) return;
      try {
        const token = localStorage.getItem('token');
        let username = localStorage.getItem('username');
        if (!username && token) {
          try {
            const payload = JSON.parse(atob(token.split('.')[1].replace(/-/g, '+').replace(/_/g, '/')));
            username = payload.sub || payload.username || payload.name || '';
            if (username) localStorage.setItem('username', username);
          } catch (err) {
            console.warn('Failed to decode token for username', err);
          }
        }

        if (!username) {
          this.error = 'No username available to leave group';
          console.warn('leave group aborted: no username');
          return;
        }

        const targetGroupId = this.groupId || this.id;
        const url = `${__API_URL__}/groups/${targetGroupId}/members/me`;
        const headers = token ? { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' } : { 'Content-Type': 'application/json' };
        const payload = { username };
        console.log('Leaving group request', url, payload);
        // Use axios.delete with config to include request body
        await axios.delete(url, { headers, data: payload });

        this.$router.push({ path: '/' });
      } catch (e) {
        if (e.response) {
          console.warn('leave group failed', e.response.status, e.response.data);
          this.error = e.response.data?.error || JSON.stringify(e.response.data) || `Request failed: ${e.response.status}`;
        } else {
          this.error = e.message || 'Failed to leave group';
        }
      }
    },

    startAutoRefresh() {
      this.refreshInterval = setInterval(() => {
        this.loadMessages();
      }, 1000);
    },

    stopAutoRefresh() {
      if (this.refreshInterval) {
        clearInterval(this.refreshInterval);
        this.refreshInterval = null;
      }
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },

    async updateReceivedMessagesStatus(currentUsername) {
      const token = localStorage.getItem('token');
      
      // Find messages sent by others that need status update
      const messagesToUpdate = this.messages.filter(msg => 
        msg.sender !== currentUsername && 
        msg.status === 'sent'
      );

      // Update each message status to "received"
      for (const msg of messagesToUpdate) {
        try {
          await axios({
            method: 'patch',
            url: `${__API_URL__}/messages/${msg.id}/status`,
            headers: {
              Authorization: `Bearer ${token}`,
              'Content-Type': 'application/json'
            },
            data: {
              username: currentUsername,
              status: 'received'
            }
          });
          // Update local message status
          msg.status = 'received';
        } catch (e) {
          console.error(`Failed to update status for message ${msg.id}:`, e);
        }
      }
    },

    async loadUserPhoto(username) {
      try {
        const token = localStorage.getItem('token');
        const res = await axios({
          method: 'get',
          url: `${__API_URL__}/users/${username}/photo`,
          headers: {
            Authorization: `Bearer ${token}`
          },
          responseType: 'blob'
        });
        const url = URL.createObjectURL(res.data);
        this.otherUserPhoto = url;
      } catch (e) {
        console.error('Failed to load user photo:', e);
      }
    },

    async handleSelectUser(username, item) {
    console.log('Wybrano username:', username);
    console.log('Wybrany obiekt:', item);
    const my_username = localStorage.getItem('username');
    
    try {
      const response = await axios.post(`${__API_URL__}/groups/${this.groupId}/members`, {
        username: my_username,
        user_to_add: username
      }, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('token')}`,
          'Content-Type': 'application/json'
        }
      });
      console.log('dodany:', response.data);
      this.error = null;
      
      // Reload conversation data to update participants list
      await this.loadConversationData();
    } catch (error) {
      console.error('Błąd:', error.response ? error.response.data : error.message);
      this.error = error.response?.data || error.message;
    }

    this.selectedUser = item;
  }
  }
}
</script>

<template>
  <div class="conversation-messages-view d-flex flex-column p-2">
    
    <div class="p-3 mb-2 bg-secondary text-white rounded-4 flex-shrink-0 d-flex align-items-center">
      <div v-if="!isGroup && otherUserPhoto" class="me-3">
        <img :src="otherUserPhoto" alt="User photo" class="rounded-circle" style="width: 50px; height: 50px; object-fit: cover;">
      </div>
      <h2 class="m-0">{{ conversationName }}</h2>
    </div>

    <div class="messages flex-grow-1 mb-2" ref="messagesContainer">
      <div v-if="loading">Loading...</div>
      <div v-else-if="error" class="text-danger">{{ error }}</div>
      <MessageView v-else :messages="messages" @message-deleted="handleMessageDeleted" @reply-message="handleReplyMessage" />
    </div>

    <MessageInputForm
      ref="messageForm"
      :is-group="isGroup"
      :sending="sending"
      :replying-to="replyingTo"
      @send-message="handleSendMessage"
      @add-user="openAddUserModal"
      @leave-group="handleLeaveGroup"
      @cancel-reply="handleCancelReply"
    />

    <UserList ref="userList" :items="users" @select="handleSelectUser" />

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

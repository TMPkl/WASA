<template>
  <div>
    <teleport to="body">
      <div
        class="modal fade"
        :class="{ show: showModal }"
        style="display: block;"
        tabindex="-1"
        role="dialog"
        v-if="showModal"
      >
        <div class="modal-dialog modal-dialog-centered" role="document">
          <div class="modal-content">
            <div class="modal-header">
              <h5 class="modal-title">{{ title }}</h5>
              <button type="button" class="btn-close" @click="close"></button>
            </div>
            <div class="modal-body">
              <div class="d-flex flex-column gap-2">
                <button class="btn btn-success btn-sm mb-3" @click="openNewContactList">+ Add new contact</button>

                <label class="form-label small text-muted">Select existing conversation:</label>
                <input 
                  type="text" 
                  class="form-control mb-3" 
                  placeholder="Search conversations..." 
                  v-model="searchQuery" 
                />

                <div v-if="loading" class="text-center p-3">
                  <div class="spinner-border spinner-border-sm" role="status"></div>
                  Loading...
                </div>

                <div v-else-if="error" class="alert alert-danger p-2">{{ error }}</div>

                <ul v-else class="list-group overflow-auto" style="max-height: 400px; cursor: pointer;">
                  <li
                    v-for="conv in filteredConversations"
                    :key="conv.id"
                    class="list-group-item d-flex justify-content-between align-items-center"
                    @click="selectConversation(conv)"
                  >
                    <div>
                      <strong>{{ conv.title }}</strong>
                      <small v-if="conv.lastMessage" class="d-block text-muted text-truncate" style="max-width: 250px;">
                        {{ conv.lastMessage }}
                      </small>
                    </div>
                    <span class="badge bg-secondary">{{ conv.type }}</span>
                  </li>
                  <li v-if="filteredConversations.length === 0" class="list-group-item text-muted text-center">
                    No conversations
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div v-if="showModal" class="modal-backdrop fade show"></div>
    </teleport>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'ConversationListModal',
  props: {
    title: {
      type: String,
      default: 'Select conversation'
    }
  },
  data() {
    return {
      showModal: false,
      searchQuery: '',
      conversations: [],
      loading: false,
      error: null
    };
  },
  computed: {
    filteredConversations() {
      if (!this.searchQuery) return this.conversations;
      const q = this.searchQuery.toLowerCase();
      return this.conversations.filter(conv =>
        (conv.title || '').toLowerCase().includes(q) ||
        (conv.lastMessage || '').toLowerCase().includes(q)
      );
    }
  },
  methods: {
    async open() {
      this.error = null;
      this.loading = true;
      this.searchQuery = '';
      this.showModal = true;

      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');

        if (!username) {
          this.error = 'No username found';
          this.loading = false;
          return;
        }

        const res = await axios({
          method: 'get',
          url: `${__API_URL__}/conversations/${username}`,
          headers: {
            Authorization: token ? `Bearer ${token}` : undefined
          }
        });

        // Map response to uniform format
        // Expected: array of { ConversationID, ConversationType, GroupName, OtherUsername, LastMessage, ... }
        const data = Array.isArray(res.data) ? res.data : [];
        this.conversations = data.map(conv => ({
          id: conv.ConversationID,
          title: conv.ConversationType === 'private' ? conv.OtherUsername : conv.GroupName,
          type: conv.ConversationType === 'private' ? 'Private' : 'Group',
          lastMessage: conv.LastMessage || ''
        }));
      } catch (e) {
        this.error = e?.response?.data?.error || e.message || 'Error loading conversations';
        console.warn('ConversationListModal.open failed', e);
      } finally {
        this.loading = false;
      }
    },
    close() {
      this.showModal = false;
    },
    selectConversation(conv) {
      // Emit conversation id (integer) for forwarding
      this.$emit('select', conv.id, conv);
      this.close();
    },
    openNewContactList() {
      this.$emit('open-new-contact-list');
    }
  }
};
</script>

<script>
import axios from 'axios';

export default {
  name: 'NewConversationModal',
  data() {
    return {
      mode: null,
      users: [],
      selectedUser: null,
      groupName: '',
      initialMessage: '',
      loading: false,
      error: null
    };
  },
  methods: {
    open() {
      this.mode = null;
      this.selectedUser = null;
      this.groupName = '';
      this.initialMessage = '';
      this.error = null;
      this.users = [];
      this.$refs.modal.showModal();
    },
    async selectMode(selectedMode) {
      this.mode = selectedMode;
      if (selectedMode === 'private') {
        await this.loadUsers();
      }
    },
    async loadUsers() {
      this.loading = true;
      try {
        const token = localStorage.getItem('token');
        const res = await axios({
          method: 'get',
          url: `${__API_URL__}/users`,
          headers: {
            Authorization: `Bearer ${token}`
          }
        });
        this.users = res.data || [];
      } catch (e) {
        this.error = 'Failed to load users';
        console.error(e);
      } finally {
        this.loading = false;
      }
    },
    async createPrivateConversation() {
      if (!this.selectedUser) {
        this.error = 'Please select a user';
        return;
      }
      if (!this.initialMessage.trim()) {
        this.error = 'Please enter a message';
        return;
      }
      this.loading = true;
      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        const formData = new FormData();
        formData.append('senderUsername', username);
        formData.append('receiverUsername', this.selectedUser);
        formData.append('content', this.initialMessage.trim());

        const res = await axios({
          method: 'post',
          url: `${__API_URL__}/messages`,
          headers: {
            Authorization: `Bearer ${token}`
          },
          data: formData
        });
        this.$refs.modal.close();
        this.$emit('conversation-created', res.data.conversation_id);
      } catch (e) {
        this.error = e?.response?.data?.error || 'Failed to create conversation';
        console.error(e);
      } finally {
        this.loading = false;
      }
    },
    async createGroup() {
      if (!this.groupName.trim()) {
        this.error = 'Please enter a group name';
        return;
      }
      this.loading = true;
      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        const res = await axios({
          method: 'post',
          url: `${__API_URL__}/groups`,
          headers: {
            Authorization: `Bearer ${token}`,
            'Content-Type': 'application/json'
          },
          data: {
            username: username,
            group_name: this.groupName.trim()
          }
        });
        this.$refs.modal.close();
        this.$emit('group-created', res.data.group_id);
      } catch (e) {
        this.error = e?.response?.data?.error || 'Failed to create group';
        console.error(e);
      } finally {
        this.loading = false;
      }
    },
    cancel() {
      this.$refs.modal.close();
    }
  }
}
</script>

<template>
  <dialog ref="modal" class="modal" @keydown.escape="cancel">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Create New Conversation</h5>
          <button type="button" class="btn-close" @click="cancel" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <div v-if="error" class="alert alert-danger" role="alert">{{ error }}</div>

          <div v-if="!mode" class="d-grid gap-2">
            <button class="btn btn-primary btn-lg" @click="selectMode('private')">
              Private Message
            </button>
            <button class="btn btn-secondary btn-lg" @click="selectMode('group')">
              Create Group
            </button>
          </div>

          <div v-else-if="mode === 'private'" class="mb-3">
            <label class="form-label">Select user:</label>
            <div v-if="loading" class="text-center">Loading users...</div>
            <div v-else class="list-group" style="max-height: 300px; overflow-y: auto;">
              <button
                v-for="user in users"
                :key="user"
                type="button"
                class="list-group-item list-group-item-action"
                :class="{ active: selectedUser === user }"
                @click="selectedUser = user"
              >
                {{ user }}
              </button>
            </div>
            <div v-if="selectedUser" class="mt-3">
              <label for="initialMessage" class="form-label">Your message:</label>
              <textarea
                id="initialMessage"
                v-model="initialMessage"
                class="form-control"
                rows="3"
                placeholder="Say hello..."
                autofocus
              ></textarea>
            </div>
            <div class="mt-3">
              <button
                class="btn btn-primary"
                @click="createPrivateConversation"
                :disabled="!selectedUser || !initialMessage.trim() || loading"
              >
                {{ loading ? 'Creating...' : 'Create' }}
              </button>
              <button class="btn btn-secondary ms-2" @click="selectMode(null)">
                Back
              </button>
            </div>
          </div>

          <div v-else-if="mode === 'group'" class="mb-3">
            <label for="groupName" class="form-label">Group name:</label>
            <input
              id="groupName"
              v-model="groupName"
              type="text"
              class="form-control"
              placeholder="Enter group name"
              autofocus
            />
            <div class="mt-3">
              <button
                class="btn btn-secondary"
                @click="createGroup"
                :disabled="!groupName.trim() || loading"
              >
                {{ loading ? 'Creating...' : 'Create Group' }}
              </button>
              <button class="btn btn-secondary ms-2" @click="selectMode(null)">
                Back
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </dialog>
</template>

<style scoped>
dialog::backdrop {
  background-color: rgba(0, 0, 0, 0.5);
}

dialog {
  border: none;
  border-radius: 0.5rem;
  padding: 0;
}

dialog[open] {
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>

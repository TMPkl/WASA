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
              <h5 class="modal-title">Wybierz element</h5>
              <button type="button" class="btn-close" @click="showModal = false"></button>
            </div>
            <div class="modal-body">
              <div class="d-flex flex-column gap-2">
                <input 
                  type="text" 
                  class="form-control mb-3" 
                  placeholder="Search user..." 
                  v-model="searchQuery" 
                />

                <ul class="list-group overflow-auto" style="max-height: 400px; cursor: pointer;">
                  <li 
                    v-for="item in filteredItems"
                    :key="(item && (item.id || item.username || item.name)) || item"
                    class="list-group-item"
                    @click="selectItem(item)"
                  >
                    {{ (item && (item.name || item.username)) || item }}
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
  props: {
    items: {
      type: Array,
      required: false,
      default: () => []
    }
  },
  data() {
    return {
      showModal: false,
      searchQuery: '',
      itemsInternal: [],
      loading: false,
      error: null
    }
  },
  computed: {
    filteredItems() {
      const list = (this.itemsInternal && this.itemsInternal.length) ? this.itemsInternal : this.items;
      if (!this.searchQuery) return list;
      return list.filter(item =>
        (String((item && (item.name || item.username)) || item)).toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    }
  },
  methods: {
    async open(filterExistingContacts = true, groupMembers = []) {
      this.error = null;
      this.loading = true;
      try {
        const token = localStorage.getItem('token');
        const currentUsername = localStorage.getItem('username');

        const usersRes = await axios({
          method: 'get',
          url: `${__API_URL__}/users`,
          headers: token ? { Authorization: `Bearer ${token}` } : {}
        });
        let allUsers = Array.isArray(usersRes.data) ? usersRes.data : (usersRes.data.users || []);

        if (filterExistingContacts) {
          const convRes = await axios({
            method: 'get',
            url: `${__API_URL__}/conversations/${currentUsername}`,
            headers: {
              Authorization: token ? `Bearer ${token}` : undefined
            }
          });
          const conversations = Array.isArray(convRes.data) ? convRes.data : [];

          const existingContacts = new Set(
            conversations
              .filter(conv => conv.ConversationType === 'private')
              .map(conv => conv.OtherUsername)
          );

          this.itemsInternal = allUsers.filter(user => 
            user !== currentUsername && !existingContacts.has(user)
          );
        } else {
          // Filter out current user and existing group members
          const groupMembersSet = new Set(groupMembers || []);
          this.itemsInternal = allUsers.filter(user => 
            user !== currentUsername && !groupMembersSet.has(user)
          );
        }
      } catch (e) {
        this.itemsInternal = [];
        this.error = e?.response?.data?.error || e.message || 'Failed to load users';
        console.warn('UserList.open: could not fetch users', e);
      } finally {
        this.loading = false;
        this.showModal = true;
      }
    },
    async selectItem(item) {
      const username = (item && (item.username || item.name)) || item;
      this.$emit('select', username, item);
      this.showModal = false;
    }
  }
}
</script>

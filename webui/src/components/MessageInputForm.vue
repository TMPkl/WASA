<script>
export default {
  name: 'MessageInputForm',
  props: {
    isGroup: {
      type: Boolean,
      default: false
    },
    sending: {
      type: Boolean,
      default: false
    }
  },
  emits: ['send-message', 'add-user', 'leave-group'],
  data() {
    return {
      messageContent: '',
      selectedFile: null
    };
  },
  methods: {
    handleFileSelect(event) {
      const file = event.target.files[0];
      this.selectedFile = file || null;
    },
    openFileDialog() {
      this.$refs.fileInput.click();
    },
    removeFile() {
      this.selectedFile = null;
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = '';
      }
    },
    handleSubmit(event) {
      event.preventDefault();
      
      if (!this.messageContent.trim() && !this.selectedFile) {
        return;
      }

      this.$emit('send-message', {
        content: this.messageContent,
        files: this.selectedFile ? [this.selectedFile] : []
      });
    },
    clearForm() {
      this.messageContent = '';
      this.selectedFile = null;
      if (this.$refs.fileInput) {
        this.$refs.fileInput.value = '';
      }
    }
  }
};
</script>
<template>
  <div class="message-input-form">
    <!-- Selected file preview -->
    <div v-if="selectedFile" class="mb-2">
      <div class="badge bg-info d-flex align-items-center">
        <span>{{ selectedFile.name }}</span>
        <button 
          type="button" 
          class="btn-close btn-close-white ms-2" 
          @click="removeFile"
          style="font-size: 0.7rem;"
        ></button>
      </div>
    </div>

    <form @submit="handleSubmit" class="d-flex flex-column">
      <div class="d-flex mb-2">
        <input
          v-model="messageContent"
          type="text"
          class="form-control me-2"
          placeholder="Write a message..."
          :disabled="sending"
        />
        <input 
          ref="fileInput"
          type="file" 
          @change="handleFileSelect" 
          style="display: none;"
        />
        <button 
          type="button" 
          class="btn btn-secondary me-2"
          @click="openFileDialog"
          :disabled="sending"
          title="Attach files"
        >
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#paperclip"/></svg>
        </button>
        <button 
          type="submit" 
          class="btn btn-primary"
          :disabled="sending || (!messageContent.trim() && !selectedFile)"
        >
          {{ sending ? 'Sending...' : 'Send' }}
        </button>
      </div>
      
      <div v-if="isGroup" class="d-flex">
        <button type="button" class="btn btn-info btn-sm me-2" @click="$emit('add-user')">Add User</button>
        <button type="button" class="btn btn-danger btn-sm" @click="$emit('leave-group')">Leave Group</button>
      </div>
    </form>
  </div>
</template>



<style scoped>
.dalekooooo {
  gap: 0.5rem;
}
</style>

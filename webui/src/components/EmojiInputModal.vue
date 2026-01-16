<script>
export default {
  name: 'EmojiInputModal',
  data() {
    return {
      emoji: '',
      messageId: null
    };
  },
  methods: {
    open(messageId) {
      this.messageId = messageId;
      this.emoji = '';
      this.$refs.modal.showModal();
    },
    isSingleEmoji(str) {
      const trimmed = str.trim();
      const emojiRegex = /^(\p{Emoji_Presentation}|\p{Extended_Pictographic})[\p{Emoji_Modifier}]*$/u;
      return emojiRegex.test(trimmed);
    },
    addReaction() {
      if (!this.emoji || !this.emoji.trim()) {
        alert('Please enter an emoji');
        return;
      }
      if (!this.isSingleEmoji(this.emoji)) {
        alert('Please enter a single emoji only');
        return;
      }
      this.$emit('add-reaction', this.messageId, this.emoji.trim());
      this.$refs.modal.close();
    },
    cancel() {
      this.$refs.modal.close();
    }
  }
}
</script>

<template>
  <dialog ref="modal" class="modal" @keydown.enter="addReaction" @keydown.escape="cancel">
    <div class="modal-dialog">
      <div class="modal-content">
        <div class="modal-header">
          <h5 class="modal-title">Add reaction</h5>
          <button type="button" class="btn-close" @click="cancel" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          <div class="mb-3">
            <label for="emojiInput" class="form-label">Enter emoji:</label>
            <input 
              id="emojiInput"
              v-model="emoji" 
              type="text" 
              class="form-control form-control-lg text-center" 
              placeholder=""
              maxlength="10"
              autofocus
            />
          </div>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="cancel">Cancel</button>
          <button type="button" class="btn btn-primary" @click="addReaction">Add</button>
        </div>
      </div>
    </div>
  </dialog>
</template>

<style scoped>
dialog::backdrop {
  background-color: transparent;
}

dialog {
  border: none;
  border-radius: 0.5rem;
  padding: 0;
  width: auto;
  max-width: 200px;
  background-color: transparent;
}

dialog[open] {
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal-dialog {
  width: 100%;
  margin: 0;
}
</style>

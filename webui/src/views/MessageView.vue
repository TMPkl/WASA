<script>
import MessageItem from '../components/MessageItem.vue';
import ConversationListModal from '../components/ConversationListModal.vue';
import UserList from '../components/UserList.vue';
import EmojiInputModal from '../components/EmojiInputModal.vue';
import axios from 'axios';

export default {
    name: 'MessageView',
    components: {
        MessageItem,
        ConversationListModal,
        UserList,
        EmojiInputModal
    },
    props: {
        messages: {
            type: Array,
            default: () => []
        }
    },
    data() {
        return {
            pendingForwardMessageId: null,
            forwardError: null
        };
    },
    methods: {
        forwardMessage(messageId) {
            this.pendingForwardMessageId = messageId;
            this.$nextTick(() => {
                if (this.$refs.convListModal && this.$refs.convListModal.open) {
                    this.$refs.convListModal.open();
                }
            });
        },
        openEmojiModal(messageId) {
            this.$refs.emojiModal.open(messageId);
        },
        openNewContactList() {
            if (this.$refs.userListModal && this.$refs.userListModal.open) {
                this.$refs.userListModal.open();
            }
        },
        async handleNewContactSelect(username) {
            if (!this.pendingForwardMessageId) return;

            const token = localStorage.getItem('token');
            const currentUsername = localStorage.getItem('username');

            try {
                await axios({
                    method: 'post',
                    url: `${__API_URL__}/messages/${this.pendingForwardMessageId}/forwards`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: currentUsername,
                        newContactUsername: username
                    }
                });
                this.forwardError = null;
            } catch (e) {
                console.error('Failed to forward message:', e);
                this.forwardError = e?.response?.data?.error || e.message;
            } finally {
                this.pendingForwardMessageId = null;
            }
        },
        async handleConversationSelect(conversationId, conv) {
            if (!this.pendingForwardMessageId) return;

            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');

            try {
                await axios({
                    method: 'post',
                    url: `${__API_URL__}/messages/${this.pendingForwardMessageId}/forwards`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: username,
                        addressingConversationID: conversationId
                    }
                });
                this.forwardError = null;
            } catch (e) {
                console.error('Failed to forward message:', e);
                this.forwardError = e?.response?.data?.error || e.message;
            } finally {
                this.pendingForwardMessageId = null;
            }
        },
        async handleDeleteMessage(messageId) {
            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');
            try {
                await axios({
                    method: 'delete',
                    url: `${__API_URL__}/messages/${messageId}`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: username
                    }
                });
                this.$emit('message-deleted', messageId);

            } catch (e) {
                console.error('Failed to delete message:', e);
                alert('Failed to delete message: ' + (e?.response?.data?.error || e.message));
            }
        },
        async removeReaction(messageId) {
            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');
            try {
                await axios({
                    method: 'delete',
                    url: `${__API_URL__}/messages/${messageId}/reactions`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: username
                    }
                });
                this.$emit('reaction-removed', messageId);
            } catch (e) {
                console.error('Failed to remove reaction:', e);
            }
        },
        isSingleEmoji(str) {
            const trimmed = str.trim();
            const emojiRegex = /^(\p{Emoji_Presentation}|\p{Extended_Pictographic})[\p{Emoji_Modifier}]*$/u;
            return emojiRegex.test(trimmed);
        },
        async addReaction(messageId, emoji) {
            if (!this.isSingleEmoji(emoji)) {
                console.error('Invalid reaction: must be a single emoji');
                alert('Reaction must be a single emoji');
                return;
            }
            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');
            try {
                await axios({
                    method: 'post',
                    url: `${__API_URL__}/messages/${messageId}/reactions`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: username,
                        emoji: emoji
                    }
                });
                this.$emit('reaction-added', messageId);
            } catch (e) {
                console.error('Failed to add reaction:', e);
            }
        },
        async handleDownloadAttachment(messageId) {
            try {
                const token = localStorage.getItem('token');
                const username = localStorage.getItem('username');
                
                const response = await axios({
                    method: 'post',
                    url: `${__API_URL__}/attachments/${messageId}`,
                    headers: {
                        Authorization: `Bearer ${token}`
                    },
                    responseType: 'blob'
                });

                const contentDisposition = response.headers['content-disposition'];
                let filename = `attachment_${messageId}`;
                
                if (contentDisposition) {
                    const filenameMatch = contentDisposition.match(/filename="?(.+?)"?$/i);
                    if (filenameMatch) {
                        filename = filenameMatch[1];
                    }
                }

                
                const contentType = response.headers['content-type'];
                if (contentType && !filename.includes('.')) {
                    const ext = this.getExtensionFromMimeType(contentType);
                    if (ext) {
                        filename += ext;
                    }
                }

                const url = window.URL.createObjectURL(new Blob([response.data]));
                const link = document.createElement('a');
                link.href = url;
                link.setAttribute('download', filename);
                document.body.appendChild(link);
                link.click();
                link.remove();
                window.URL.revokeObjectURL(url);
            } catch (e) {
                console.error('Failed to download attachment:', e);
                alert('Failed to download attachment: ' + (e?.response?.data?.error || e.message));
            }
        },
        getExtensionFromMimeType(mimeType) {
            const mimeToExt = {
                'image/jpeg': '.jpg',
                'image/png': '.png',
                'image/gif': '.gif',
                'image/webp': '.webp',
                'application/pdf': '.pdf',
                'text/plain': '.txt',
                'application/zip': '.zip',
                'application/json': '.json',
                'video/mp4': '.mp4',
                'audio/mpeg': '.mp3'
            };
            return mimeToExt[mimeType] || '';
        }
    }
}
</script>
<template>
    <div class="message-view">
        <div v-if="messages.length === 0" class="text-center text-muted p-4">
            No messages here yet
        </div>
        <MessageItem
            v-else
            v-for="message in messages"
            :key="message.id"
            :message-id="message.id"
            :sender="message.sender"
            :content="message.content"
            :timestamp="message.timestamp"
            :attachment="message.attachment"
            :reactions="message.reactions"
            :replying-to-id="message.replyingToId"
            :replying-to-sender="message.replyingToSender"
            :replying-to-content="message.replyingToContent"
            @download-attachment="handleDownloadAttachment"
            @delete-message="handleDeleteMessage"
            @forward-message="forwardMessage"
            @reply-message="$emit('reply-message', $event)"
            @open-emoji-modal="openEmojiModal"
        />

        <ConversationListModal
            ref="convListModal"
            title="Forward message to"
            @select="handleConversationSelect"
            @open-new-contact-list="openNewContactList"
        />

        <UserList
            ref="userListModal"
            @select="handleNewContactSelect"
        />

        <EmojiInputModal
            ref="emojiModal"
            @add-reaction="addReaction"
        />
    </div>
</template>
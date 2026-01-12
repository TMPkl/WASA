<script>
import MessageItem from '../components/MessageItem.vue';
import axios from 'axios';

export default {
    name: 'MessageView',
    components: {
        MessageItem
    },
    props: {
        messages: {
            type: Array,
            default: () => []
        }
    },
    methods: {
        async forwardMessage(messageId) {
            const conversationId = prompt('Enter the conversation ID to forward this message to:');
            
            if (!conversationId) {
                return;
            }
            
            const token = localStorage.getItem('token');
            const username = localStorage.getItem('username');
            
            try {
                await axios({
                    method: 'post',
                    url: `${__API_URL__}/messages/${messageId}/forwards`,
                    headers: {
                        Authorization: `Bearer ${token}`,
                        'Content-Type': 'application/json'
                    },
                    data: {
                        username: username,
                        addressingConversationID: parseInt(conversationId)
                    }
                });
                alert('Message forwarded successfully!');
            } catch (e) {
                console.error('Failed to forward message:', e);
                alert('Failed to forward message: ' + (e?.response?.data?.error || e.message));
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
            @download-attachment="handleDownloadAttachment"
            @delete-message="handleDeleteMessage"
            @forward-message="forwardMessage"
        />
    </div>
</template>
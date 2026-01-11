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

                // Get filename from Content-Disposition header or use default
                const contentDisposition = response.headers['content-disposition'];
                let filename = `attachment_${messageId}`;
                
                if (contentDisposition) {
                    const filenameMatch = contentDisposition.match(/filename="?(.+?)"?$/i);
                    if (filenameMatch) {
                        filename = filenameMatch[1];
                    }
                }

                // Determine file extension from content type
                const contentType = response.headers['content-type'];
                if (contentType && !filename.includes('.')) {
                    const ext = this.getExtensionFromMimeType(contentType);
                    if (ext) {
                        filename += ext;
                    }
                }

                // Create blob link to download
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
        <MessageItem
            v-for="message in messages"
            :key="message.id"
            :message-id="message.id"
            :sender="message.sender"
            :content="message.content"
            :timestamp="message.timestamp"
            :attachment="message.attachment"
            @download-attachment="handleDownloadAttachment"
        />
    </div>
</template>
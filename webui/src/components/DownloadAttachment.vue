<template>
  <div class="download-attachment">
    <button 
      class="btn btn-sm btn-success" 
      @click="downloadAttachment"
      :disabled="downloading"
    >
      {{ downloading ? 'Downloading...' : 'Download Attachments' }}
    </button>
  </div>
</template>

<script>
import axios from 'axios';

export default {
  name: 'DownloadAttachment',
  props: {
    messageId: {
      type: [String, Number],
      required: true
    }
  },
  data() {
    return {
      downloading: false,
      error: null
    };
  },
  methods: {
    async downloadAttachment() {
      this.downloading = true;
      this.error = null;
      
      try {
        const token = localStorage.getItem('token');
        const username = localStorage.getItem('username');
        
        const response = await axios({
          method: 'post',
          url: `${__API_URL__}/attachments/${this.messageId}`,
          headers: {
            Authorization: `Bearer ${token}`
          },
          responseType: 'blob'
        });

        // Get filename from Content-Disposition header or use default
        const contentDisposition = response.headers['content-disposition'];
        let filename = `attachment_${this.messageId}`;
        
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
        
        this.$emit('download-success');
      } catch (e) {
        this.error = e?.response?.data?.error || e.message || 'Failed to download attachment';
        this.$emit('download-error', this.error);
        console.error('Download error:', this.error);
      } finally {
        this.downloading = false;
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
};
</script>

<style scoped>
.download-attachment {
  display: inline-block;
}
</style>
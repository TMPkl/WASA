<script>
export default {
  name: 'MessageItem',
  props: {
    messageId: [String, Number],
    sender: String,
    content: String,
    timestamp: String,
    attachment: Boolean
  },
  methods: {
    downloadAttachments() {
      this.$emit('download-attachment', this.messageId);
    }
  }
}
</script>

    <template>
        <div class="d-flex flex-column p-2 border-bottom m-2 bg-light rounded-3 overflow-x-hidden">
            <div class="d-flex justify-content-between"> 
            <div class="p-2  small ">{{ sender }}</div>
            
            <div v-if="attachment" class="p-2">
                <button class="btn btn-sm btn-success" @click="downloadAttachments">Download Attachments</button> 
            </div>
        <div class="p-2 small ">{{ timestamp}}</div>   
        </div>
            <div class=" p-2 Large">{{ content }}</div>        
        </div>
        
    </template>
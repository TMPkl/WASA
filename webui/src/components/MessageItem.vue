<script>
export default {
  name: 'MessageItem',
  props: {
    messageId: [String, Number],
    sender: String,
    content: String,
    timestamp: String,
    attachment: Boolean,
    reactions: Array
  },
  methods: {
    downloadAttachments() {
      this.$emit('download-attachment', this.messageId);
    },
    getUsername() {
      return localStorage.getItem('username');
    },
    deleteMessage(){
        this.$emit('delete-message', this.messageId);
    },
    forwardMessage(){
        this.$emit('forward-message', this.messageId);
    },
    addReaction() {
      this.$emit('open-emoji-modal', this.messageId);
    }
  }
}
</script>

    <template>
        <div class="d-flex flex-column p-2 border-bottom m-2 bg-light rounded-3 overflow-x-hidden">
            <div class="d-flex justify-content-between"> 
            <div class="p-2  small ">{{ sender }}</div>

            <div class="d-flex justify-content-center"> 
                <div v-if="attachment" class="p-2">
                    <button class="btn btn-sm btn-success" @click="downloadAttachments">Download Attachments</button> 
                </div>

                <div v-if="getUsername()==sender" class="p-2">
                    <span class="">
                        <button class="btn btn-sm btn-danger" @click="deleteMessage">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#trash-2"/>
                            </svg>
                </button></span></div>
                <div v-if="getUsername()!=sender" class="p-2" >
                    <span class="">
                        <button class="btn btn-sm btn-danger" @click="deleteMessage" disabled>
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#trash-2"/>
                            </svg>
                </button></span></div>
                <div  class="p-2">
                    <span class="">
                        <button class="btn btn-sm btn-secondary" @click="forwardMessage" >
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#repeat"/>
                            </svg>
                        </button>
                    </span>
                </div>
                <div class="p-2">
                    <span class="">
                        <button class="btn btn-sm btn-primary" @click="addReaction" title="Add reaction">
                            <svg class="feather">
                                <use href="/feather-sprite-v4.29.0.svg#smile"/>
                            </svg>
                        </button>
                    </span>
                </div>
            </div>
        <div class="p-2 small ">{{ timestamp}}</div>   
        </div>
            <div class=" p-2 Large">{{ content }}</div>
            <div v-if="reactions && reactions.length > 0" class="p-2 small">
                <div class="badge bg-light text-dark" v-for="reaction in reactions" :key="reaction">
                    {{ reaction }}
                </div>
            </div>
        </div>
        
    </template>
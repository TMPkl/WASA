<script>
export default {
  name: 'MessageItem',
  props: {
    sender: String,
    content: String,
    timestamp: String,
    attachment: Boolean //jesli bedzie attachment to dodaj przycisk do pobrania go
  }
}
</script>

    <template>
        <div class="d-flex flex-column p-2 border-bottom m-2 bg-light rounded-3 overflow-x-hidden">
            <div class="d-flex justify-content-between"> 
            <div class="p-2  small ">{{ sender }}</div>
            
            <div v-if="attachment" class="p-2">
                <button class="btn btn-sm btn-success">Pobierz załącznik</button> 
            </div>
        <div class="p-2 small ">{{ timestamp}}</div>   
        </div>
            <div class=" p-2 Large">{{ content }}</div>        
        </div>
        
    </template>
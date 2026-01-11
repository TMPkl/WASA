<script>
    import MessageItem from '../components/MessageItem.vue';

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
        }
    }
</script>
<template>
    <div class="message-view">
        <MessageItem
            v-for="message in messages"
            :key="message.id"
            :sender="message.sender"
            :content="message.content"
            :timestamp="message.timestamp"
            :attachment="message.attachment"
        />
    </div>
</template>
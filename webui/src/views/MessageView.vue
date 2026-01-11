<script>
    import MessageItem from '../components/MessageItem.vue';

    export default {
        name: 'MessageView',
        components: {
            MessageItem
        },
        data() {
            return {
                hardcodedMessages: [
                    {
                        id: 1,
                        sender: 'Jan ktos ',
                        content: 'Cześć! Jak się masz?AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAASSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSSAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA AAAAaaAAAAAAAAAAaa',
                        timestamp: '2026-01-11 12:34',
                        attachment: null
                    },
                    {
                        id: 2,
                        sender: 'Anna Nowak',
                        content: 'plik',
                        timestamp: '026-01-11 12:34',
                        attachment: 'plik.pdf'
                    },
                    {
                        id: 3,
                        sender: 'Jan ktos ',
                        content: 'bla bla bla ',
                        timestamp: '026-01-11 12:34',
                        attachment: null
                    }
                ]
            }
        }
    }
</script>
<template>
    <div class="message-view">
        <MessageItem
            v-for="message in hardcodedMessages"
            :key="message.id"
            :sender="message.sender"
            :content="message.content"
            :timestamp="message.timestamp"
            :attachment="message.attachment"
        />
    </div>
</template>
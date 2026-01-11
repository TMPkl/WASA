<script>
export default {
  name: 'ContactItem',
  props: {
    title: String,
    lastMessage: String,
    time: String,
    status: String
  }
}
</script>

<template>
  <div class="list-group-item">
    <div class="fw-bold">{{ title }}</div>
    <small class="text-muted">{{ lastMessage }}</small>
      <div class = "list-group-item-extra text-end">
        <div class="text-muted">{{ time }}</div>
        <div class="text-muted">{{ status }}</div>
      </div>
  </div>
  <div class="d-flex p-3 bg-primary text-white">
  <div class="p-2 bg-info">Flex item 1</div>
    <div class="d-flex flex-column">
      <div class="p-2 bg-info">Flex item 1</div>
      <div class="p-2 bg-warning">Flex item 2</div>
    </div>
</div>
</template>

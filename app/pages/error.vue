<script setup>
const route = useRoute();

const status = ref(route.query.status);
const message = ref(route.query.error);
</script>

<template>
  <Frame :page-title="status.toUpperCase()" :page-message="message"></Frame>
</template>

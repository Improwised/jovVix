<script setup>
const route = useRoute();

const status = ref(route.query.status);
const message = ref(route.query.error);
</script>

<template>
  <FrameLayout
    :page-title="status.toUpperCase()"
    :page-welcome-message="message"
  >
  </FrameLayout>
</template>

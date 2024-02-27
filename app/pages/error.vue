<script setup>
import Frame from "../components/Frame.vue";
const route = useRoute();

const status = ref(route.query.status);
const message = ref(route.query.error);
console.log(route.query, route.query.error);
</script>

<template>
  <Frame :page-title="status.toUpperCase()" :page-message="message"></Frame>
</template>

<script setup>
definePageMeta({
  layout: "empty",
});

const route = useRoute();

const status = ref(
  route.query.status || "Something went wrong. Please try again later."
);
const message = ref(route.query.error);
</script>

<template>
  <div class="bg-jv-canvas">
    <Frame :page-title="status?.toUpperCase()" :page-message="message" />
  </div>
</template>

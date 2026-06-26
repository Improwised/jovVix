<script setup>
definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "Page Not Found - jovVix",
  description:
    "The page you are looking for could not be found. Return to jovVix to host or join a quiz.",
  robots: "noindex, nofollow",
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

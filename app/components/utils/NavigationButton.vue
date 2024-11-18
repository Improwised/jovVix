<template>
  <NuxtLink class="nav-link" :to="props.navigateTo">
    <button type="button" class="btn text-white btn-primary me-2">
      {{ props.title }}
    </button>
  </NuxtLink>
</template>

<script setup>
const props = defineProps({
  title: {
    type: String,
    required: true,
    default: "",
  },
  navigateTo: {
    type: String,
    required: true,
    default: "",
  },
});
</script>

<script setup>
// Define props using defineProps
const props = defineProps({
  pageTitle: {
    type: String,
    required: true,
    default: "title",
  },
  pageWelcomeMessage: {
    type: String,
    required: false,
    default: "message",
  },
});
</script>

<template>
  <div class="d-flex justify-content-center">
    <div class="border p-2 m-0 m-sm-5 p-sm-5 max-width rounded">
      <h1>{{ pageTitle }}</h1>
      <h6 v-if="props.pageWelcomeMessage">{{ pageWelcomeMessage }}</h6>

      <hr class="m-2" />
      <slot></slot>
    </div>
  </div>
</template>

<style>
.max-width {
  width: 700px;
}
</style>

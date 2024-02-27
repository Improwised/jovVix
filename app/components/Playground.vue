<script setup>
const bodyRef = ref();

const props = defineProps({
  fullScreenEnabled: {
    default: false,
    type: Boolean,
  },
});
const emits = defineEmits(["isFullScreen"]);

onMounted(() => {
  if (process.client) {
    bodyRef.value = document.documentElement;

    addEventListener("fullscreenchange", () => {
      console.log(document.fullscreenElement);
      emits("isFullScreen", document.fullscreenElement != null ? true : false);
    });
  }
});

watch(
  () => props.fullScreenEnabled,
  () => {
    if (!bodyRef.value) {
      return;
    }

    if (document.fullscreenElement) {
      document.exitFullscreen();
      emits("isFullScreen", false);
    } else {
      bodyRef.value.requestFullscreen();
      emits("isFullScreen", true);
    }
  }
);
</script>

<template>
  <div class="full-screen">
    <slot></slot>
  </div>
</template>

<style scoped>
.full-screen {
  max-height: 100vh;
  width: 100%;
}
</style>

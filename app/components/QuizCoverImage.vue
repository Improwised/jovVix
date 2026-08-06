<template>
  <NuxtImg
    :src="currentSrc"
    :alt="alt"
    :width="width"
    :height="height"
    fit="cover"
    :loading="eager ? 'eager' : 'lazy'"
    :preload="eager"
    :fetchpriority="eager ? 'high' : 'auto'"
    decoding="async"
    @error="failed = true"
  />
</template>

<script setup>
import { computed, ref, watch } from "vue";

const props = defineProps({
  src: {
    type: String,
    default: "",
  },
  fallback: {
    type: String,
    default: "",
  },
  alt: {
    type: String,
    default: "",
  },
  width: {
    type: [String, Number],
    default: 320,
  },
  height: {
    type: [String, Number],
    default: 124,
  },
  eager: {
    type: Boolean,
    default: false,
  },
});

const failed = ref(false);

const currentSrc = computed(() =>
  failed.value && props.fallback ? props.fallback : props.src
);

watch(
  () => props.src,
  () => {
    failed.value = false;
  }
);
</script>

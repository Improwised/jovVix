<template>
  <NuxtLink
    :to="props.navigateTo"
    class="jv-card inline-flex h-11 items-center justify-center gap-2 border-2 border-jv-ink bg-jv-coral px-6 font-headings text-base text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
  >
    {{ props.title }}
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

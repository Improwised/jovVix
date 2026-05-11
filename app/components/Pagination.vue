<template>
  <div class="flex items-center justify-center gap-2 mt-8">
    <NuxtLink
      v-if="props.page > 1"
      class="flex h-10 w-10 rotate-[-1deg] items-center justify-center jv-border-rough bg-jv-white text-jv-ink shadow-brutal-sm hover:bg-jv-yellow/50 transition-colors"
      :to="{
        path: route.path,
        query: { ...route.query, page: props.page - 1 },
      }"
    >
      <span aria-hidden="true" class="text-xl font-black">&lt;</span>
    </NuxtLink>
    <div
      v-else
      class="flex h-10 w-10 rotate-[-1deg] items-center justify-center jv-border-rough bg-jv-canvas text-jv-ink/30 shadow-brutal-sm cursor-not-allowed"
    >
      <span aria-hidden="true" class="text-xl font-black">&lt;</span>
    </div>

    <div
      class="flex h-10 w-10 items-center justify-center jv-border-rough bg-jv-coral text-white text-[16px] font-black shadow-brutal-sm rotate-[1deg]"
    >
      {{ props.page }}
    </div>

    <NuxtLink
      v-if="props.page < props.numOfRecords"
      class="flex h-10 w-10 rotate-[1deg] items-center justify-center jv-border-rough bg-jv-white text-jv-ink shadow-brutal-sm hover:bg-jv-yellow/50 transition-colors"
      :to="{
        path: route.path,
        query: { ...route.query, page: props.page + 1 },
      }"
    >
      <span aria-hidden="true" class="text-xl font-black">&gt;</span>
    </NuxtLink>
    <div
      v-else
      class="flex h-10 w-10 rotate-[1deg] items-center justify-center jv-border-rough bg-jv-canvas text-jv-ink/30 shadow-brutal-sm cursor-not-allowed"
    >
      <span aria-hidden="true" class="text-xl font-black">&gt;</span>
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const props = defineProps({
  page: {
    type: Number,
    required: true,
    default: 0,
  },
  numOfRecords: {
    type: Number,
    required: true,
    default: 0,
  },
});
</script>

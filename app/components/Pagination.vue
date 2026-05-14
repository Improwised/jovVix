<template>
  <div class="flex items-center justify-center gap-2 mt-8">
    <NavigationLink
      v-if="props.page > 1"
      :url="{
        path: route.path,
        query: { ...route.query, page: props.page - 1 },
      }"
      class="h-10 md:p-2 bg-jv-white text-jv-ink"
    >
      <span aria-hidden="true" class="text-xl font-black">&lt;</span>
    </NavigationLink>
    <NavigationLink
      v-else
      disabled
      class="h-10 md:p-2 bg-jv-white text-jv-ink cursor-not-allowed"
    >
      <span aria-hidden="true" class="text-xl font-black">&lt;</span>
    </NavigationLink>

    <div
      class="flex h-10 w-10 items-center justify-center jv-border-rough bg-jv-coral text-white text-[16px] font-black shadow-brutal-sm rotate-[1deg]"
    >
      {{ props.page }}
    </div>

    <NavigationLink
      v-if="props.page < props.numOfRecords"
      :url="{
        path: route.path,
        query: { ...route.query, page: props.page + 1 },
      }"
      class="h-10 md:p-2 bg-jv-white text-jv-ink"
    >
      <span aria-hidden="true" class="text-xl font-black">&gt;</span>
    </NavigationLink>
    <NavigationLink
      v-else
      disabled
      class="h-10 md:p-2 bg-jv-white text-jv-ink cursor-not-allowed"
    >
      <span aria-hidden="true" class="text-xl font-black">&gt;</span>
    </NavigationLink>
  </div>
</template>

<script setup>
import NavigationLink from "./common/NavigationLink.vue";
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

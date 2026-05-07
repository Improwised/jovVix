<template>
  <NuxtLink
    :to="url"
    :class="
      cn(
        'inline-flex items-center justify-center font-headings text-center whitespace-nowrap jv-border-uneven shadow-brutal transition-transform transform transform-gpu backface-hidden active:shadow-none bg-jv-yellow gap-1.5 px-5 sm:px-6 md:px-7 py-2.5 md:py-3 text-sm md:text-lg hover:rotate-[2deg]',
        props.class
      )
    "
  >
    <slot></slot>
    <span>{{ urlName }}</span>
  </NuxtLink>
</template>

<script setup>
import { cn } from "@/lib/utils";

const props = defineProps({
  url: {
    default: "",
    type: String,
    required: true,
  },
  urlName: {
    default: "",
    type: String,
    required: true,
  },
  class: {
    type: [String, Array, Object],
    default: "",
  },
});
</script>

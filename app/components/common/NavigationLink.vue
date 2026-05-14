<template>
  <NuxtLink
    v-if="url"
    :to="url"
    v-bind="$attrs"
    :class="cn(baseClasses, props.class)"
  >
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </NuxtLink>
  <Button v-else v-bind="$attrs" :class="cn(baseClasses, props.class)">
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </Button>
</template>

<script setup>
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

defineOptions({ inheritAttrs: false });

const props = defineProps({
  url: {
    type: String,
    default: "",
  },
  urlName: {
    type: String,
    default: "",
  },
  class: {
    type: [String, Array, Object],
    default: "",
  },
});

const baseClasses =
  "inline-flex items-center justify-center font-headings text-center whitespace-nowrap jv-border-uneven shadow-brutal transition-transform transform transform-gpu backface-hidden active:shadow-none bg-jv-yellow gap-1.5 px-5 sm:px-6 md:px-7 py-2.5 md:py-3 text-sm md:text-lg hover:rotate-[2deg] disabled:opacity-60 disabled:cursor-not-allowed disabled:hover:rotate-0 disabled:active:shadow-brutal aria-disabled:opacity-60 aria-disabled:cursor-not-allowed aria-disabled:pointer-events-none";
</script>

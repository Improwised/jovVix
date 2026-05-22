<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";

interface Props {
  for?: string;
  required?: boolean;
  class?: HTMLAttributes["class"];
}

const props = defineProps<Props>();
</script>

<template>
  <label
    :for="props.for"
    data-slot="label"
    :class="
      cn(
        'mb-1.5 inline-flex items-center gap-1 font-body text-sm font-bold text-jv-ink',
        props.class
      )
    "
  >
    <slot />
    <span v-if="props.required" class="text-jv-coral" aria-hidden="true"
      >*</span
    >
  </label>
</template>

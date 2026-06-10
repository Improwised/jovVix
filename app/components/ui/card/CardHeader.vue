<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";

interface Props {
  class?: HTMLAttributes["class"];
}

const props = defineProps<Props>();
</script>

<template>
  <header
    data-slot="card-header"
    :class="cn('mb-4 flex items-start justify-between gap-3', props.class)"
  >
    <slot />
  </header>
</template>

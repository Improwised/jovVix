<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";

interface Props {
  class?: HTMLAttributes["class"];
}

const props = defineProps<Props>();
</script>

<template>
  <div data-slot="card-body" :class="cn('font-body text-jv-ink', props.class)">
    <slot />
  </div>
</template>

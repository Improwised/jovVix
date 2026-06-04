<script setup lang="ts">
import type { TooltipProviderProps } from "reka-ui";
import { TooltipProvider } from "reka-ui";

const props = withDefaults(defineProps<TooltipProviderProps>(), {
  delayDuration: 200,
});
</script>

<template>
  <TooltipProvider data-slot="tooltip-provider" v-bind="props">
    <slot />
  </TooltipProvider>
</template>

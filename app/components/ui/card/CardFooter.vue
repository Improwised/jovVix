<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";

interface Props {
  class?: HTMLAttributes["class"];
}

const props = defineProps<Props>();
</script>

<template>
  <footer
    data-slot="card-footer"
    :class="
      cn(
        'mt-4 flex flex-col-reverse items-stretch justify-end gap-3 sm:flex-row sm:items-center',
        props.class
      )
    "
  >
    <slot />
  </footer>
</template>

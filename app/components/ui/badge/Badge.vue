<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { badgeVariants, type BadgeVariants } from ".";

interface Props {
  variant?: BadgeVariants["variant"];
  size?: BadgeVariants["size"];
  shadow?: BadgeVariants["shadow"];
  as?: string;
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<Props>(), {
  variant: "default",
  size: "sm",
  shadow: "none",
  as: "span",
  class: undefined,
});
</script>

<template>
  <component
    :is="as"
    data-slot="badge"
    :class="
      cn(
        badgeVariants({
          variant: props.variant,
          size: props.size,
          shadow: props.shadow,
        }),
        props.class
      )
    "
  >
    <slot />
  </component>
</template>

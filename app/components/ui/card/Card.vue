<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { cardVariants, type CardVariants } from ".";

interface Props {
  variant?: CardVariants["variant"];
  tone?: CardVariants["tone"];
  shadow?: CardVariants["shadow"];
  padding?: CardVariants["padding"];
  tilt?: CardVariants["tilt"];
  as?: string;
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<Props>(), {
  variant: "default",
  tone: "white",
  shadow: "md",
  padding: "md",
  tilt: "none",
  as: "div",
  class: undefined,
});
</script>

<template>
  <component
    :is="as"
    data-slot="card"
    :class="
      cn(
        cardVariants({
          variant: props.variant,
          tone: props.tone,
          shadow: props.shadow,
          padding: props.padding,
          tilt: props.tilt,
        }),
        props.class
      )
    "
  >
    <slot />
  </component>
</template>

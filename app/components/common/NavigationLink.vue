<template>
  <NuxtLink
    v-if="url"
    :to="url"
    v-bind="$attrs"
    :class="cn(sharedClasses, variantClasses, props.class)"
  >
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </NuxtLink>
  <button
    v-else
    v-bind="$attrs"
    :class="cn(sharedClasses, variantClasses, props.class)"
  >
    <slot />
    <span v-if="urlName">{{ urlName }}</span>
  </button>
</template>

<script setup>
import { computed } from "vue";
import { cn } from "@/lib/utils";

defineOptions({ inheritAttrs: false });

const props = defineProps({
  url: {
    type: [String, Object],
    default: "",
  },
  urlName: {
    type: String,
    default: "",
  },
  variant: {
    type: String,
    default: "brutal",
  },
  class: {
    type: [String, Array, Object],
    default: "",
  },
});

const sharedClasses =
  "inline-flex items-center justify-center whitespace-nowrap disabled:opacity-60 disabled:cursor-not-allowed aria-disabled:opacity-60 aria-disabled:cursor-not-allowed aria-disabled:pointer-events-none";

const brutalClasses =
  "font-headings text-center jv-border-uneven shadow-brutal transition-transform transform transform-gpu backface-hidden active:shadow-none bg-jv-yellow gap-1.5 px-5 sm:px-6 md:px-7 py-2.5 md:py-3 text-sm md:text-lg hover:rotate-[2deg] disabled:hover:rotate-0 disabled:active:shadow-brutal";

const toggleClasses =
  "gap-1.5 rounded-md border px-2.5 h-8 text-[13px] font-semibold transition-colors";

const resultClasses =
  "h-11 gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 font-body text-sm font-black text-jv-ink shadow-brutal-sm cursor-pointer transition-transform hover:-rotate-1 active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:hover:rotate-0 focus-visible:outline-2 focus-visible:outline-offset-4 focus-visible:outline-jv-ink [&_svg]:size-4 [&_svg]:shrink-0";

const variantClasses = computed(
  () =>
    ({ toggle: toggleClasses, result: resultClasses })[props.variant] ||
    brutalClasses
);
</script>

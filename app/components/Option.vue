<script setup>
import { Check, X, User } from "lucide-vue-next";

const props = defineProps({
  order: {
    type: Number,
    required: false,
    default: 0,
  },
  selected: {
    type: Number,
    required: false,
    default: 0,
  },
  option: {
    type: String,
    required: false,
    default: "",
  },
  isCorrect: {
    type: Boolean,
    required: false,
    default: false,
  },
  isPicked: {
    type: Boolean,
    required: false,
    default: false,
  },
  optionsMedia: {
    type: String,
    required: true,
    default: "",
  },
  isAdminAnalysis: {
    type: Boolean,
    required: false,
    default: false,
  },
});

// Server returns presigned URLs for image options, but if generation fails it
// falls back to the raw object key. Treating that as <img src> resolves to the
// current path and 404s. Render an actual image only when the value is a
// fetchable URL.
const isFetchableUrl = (val) =>
  typeof val === "string" && /^(https?:|data:|blob:)/i.test(val.trim());
const renderAsImage = computed(
  () => props.optionsMedia === "image" && isFetchableUrl(props.option)
);
</script>

<template>
  <div class="flex w-full items-center gap-3 sm:gap-4">
    <span
      v-if="props.isCorrect"
      class="grid size-9 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white text-jv-ink sm:size-11"
      aria-label="Correct"
    >
      <Check class="size-4 sm:size-5" :stroke-width="3" />
    </span>
    <span
      v-else-if="props.isPicked"
      class="grid size-9 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white text-jv-coral sm:size-11"
      aria-label="Your incorrect answer"
    >
      <X class="size-4 sm:size-5" :stroke-width="3" />
    </span>
    <span
      v-else
      class="grid size-9 shrink-0 place-items-center rounded-[6px] border-[2px] border-jv-ink bg-jv-white font-feature text-[14px] font-black text-jv-ink sm:size-11 sm:text-[16px]"
    >
      {{ String.fromCharCode(64 + Number(props.order)) }}
    </span>

    <div class="min-w-0 flex-1">
      <div v-if="renderAsImage" class="flex items-center justify-center">
        <img
          :src="props.option"
          :alt="`Option ${props.order}`"
          class="max-h-[120px] w-auto object-contain"
        />
      </div>
      <CodeBlockComponent
        v-else-if="props.optionsMedia === 'code'"
        :code="props.option"
      />
      <span
        v-else
        class="block break-words font-body text-[15px] font-black leading-snug text-jv-ink sm:text-[17px]"
      >
        {{ props.option }}
      </span>
    </div>

    <span
      v-if="props.isAdminAnalysis"
      :class="[
        'inline-flex shrink-0 items-center gap-1 rounded-full border-[2px] border-jv-ink px-2.5 py-0.5 font-feature text-[12px] font-black text-jv-ink shadow-brutal-sm sm:text-[13px]',
        props.isCorrect ? 'bg-jv-mint' : 'bg-jv-white',
      ]"
    >
      <User class="size-3" :stroke-width="2.6" />
      {{ props.selected }}
    </span>
  </div>
</template>

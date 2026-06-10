<script setup lang="ts">
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { inputVariants, type InputVariants } from ".";

interface Props {
  modelValue?: string | null;
  size?: InputVariants["size"];
  invalid?: boolean;
  rows?: number;
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  size: "md",
  invalid: false,
  rows: 4,
  class: undefined,
});

defineEmits<{ "update:modelValue": [value: string] }>();
</script>

<template>
  <textarea
    :value="props.modelValue ?? ''"
    :rows="props.rows"
    data-slot="textarea"
    :class="
      cn(
        inputVariants({ size: props.size, invalid: props.invalid }),
        'min-h-[5.5rem] py-3 leading-relaxed',
        props.class
      )
    "
    @input="
      (e) =>
        $emit('update:modelValue', (e.target as HTMLTextAreaElement).value)
    "
  />
</template>

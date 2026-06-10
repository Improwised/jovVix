<script setup lang="ts">
import type { HTMLAttributes, InputHTMLAttributes } from "vue";
import { cn } from "@/lib/utils";
import { inputVariants, type InputVariants } from ".";

interface Props {
  modelValue?: string | number | null;
  size?: InputVariants["size"];
  invalid?: boolean;
  type?: InputHTMLAttributes["type"];
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: "",
  size: "md",
  invalid: false,
  type: "text",
  class: undefined,
});

defineEmits<{ "update:modelValue": [value: string] }>();
</script>

<template>
  <input
    :value="props.modelValue ?? ''"
    :type="props.type"
    data-slot="input"
    :class="
      cn(
        inputVariants({ size: props.size, invalid: props.invalid }),
        props.class
      )
    "
    @input="
      (e) =>
        $emit('update:modelValue', (e.target as HTMLInputElement).value)
    "
  />
</template>

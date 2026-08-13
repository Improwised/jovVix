<script setup>
import { computed } from "vue";
import { Check, ChevronDown } from "lucide-vue-next";
import {
  SelectContent,
  SelectItem,
  SelectItemIndicator,
  SelectItemText,
  SelectPortal,
  SelectRoot,
  SelectTrigger,
  SelectViewport,
} from "reka-ui";
import { cn } from "@/lib/utils";

const props = defineProps({
  modelValue: { type: String, default: "" },
  // Either {value, label} objects or plain strings.
  options: { type: Array, required: true },
  placeholder: { type: String, default: "Choose…" },
  disabled: { type: Boolean, default: false },
  class: { type: String, default: "" },
});

const emit = defineEmits(["update:modelValue"]);

const items = computed(() =>
  props.options.map((option) =>
    typeof option === "string" ? { value: option, label: option } : option
  )
);

// SelectValue only learns an item's text once the portal content has mounted, so
// the label is resolved here instead and shows before the first open.
const selectedLabel = computed(
  () => items.value.find((item) => item.value === props.modelValue)?.label || ""
);
</script>

<template>
  <SelectRoot
    :model-value="modelValue"
    :disabled="disabled"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <SelectTrigger
      :class="
        cn(
          'flex h-14 w-full items-center justify-between gap-3 border-[3px] border-jv-ink bg-jv-canvas px-4 text-left text-[17px] font-semibold text-jv-ink outline-none transition-shadow focus:shadow-brutal-sm disabled:cursor-not-allowed disabled:opacity-60',
          props.class
        )
      "
    >
      <span v-if="selectedLabel" class="truncate">{{ selectedLabel }}</span>
      <span v-else class="truncate text-jv-muted">{{ placeholder }}</span>
      <ChevronDown class="size-5 shrink-0 text-jv-ink" :stroke-width="2.6" />
    </SelectTrigger>

    <SelectPortal>
      <SelectContent
        position="popper"
        :side-offset="6"
        class="z-[1100] max-h-[280px] w-[var(--reka-select-trigger-width)] overflow-hidden jv-border-rough bg-jv-white shadow-brutal"
      >
        <SelectViewport class="max-h-[274px] overflow-y-auto p-1">
          <SelectItem
            v-for="item in items"
            :key="item.value"
            :value="item.value"
            class="flex cursor-pointer items-center justify-between gap-3 px-3 py-2.5 text-[16px] font-semibold text-jv-ink outline-none data-[highlighted]:bg-jv-yellow/50 data-[state=checked]:bg-jv-mint/50"
          >
            <SelectItemText class="truncate">{{ item.label }}</SelectItemText>
            <SelectItemIndicator>
              <Check class="size-4 shrink-0" :stroke-width="3" />
            </SelectItemIndicator>
          </SelectItem>
        </SelectViewport>
      </SelectContent>
    </SelectPortal>
  </SelectRoot>
</template>

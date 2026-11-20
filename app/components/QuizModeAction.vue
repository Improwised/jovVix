<script setup>
defineProps({
  label: { type: String, required: true },
  secondary: Boolean,
  disabled: Boolean,
});
defineEmits(["click"]);
</script>
<template>
  <button
    type="button"
    :disabled="disabled"
    class="inline-flex h-8 min-w-0 w-full cursor-pointer items-center justify-center whitespace-nowrap rounded-full border-2 border-jv-ink px-1 text-[11px] font-bold leading-none shadow-brutal-sm transition-all duration-150 ease-out enabled:hover:-translate-y-0.5 enabled:hover:-rotate-1 enabled:hover:brightness-95 enabled:active:translate-x-[2px] enabled:active:translate-y-[2px] enabled:active:rotate-0 enabled:active:shadow-none focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-jv-coral/40 disabled:cursor-wait disabled:opacity-50 sm:text-[12px]"
    :class="secondary ? 'bg-jv-white text-jv-ink' : 'bg-jv-coral text-white'"
    @click="$emit('click')"
  >
    {{ label }}
  </button>
</template>

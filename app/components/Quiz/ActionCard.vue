<template>
  <button
    type="button"
    class="jv-border-rough group flex w-full max-w-[300px] flex-col items-center justify-between gap-4 bg-jv-white p-6 text-left shadow-brutal-sm transition-transform hover:rotate-[-1deg] hover:shadow-brutal active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:p-8"
    @click="$emit('cardClick')"
  >
    <div
      class="flex h-[180px] w-full items-center justify-center"
      :aria-hidden="!!$slots.image"
    >
      <slot name="image" />
    </div>
    <h3
      class="font-headings text-[20px] leading-tight text-jv-ink sm:text-[22px]"
    >
      {{ props.actionTitle }}
    </h3>
  </button>
</template>

<script setup>
const props = defineProps({
  actionTitle: {
    default: "",
    type: String,
    required: true,
  },
});

defineEmits(["cardClick"]);
</script>

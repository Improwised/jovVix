<template>
  <article
    class="relative min-h-[64px] sm:min-h-[72px] bg-jv-white px-3 sm:px-4 md:px-5 xl:px-7 py-3 text-center shadow-brutal-sm jv-border-rough"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-10px] h-2.5 w-6 -translate-x-1/2 rounded-[3px] border-2 border-jv-ink bg-jv-yellow"
      aria-hidden="true"
    ></span>
    <h3
      class="font-body text-[15px] sm:text-[17px] md:text-[19px] xl:text-[20px] font-black leading-tight text-jv-ink"
    >
      {{ title }}
    </h3>
    <p
      class="mt-1 text-[11px] sm:text-[14px] leading-tight text-jv-muted break-words"
    >
      {{ description }}
    </p>
  </article>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
  tiltClass: {
    type: String,
    default: "",
  },
});
</script>

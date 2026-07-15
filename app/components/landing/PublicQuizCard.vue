<template>
  <button
    type="button"
    class="group relative flex w-full cursor-pointer flex-col bg-jv-white p-3 text-left shadow-brutal-sm jv-border-rough transition-transform duration-150 hover:-translate-y-1 hover:scale-[1.03] focus-visible:outline-none focus-visible:ring-4 focus-visible:ring-jv-coral/40 disabled:cursor-wait disabled:opacity-70"
    :class="tiltClass"
    :disabled="starting"
    :aria-label="`Start quiz: ${title}`"
    @click="$emit('start-quiz')"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <span class="relative block border-[2px] border-jv-ink bg-jv-slate p-2">
      <span class="block h-[118px] overflow-hidden sm:h-[128px]">
        <img
          :src="image"
          :alt="title"
          class="size-full object-cover transition-transform duration-200 group-hover:scale-[1.05]"
        />
      </span>
      <span
        v-if="starting"
        class="absolute inset-0 grid place-items-center bg-jv-ink/50 text-[13px] font-black uppercase tracking-[0.14em] text-jv-white"
      >
        Starting...
      </span>
    </span>

    <span
      class="mt-3 block break-words font-body text-[15px] font-bold leading-tight text-jv-ink sm:text-[16px]"
    >
      {{ title }}
    </span>
  </button>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  image: {
    type: String,
    default: "",
  },
  tiltClass: {
    type: String,
    default: "",
  },
  starting: {
    type: Boolean,
    default: false,
  },
});

defineEmits(["start-quiz"]);
</script>

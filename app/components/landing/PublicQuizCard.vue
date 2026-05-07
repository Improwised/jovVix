<template>
  <article
    class="relative flex min-h-[280px] sm:min-h-[306px] flex-col bg-jv-white p-3 sm:p-4 md:p-5 shadow-brutal-sm jv-border-rough"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <div class="border-[2px] border-jv-ink bg-jv-slate p-2">
      <div class="relative h-[88px] sm:h-[98px] overflow-hidden">
        <img :src="image" :alt="title" class="size-full object-cover" />
      </div>
    </div>

    <h3
      class="mt-4 sm:mt-5 font-body text-[20px] sm:text-[22px] md:text-[24px] font-normal leading-tight text-jv-ink"
    >
      {{ title }}
    </h3>
    <p
      class="mt-2 sm:mt-3 min-h-[48px] sm:min-h-[52px] text-[12px] sm:text-[13px] leading-[1.5] text-jv-muted"
    >
      {{ description }}
    </p>

    <div
      class="mt-auto flex items-center justify-between border-t border-dashed border-jv-ink/20 pt-4 sm:pt-5"
    >
      <div class="flex min-w-0 items-center gap-2">
        <span
          class="size-6 sm:size-7 shrink-0 rounded-full border border-jv-ink bg-[linear-gradient(135deg,#f6d7b0_0_45%,#6b4f3a_46%_100%)]"
          aria-hidden="true"
        ></span>
        <span
          class="truncate text-[13px] sm:text-[14px] font-semibold text-jv-ink"
          >{{ author }}</span
        >
      </div>
      <NuxtLink
        to="/join"
        class="ml-2 inline-flex h-8 min-w-[54px] sm:min-w-[58px] items-center justify-center rounded-[8px] border-[3px] border-jv-ink !bg-jv-yellow px-2.5 sm:px-3 text-[13px] sm:text-[14px] font-bold !text-jv-ink no-underline shadow-brutal-sm"
      >
        Play
      </NuxtLink>
    </div>
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
  author: {
    type: String,
    required: true,
  },
  tiltClass: {
    type: String,
    default: "",
  },
  image: {
    type: String,
    default: "",
  },
});
</script>

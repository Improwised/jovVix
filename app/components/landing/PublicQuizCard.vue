<template>
  <article
    class="relative flex min-h-[306px] flex-col bg-jv-white p-5 shadow-brutal-sm jv-border-rough sm:p-6"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <div class="border-[2px] border-jv-ink bg-jv-slate p-2">
      <div class="relative h-[98px] overflow-hidden" :class="imageClass">
        <div
          v-if="variant === 'pop'"
          class="absolute inset-x-5 top-5 rotate-[-3deg] text-center"
        >
          <p
            class="font-headings text-2xl leading-none text-fuchsia-200 drop-shadow-[2px_2px_0_#713f12]"
          >
            Starlight
          </p>
          <p
            class="font-feature text-sm font-black uppercase tracking-[0.12em] text-yellow-200"
          >
            Lounge
          </p>
        </div>
        <div
          v-else-if="variant === 'global'"
          class="absolute inset-0 flex items-center justify-center"
        >
          <div
            class="size-16 rounded-full border-[3px] border-[#9b7a4a] bg-[#d9b784] shadow-brutal-sm"
          ></div>
          <div class="absolute h-[88px] w-2 rotate-[25deg] bg-[#9b7a4a]"></div>
        </div>
        <div
          v-else-if="variant === 'movie'"
          class="absolute inset-x-8 bottom-4 rounded-t-[6px] border-[3px] border-[#cf2f4b] bg-[#16325c] p-2 text-center shadow-brutal-sm"
        >
          <p class="font-feature text-[9px] font-black uppercase text-jv-coral">
            The Grand
          </p>
          <p
            class="font-feature text-[11px] font-black uppercase text-yellow-200"
          >
            Cinema Palace
          </p>
        </div>
        <div
          v-else
          class="absolute inset-0 flex items-end justify-center gap-3 pb-4"
        >
          <span
            class="grid size-16 place-items-center rounded-full border-[3px] border-[#b8860b] bg-[#ffd76d]"
          >
            <span
              class="size-8 rounded-full border-[3px] border-[#b8860b] bg-[#ffe9a6]"
            ></span>
          </span>
          <span
            class="h-12 w-20 rounded-[45%] border-[3px] border-[#6aa7a4] bg-[#9fd8dc]"
          ></span>
          <span
            class="h-10 w-12 rounded-[45%] border-[3px] border-[#b96d3b] bg-[#f4a15d]"
          ></span>
        </div>
      </div>
    </div>

    <h3
      class="mt-5 font-body text-[24px] font-normal leading-tight text-jv-ink"
    >
      {{ title }}
    </h3>
    <p class="mt-3 min-h-[52px] text-[13px] leading-[1.5] text-jv-muted">
      {{ description }}
    </p>

    <div
      class="mt-auto flex items-center justify-between border-t border-dashed border-jv-ink/20 pt-5"
    >
      <div class="flex min-w-0 items-center gap-2">
        <span
          class="size-7 shrink-0 rounded-full border border-jv-ink bg-[linear-gradient(135deg,#f6d7b0_0_45%,#6b4f3a_46%_100%)]"
          aria-hidden="true"
        ></span>
        <span class="truncate text-[14px] font-semibold text-jv-ink">{{
          author
        }}</span>
      </div>
      <NuxtLink
        to="/join"
        class="ml-2 inline-flex h-8 min-w-[58px] items-center justify-center rounded-[8px] border-[3px] border-jv-ink !bg-jv-yellow px-3 text-[14px] font-bold !text-jv-ink no-underline shadow-brutal-sm"
      >
        Play
      </NuxtLink>
    </div>
  </article>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
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
  variant: {
    type: String,
    required: true,
  },
  tiltClass: {
    type: String,
    default: "",
  },
});

const imageClass = computed(() => {
  if (props.variant === "pop") {
    return "bg-[radial-gradient(circle_at_bottom,#f59e0b_0_2px,transparent_3px),linear-gradient(160deg,#111827_0%,#312e81_45%,#1f2937_100%)] bg-[length:12px_12px,auto]";
  }

  if (props.variant === "global") {
    return "bg-[linear-gradient(90deg,rgba(90,59,30,0.28)_0_7px,transparent_7px_18px),linear-gradient(135deg,#f0d7b5_0%,#caa36f_50%,#f7e6c8_100%)]";
  }

  if (props.variant === "movie") {
    return "bg-[radial-gradient(circle_at_top,#334155_0%,#111827_45%,#27344d_100%)]";
  }

  return "bg-[radial-gradient(circle_at_20%_18%,#f97316_0_3px,transparent_4px),radial-gradient(circle_at_80%_24%,#22c55e_0_3px,transparent_4px),linear-gradient(135deg,#f9f0c8_0%,#d7f7d7_100%)] bg-[length:24px_24px,22px_22px,auto]";
});
</script>

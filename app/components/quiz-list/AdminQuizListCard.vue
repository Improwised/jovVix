<template>
  <article
    class="group relative flex min-h-[360px] flex-col bg-jv-white p-3 sm:p-4 shadow-brutal-sm jv-border-rough"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <div class="relative border-[2px] border-jv-ink bg-jv-slate p-2">
      <div class="relative h-[118px] overflow-hidden">
        <img :src="image" :alt="title" class="size-full object-cover" />
      </div>

      <div
        class="absolute inset-2 grid place-items-center bg-jv-white/70 opacity-0 backdrop-blur-sm transition-opacity group-hover:opacity-100 group-focus-within:opacity-100"
      >
        <div class="flex flex-col gap-3">
          <!-- <NuxtLink
            :to="viewUrl"
            class="inline-flex h-11 min-w-44 items-center justify-center rounded-[999px] border-[3px] border-jv-ink bg-jv-coral px-6 font-headings text-[16px] text-white no-underline shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none"
          >
            View Quiz
          </NuxtLink> -->
          <NavigationLink
            :url="viewUrl"
            url-name="View Quiz"
            class="bg-jv-coral text-white rounded-[999px]"
          />
          <NavigationLink
            :url-name="starting ? 'Starting...' : 'Start Quiz'"
            class="rounded-[999px]"
            :disabled="starting"
            @click="$emit('start-quiz')"
          />

          <!-- <button
            type="button"
            class="inline-flex h-11 min-w-44 items-center justify-center rounded-[999px] border-[3px] border-jv-ink bg-jv-yellow px-6 font-headings text-[16px] text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[-1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none disabled:cursor-not-allowed disabled:opacity-70"
            :disabled="starting"
            @click="$emit('start-quiz')"
          >
            {{ starting ? "Starting..." : "Start Quiz" }}
          </button> -->
        </div>
      </div>
    </div>

    <div class="relative mt-4 flex items-start justify-between gap-3">
      <h3 class="font-body text-[24px] font-black leading-tight text-jv-ink">
        {{ title }}
      </h3>
      <button
        type="button"
        class="grid size-8 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg]"
        aria-label="Open quiz actions"
        @click="actionsOpen = !actionsOpen"
      >
        <MoreVertical class="size-4" :stroke-width="2.5" />
      </button>

      <div
        v-if="actionsOpen"
        class="absolute right-0 top-10 z-20 w-32 rotate-[1deg] border-[3px] border-jv-ink bg-jv-yellow p-2 shadow-brutal-sm jv-card"
      >
        <button
          type="button"
          class="flex w-full items-center gap-2 border-b border-dashed border-jv-ink/25 px-1 py-2 text-left text-[14px] font-bold transition-colors hover:text-jv-coral"
          @click="handleShare"
        >
          <Share2 class="size-4" :stroke-width="2.4" />
          <span>Share</span>
        </button>
        <button
          type="button"
          class="flex w-full items-center gap-2 px-1 py-2 text-left text-[14px] font-bold transition-colors hover:text-jv-coral"
          @click="handleDelete"
        >
          <Trash2 class="size-4" :stroke-width="2.4" />
          <span>Delete</span>
        </button>
      </div>
    </div>
    <p class="mt-1 text-[13px] leading-[1.4] text-jv-muted">
      {{ createdAt }}
    </p>
    <p class="mt-3 min-h-[66px] text-[14px] leading-[1.55] text-jv-muted">
      {{ description }}
    </p>

    <div class="mt-auto border-t-2 border-dashed border-jv-ink/15 pt-3">
      <span
        class="inline-flex items-center gap-1.5 rounded-[5px] border border-jv-ink/30 bg-jv-white px-2.5 py-1 text-[13px] leading-none text-jv-muted shadow-[1px_1px_0_rgba(45,45,45,0.25)]"
      >
        <CircleHelp class="size-3.5" :stroke-width="2.2" />
        <span>{{ questionCount }} Questions</span>
      </span>
    </div>
  </article>
</template>

<script setup>
import { ref } from "vue";
import { CircleHelp, MoreVertical, Share2, Trash2 } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

defineProps({
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    default: "",
  },
  createdAt: {
    type: String,
    default: "",
  },
  questionCount: {
    type: Number,
    default: 0,
  },
  image: {
    type: String,
    default: "",
  },
  tiltClass: {
    type: String,
    default: "",
  },
  viewUrl: {
    type: String,
    required: true,
  },
  starting: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["share", "delete", "start-quiz"]);
const actionsOpen = ref(false);

const handleShare = () => {
  actionsOpen.value = false;
  emit("share");
};

const handleDelete = () => {
  actionsOpen.value = false;
  emit("delete");
};
</script>

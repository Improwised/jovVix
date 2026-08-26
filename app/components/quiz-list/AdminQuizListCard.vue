<template>
  <article
    class="group relative flex min-h-[360px] w-full max-w-[260px] flex-col bg-jv-white p-3.5 shadow-brutal-sm jv-border-rough"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <div class="relative border-[2px] border-jv-ink bg-jv-slate p-1.5">
      <div class="relative aspect-[5/4] overflow-hidden">
        <QuizCoverImage
          :src="image"
          :fallback="fallback"
          :alt="title"
          :width="320"
          :height="256"
          class="size-full object-cover"
        />
      </div>
    </div>

    <div
      ref="actionsMenuRef"
      class="relative mt-3.5 flex min-h-[44px] items-start justify-between gap-2"
    >
      <h3
        class="min-w-0 flex-1 break-words font-body text-[18px] font-black leading-[1.2] text-jv-ink line-clamp-2"
      >
        {{ title }}
      </h3>
      <button
        v-if="showActions"
        type="button"
        class="grid size-7 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg]"
        aria-label="Open quiz actions"
        :aria-expanded="actionsOpen"
        @click="toggleActionsMenu"
      >
        <MoreVertical class="size-3.5" :stroke-width="2.5" />
      </button>

      <div
        v-if="showActions && actionsOpen"
        class="absolute right-0 top-9 z-20 w-32 rotate-[1deg] border-[3px] border-jv-ink bg-jv-yellow p-2 shadow-brutal-sm jv-card"
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
    <div
      class="mt-3 flex items-center justify-between gap-2 border-y border-dashed border-jv-ink/20 py-2.5 text-[12px] leading-none text-jv-muted"
    >
      <span class="truncate">{{ createdAt }}</span>
      <span class="inline-flex shrink-0 items-center gap-1">
        <CircleHelp class="size-3" :stroke-width="2.2" />
        {{ questionCount }} Questions
      </span>
    </div>
    <div class="mt-auto pt-3.5">
      <div
        class="grid gap-1.5"
        :class="showActions ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <NavigationLink
          v-if="showActions"
          url-name="View Quiz"
          :url="viewUrl"
          class="h-8 rounded-full bg-jv-coral px-1 font-body text-[12px] font-semibold text-white shadow-none md:px-1 md:text-[14px]"
        />
        <NavigationLink
          url-name="Start Quiz"
          class="h-8 rounded-full px-1 font-body text-[12px] font-semibold shadow-none md:px-1 md:text-[14px]"
          @click="$emit('start-quiz')"
        />
      </div>
    </div>
  </article>
</template>

<script setup>
import { ref } from "vue";
import { onClickOutside } from "@vueuse/core";
import { CircleHelp, MoreVertical, Share2, Trash2 } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";
import QuizCoverImage from "@/components/QuizCoverImage.vue";

defineProps({
  title: {
    type: String,
    required: true,
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
  fallback: {
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
  showActions: {
    type: Boolean,
    default: true,
  },
});

const emit = defineEmits(["share", "delete", "start-quiz"]);
const actionsMenuRef = ref(null);
const actionsOpen = ref(false);

const closeActionsMenu = () => {
  actionsOpen.value = false;
};

const toggleActionsMenu = () => {
  actionsOpen.value = !actionsOpen.value;
};

onClickOutside(actionsMenuRef, closeActionsMenu);

const handleShare = () => {
  closeActionsMenu();
  emit("share");
};

const handleDelete = () => {
  closeActionsMenu();
  emit("delete");
};
</script>

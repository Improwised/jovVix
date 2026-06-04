<template>
  <article
    class="group relative flex min-h-[342px] flex-col bg-jv-white p-3 shadow-brutal-sm jv-border-rough sm:min-h-[360px] sm:p-4 md:min-h-[372px] md:p-5"
    :class="tiltClass"
  >
    <span
      class="absolute left-1/2 top-[-13px] h-4 w-8 -translate-x-1/2 rotate-[-2deg] bg-jv-salmon opacity-90"
      aria-hidden="true"
    ></span>

    <div class="relative border-[2px] border-jv-ink bg-jv-slate p-2">
      <div class="relative h-[104px] overflow-hidden sm:h-[112px] md:h-[124px]">
        <img :src="image" :alt="title" class="size-full object-cover" />
      </div>
      <span
        v-if="isPublic"
        class="absolute right-[-6px] top-[-10px] z-10 inline-flex rotate-[4deg] items-center gap-1 border-[2.5px] border-jv-ink bg-jv-mint px-2 py-[3px] text-[11px] font-black uppercase tracking-[0.12em] text-jv-ink shadow-[2px_2px_0_#2D2D2D]"
        aria-label="Public quiz"
      >
        <Globe class="size-3" :stroke-width="2.6" />
        Public
      </span>
    </div>

    <div
      ref="actionsMenuRef"
      class="relative mt-4 flex items-start justify-between gap-3"
    >
      <h3
        class="min-w-0 break-words font-body text-[21px] font-black leading-tight text-jv-ink sm:text-[22px] md:text-[24px]"
      >
        {{ title }}
      </h3>
      <button
        v-if="showActions"
        type="button"
        class="grid size-8 shrink-0 place-items-center border-2 border-jv-ink bg-jv-white text-jv-ink shadow-[1px_1px_0_#2D2D2D] transition-transform hover:rotate-[3deg]"
        aria-label="Open quiz actions"
        :aria-expanded="actionsOpen"
        @click="toggleActionsMenu"
      >
        <MoreVertical class="size-4" :stroke-width="2.5" />
      </button>

      <div
        v-if="showActions && actionsOpen"
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
    <p class="mt-1 text-[12px] leading-[1.4] text-jv-muted sm:text-[13px]">
      {{ createdAt }}
    </p>
    <p
      class="mt-2 min-h-[58px] break-words text-[13px] leading-[1.5] text-jv-muted sm:mt-3 sm:min-h-[66px] sm:text-[14px] sm:leading-[1.55]"
    >
      {{ description }}
    </p>

    <span
      class="inline-flex max-w-full items-center gap-1.5 px-2.5 py-1 text-[12px] leading-none text-jv-muted sm:text-[13px]"
    >
      <CircleHelp class="size-3.5" :stroke-width="2.2" />
      <span class="truncate">{{ questionCount }} Questions</span>
    </span>
    <div class="mt-3 border-t-2 border-dashed border-jv-ink/15 pt-3">
      <div
        class="mt-3 grid gap-2"
        :class="showActions ? 'grid-cols-2' : 'grid-cols-1'"
      >
        <NavigationLink
          v-if="showActions"
          url-name="View Quiz"
          :url="viewUrl"
          class="h-8 rounded-full bg-jv-coral text-white shadow-none"
        />
        <NavigationLink
          url-name="Start Quiz"
          class="h-8 rounded-full shadow-none"
          @click="$emit('start-quiz')"
        />
      </div>
    </div>
  </article>
</template>

<script setup>
import { ref } from "vue";
import { onClickOutside } from "@vueuse/core";
import {
  CircleHelp,
  Globe,
  MoreVertical,
  Share2,
  Trash2,
} from "lucide-vue-next";
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
  showActions: {
    type: Boolean,
    default: true,
  },
  isPublic: {
    type: Boolean,
    default: false,
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

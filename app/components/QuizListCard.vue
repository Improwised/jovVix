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
    </div>

    <h3
      class="mt-4 min-w-0 break-words font-body text-[21px] font-black leading-tight text-jv-ink sm:text-[22px] md:text-[24px]"
    >
      {{ title }}
    </h3>

    <p class="mt-1 text-[12px] leading-[1.4] text-jv-muted sm:text-[13px]">
      {{ createdAtLabel }}
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
        :class="isPlayedQuiz ? 'grid-cols-1' : 'grid-cols-2'"
      >
        <NavigationLink
          url-name="View Quiz"
          :url="viewUrl"
          class="h-8 rounded-full bg-jv-coral text-white shadow-none"
        />
        <UtilsStartQuiz v-if="!isPlayedQuiz" :quiz-id="props.details?.id" />
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from "vue";
import { CircleHelp } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

const props = defineProps({
  details: {
    type: Object,
    default: () => ({}),
    required: true,
  },
  isPlayedQuiz: {
    type: Boolean,
    required: false,
    default: false,
  },
});

const quizImages = [
  "/images/landing/homepage-public-quiz-1.png",
  "/images/landing/homepage-public-quiz-2.png",
  "/images/landing/homepage-public-quiz-3.png",
  "/images/landing/homepage-public-quiz-4.png",
];

const tiltClasses = [
  "rotate-[-0.8deg]",
  "rotate-[0.7deg]",
  "rotate-[-0.4deg]",
  "rotate-[0.5deg]",
];

const hashIndex = (value, length) => {
  if (!value) return 0;
  const str = String(value);
  let hash = 0;
  for (let i = 0; i < str.length; i += 1) {
    hash = (hash * 31 + str.charCodeAt(i)) >>> 0;
  }
  return hash % length;
};

const variantIndex = computed(() =>
  hashIndex(props.details?.id, quizImages.length)
);

const title = computed(() => decodeURI(props.details?.title || ""));
const description = computed(() => props.details?.description?.String || "");
const questionCount = computed(() => props.details?.total_questions || 0);
const image = computed(() => quizImages[variantIndex.value]);
const tiltClass = computed(() => tiltClasses[variantIndex.value]);
const createdAtLabel = computed(() =>
  props.details?.created_at ? useGetTime(props.details.created_at) : ""
);
const viewUrl = computed(() =>
  props.isPlayedQuiz
    ? `/admin/played_quiz/${props.details?.id}`
    : `/admin/quiz/list-quiz/${props.details?.id}`
);
</script>

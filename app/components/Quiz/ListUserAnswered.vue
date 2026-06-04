<script setup>
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import { useNuxtApp } from "nuxt/app";
import { usePush } from "notivue";
import { Check, Users } from "lucide-vue-next";
import { getAvatarUrlByName } from "~~/composables/avatar";

const app = useNuxtApp();
const toast = usePush();

const usersThatSubmittedAnswer = useUserThatSubmittedAnswer();
const { usersSubmittedAnswers } = usersThatSubmittedAnswer;
const totalUser = ref(0);

const emits = defineEmits(["autoSkip"]);

const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
  runningQuizJoinUser: {
    type: Number,
    required: false,
    default: 0,
  },
});

watch(
  () => props.data,
  (message) => {
    if (message.status == app.$Fail) {
      toast.error(message.data);
      return;
    }
    handleCountUser(message);
  },
  { deep: true, immediate: true }
);

function handleCountUser(message) {
  if (message.event == app.$GetQuestion) {
    totalUser.value = message.data.totalJoinUser;
  }
}

watch(
  [totalUser, usersSubmittedAnswers],
  () => {
    if (totalUser.value <= usersSubmittedAnswers.length) {
      emits("autoSkip");
    }
  },
  { immediate: true, deep: true }
);

watch(
  () => props.runningQuizJoinUser,
  (message) => {
    if (message && totalUser.value < message) {
      totalUser.value = message;
    }
  },
  { deep: true, immediate: true }
);

const answeredCount = computed(() => usersSubmittedAnswers.length);
const chipAccents = [
  "bg-jv-yellow",
  "bg-jv-mint",
  "bg-jv-salmon",
  "bg-jv-white",
];
</script>

<template>
  <section
    class="mx-auto mt-6 flex w-full max-w-[1180px] flex-col items-center px-3 pb-8 sm:mt-8 sm:px-6 sm:pb-10 md:px-10 md:pb-12"
  >
    <div
      v-if="answeredCount === 0"
      class="inline-flex items-center gap-2 rounded-full border-[2px] border-dashed border-jv-ink/30 px-4 py-2 font-body text-[13px] font-bold text-jv-muted sm:gap-2.5 sm:px-5 sm:py-2.5 sm:text-[14px]"
    >
      <span
        class="size-2 shrink-0 rounded-full bg-jv-ink/30 sm:size-2.5"
        aria-hidden="true"
      ></span>
      No one answered till now
    </div>

    <div v-else class="flex flex-col items-center gap-5 sm:gap-6">
      <div
        class="inline-flex items-center gap-3 -rotate-[0.4deg] jv-border-rough bg-jv-white px-4 py-2.5 shadow-brutal-sm sm:gap-4 sm:px-5 sm:py-3"
      >
        <span
          class="grid size-9 shrink-0 rotate-[-3deg] place-items-center border-[3px] border-jv-ink bg-jv-mint sm:size-11"
        >
          <Users class="size-4 sm:size-5" :stroke-width="2.4" />
        </span>
        <p class="font-body text-[14px] font-bold text-jv-muted sm:text-[16px]">
          <span
            class="font-feature text-[18px] font-black text-jv-ink sm:text-[20px]"
          >
            {{ answeredCount
            }}<span class="text-jv-muted">/{{ totalUser }}</span>
          </span>
          players answered
        </p>
      </div>

      <ul class="flex w-full flex-wrap justify-center gap-3 sm:gap-4">
        <li
          v-for="(user, index) in usersSubmittedAnswers"
          :key="user.UserId || user.username || index"
          :class="[
            'flex items-center gap-3 border-[3px] border-jv-ink px-3 py-2 shadow-brutal-sm sm:gap-3 sm:px-4 sm:py-2.5',
            chipAccents[index % chipAccents.length],
            index % 2 === 0 ? 'rotate-[-0.6deg]' : 'rotate-[0.6deg]',
          ]"
        >
          <span
            class="grid size-10 shrink-0 place-items-center overflow-hidden rounded-full border-[2px] border-jv-ink bg-jv-white sm:size-11"
          >
            <img
              :src="getAvatarUrlByName(user?.img_key)"
              :alt="`${user?.first_name || user?.username || 'Player'} avatar`"
              class="size-full object-cover"
            />
          </span>
          <div class="flex min-w-0 flex-col">
            <span
              class="truncate font-body text-[14px] font-black leading-tight text-jv-ink sm:text-[16px]"
            >
              {{ user.first_name }}
            </span>
            <span
              v-if="user.username"
              class="truncate font-body text-[11px] font-bold text-jv-muted sm:text-[12px]"
            >
              @{{ user.username }}
            </span>
          </div>
          <span
            class="grid size-7 shrink-0 rotate-[3deg] place-items-center rounded-full border-[2px] border-jv-ink bg-jv-mint sm:size-8"
            aria-label="Answered"
          >
            <Check class="size-4" :stroke-width="2.8" />
          </span>
        </li>
      </ul>
    </div>
  </section>
</template>

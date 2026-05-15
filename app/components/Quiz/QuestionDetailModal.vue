<script setup>
import { ArrowLeft, CheckCircle2, XCircle, X } from "lucide-vue-next";

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  question: {
    type: Object,
    default: null,
  },
  index: {
    type: Number,
    default: 0,
  },
  rows: {
    type: Array,
    default: () => [],
  },
});

const emit = defineEmits(["close"]);

const close = () => emit("close");

const correctIndices = computed(() => {
  const list = props.question?.correct_answer;
  if (!list) return new Set();
  if (Array.isArray(list)) return new Set(list.map((n) => Number(n)));
  return new Set(
    String(list)
      .split(",")
      .map((s) => Number(s.trim()))
      .filter((n) => !Number.isNaN(n))
  );
});

const optionTextFor = (idx) => {
  const opts = props.question?.options;
  if (!opts) return "";
  if (Array.isArray(opts)) return opts[idx] ?? "";
  if (typeof opts === "object") return opts[String(idx)] ?? opts[idx] ?? "";
  return "";
};

// Index analytics_board rows by username for fast enrichment lookup.
const rowByUsername = computed(() => {
  const map = {};
  for (const row of props.rows || []) {
    if (row?.username) map[row.username] = row;
  }
  return map;
});

// Primary source: quiz.selected_answers is { "0": ["userA","userB"], "1": [...] }
// where the key is the option index the user picked. Build the full user list
// from this so the modal works even if analytics_board hasn't loaded yet.
const buckets = computed(() => {
  const correct = [];
  const incorrect = [];
  const selected = props.question?.selected_answers;
  if (!selected || typeof selected !== "object") return { correct, incorrect };

  for (const [optKey, usernames] of Object.entries(selected)) {
    if (!Array.isArray(usernames)) continue;
    const selectedIdx = Number(optKey);
    if (Number.isNaN(selectedIdx)) continue;
    const letter = String.fromCharCode(64 + selectedIdx);
    const optText = optionTextFor(selectedIdx);
    const isCorrect = correctIndices.value.has(selectedIdx);

    for (const username of usernames) {
      const enrich = rowByUsername.value[username];
      const responseMs = enrich?.response_time;
      const entry = {
        username,
        firstname: enrich?.firstname || username,
        avatar: getAvatarUrlByName(enrich?.img_key),
        letter,
        optionText: optText,
        responseTimeSec:
          typeof responseMs === "number"
            ? (responseMs / 1000).toFixed(2)
            : null,
        isCorrect,
      };
      if (isCorrect) correct.push(entry);
      else incorrect.push(entry);
    }
  }
  return { correct, incorrect };
});

const correctCount = computed(() => buckets.value.correct.length);
const incorrectCount = computed(() => buckets.value.incorrect.length);
const participants = computed(
  () =>
    props.question?.userParticipants ??
    correctCount.value + incorrectCount.value
);
const avgTimeSec = computed(() => {
  const ms = props.question?.avg_response_time || 0;
  return Math.abs(ms / 1000).toFixed(2);
});

const onKeydown = (e) => {
  if (e.key === "Escape" && props.isOpen) close();
};

watch(
  () => props.isOpen,
  (open) => {
    if (typeof document === "undefined") return;
    if (open) {
      document.addEventListener("keydown", onKeydown);
      document.body.style.overflow = "hidden";
    } else {
      document.removeEventListener("keydown", onKeydown);
      document.body.style.overflow = "";
    }
  }
);

onBeforeUnmount(() => {
  if (typeof document !== "undefined") {
    document.removeEventListener("keydown", onKeydown);
    document.body.style.overflow = "";
  }
});
</script>

<template>
  <Teleport to="body">
    <div
      v-if="isOpen && question"
      class="fixed inset-0 z-50 grid place-items-center bg-jv-ink/35 px-4 py-6 backdrop-blur-[2px]"
      role="dialog"
      aria-modal="true"
      @click.self="close"
    >
      <div
        class="flex max-h-[92vh] w-full max-w-[1000px] rotate-[-0.4deg] flex-col border-[4px] border-jv-ink bg-jv-white shadow-brutal-lg"
      >
        <!-- Dark header: back button + 4-stat strip -->
        <div
          class="flex flex-col gap-4 border-b-[3px] border-jv-ink bg-jv-ink px-5 py-4 text-jv-white sm:flex-row sm:items-center sm:justify-between sm:px-6"
        >
          <button
            type="button"
            class="inline-flex items-center gap-2 self-start text-[15px] font-bold text-jv-white transition-transform hover:translate-x-[-2px] sm:text-[16px]"
            @click="close"
          >
            <ArrowLeft class="size-5" :stroke-width="2.6" />
            Back to all questions
          </button>

          <div
            class="grid grid-cols-2 gap-x-5 gap-y-2 border-[3px] border-jv-white/20 bg-jv-ink/60 px-4 py-2 sm:grid-cols-4 sm:gap-x-7"
          >
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Participants
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-white">
                {{ participants }}
              </div>
            </div>
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Avg. time
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-white">
                {{ avgTimeSec }} sec
              </div>
            </div>
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Correct
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-accent-green">
                {{ correctCount }}
              </div>
            </div>
            <div>
              <div
                class="text-[10px] font-black uppercase tracking-[0.14em] text-jv-white/60"
              >
                Wrong
              </div>
              <div class="mt-0.5 text-[18px] font-black text-jv-coral">
                {{ incorrectCount }}
              </div>
            </div>
          </div>

          <button
            type="button"
            class="absolute right-3 top-3 grid size-9 place-items-center text-jv-white transition-transform hover:rotate-[6deg] sm:hidden"
            aria-label="Close"
            @click="close"
          >
            <X class="size-6" :stroke-width="2.4" />
          </button>
        </div>

        <!-- Question title bar -->
        <div
          class="flex items-center gap-3 border-b-[3px] border-jv-ink bg-jv-white px-5 py-4 sm:px-6"
        >
          <span
            class="inline-flex items-center justify-center jv-border-rough bg-jv-yellow px-3 py-1.5 text-[14px] font-black text-jv-ink shadow-brutal-sm"
          >
            Q{{ index }}
          </span>
          <h2
            class="font-headings text-[22px] leading-tight text-jv-ink sm:text-[26px]"
          >
            {{ question.question }}
          </h2>
        </div>

        <!-- Two-column body -->
        <div class="grid flex-1 grid-cols-1 overflow-hidden md:grid-cols-2">
          <!-- Correct column -->
          <section
            class="flex min-h-0 flex-col border-b-[3px] border-jv-ink bg-[#eafff2] md:border-b-0 md:border-r-[3px]"
          >
            <header
              class="flex items-center justify-between gap-3 border-b-2 border-jv-ink/15 px-5 py-3"
            >
              <div class="inline-flex items-center gap-2">
                <CheckCircle2
                  class="size-5 text-jv-accent-green"
                  :stroke-width="2.6"
                />
                <span class="text-[15px] font-black text-jv-ink"
                  >Answered correctly</span
                >
              </div>
              <span
                class="jv-border-rough bg-jv-white px-2.5 py-0.5 text-[13px] font-black text-jv-ink shadow-brutal-sm"
              >
                {{ correctCount }} users
              </span>
            </header>

            <div class="flex-1 overflow-y-auto px-4 py-3 sm:px-5">
              <div
                v-if="!correctCount"
                class="py-8 text-center text-[14px] font-semibold text-jv-muted"
              >
                No one answered correctly.
              </div>
              <ul v-else class="flex flex-col gap-2.5">
                <li
                  v-for="(u, i) in buckets.correct"
                  :key="`c-${u.username}-${i}`"
                  class="flex items-center gap-3 jv-border-rough bg-jv-white px-3 py-2.5 shadow-brutal-sm"
                >
                  <img
                    :src="u.avatar"
                    :alt="u.firstname"
                    class="size-10 shrink-0 rounded-full border-[2px] border-jv-ink bg-jv-slate object-cover"
                  />
                  <div class="min-w-0 flex-1">
                    <div
                      class="truncate text-[15px] font-black text-jv-ink leading-tight"
                    >
                      {{ u.firstname }}
                    </div>
                    <div
                      class="mt-0.5 truncate text-[12px] font-semibold text-jv-muted"
                    >
                      Selected {{ u.letter
                      }}<template v-if="u.optionText">
                        · {{ u.optionText }}</template
                      ><template v-if="u.responseTimeSec">
                        in {{ u.responseTimeSec }} sec</template
                      >
                    </div>
                  </div>
                  <span
                    class="inline-flex shrink-0 items-center gap-1 text-[13px] font-black text-jv-accent-green"
                  >
                    <CheckCircle2 class="size-4" :stroke-width="2.6" />
                    Correct
                  </span>
                </li>
              </ul>
            </div>
          </section>

          <!-- Incorrect column -->
          <section class="flex min-h-0 flex-col bg-[#fff1f1]">
            <header
              class="flex items-center justify-between gap-3 border-b-2 border-jv-ink/15 px-5 py-3"
            >
              <div class="inline-flex items-center gap-2">
                <XCircle class="size-5 text-jv-coral" :stroke-width="2.6" />
                <span class="text-[15px] font-black text-jv-ink"
                  >Answered incorrectly</span
                >
              </div>
              <span
                class="jv-border-rough bg-jv-white px-2.5 py-0.5 text-[13px] font-black text-jv-ink shadow-brutal-sm"
              >
                {{ incorrectCount }} users
              </span>
            </header>

            <div class="flex-1 overflow-y-auto px-4 py-3 sm:px-5">
              <div
                v-if="!incorrectCount"
                class="py-8 text-center text-[14px] font-semibold text-jv-muted"
              >
                No one answered incorrectly.
              </div>
              <ul v-else class="flex flex-col gap-2.5">
                <li
                  v-for="(u, i) in buckets.incorrect"
                  :key="`i-${u.username}-${i}`"
                  class="flex items-center gap-3 jv-border-rough bg-jv-white px-3 py-2.5 shadow-brutal-sm"
                >
                  <img
                    :src="u.avatar"
                    :alt="u.firstname"
                    class="size-10 shrink-0 rounded-full border-[2px] border-jv-ink bg-jv-slate object-cover"
                  />
                  <div class="min-w-0 flex-1">
                    <div
                      class="truncate text-[15px] font-black text-jv-ink leading-tight"
                    >
                      {{ u.firstname }}
                    </div>
                    <div
                      class="mt-0.5 truncate text-[12px] font-semibold text-jv-muted"
                    >
                      Picked {{ u.letter
                      }}<template v-if="u.optionText">
                        · {{ u.optionText }}</template
                      ><template v-if="u.responseTimeSec">
                        in {{ u.responseTimeSec }} sec</template
                      >
                    </div>
                  </div>
                  <span
                    class="inline-flex shrink-0 items-center gap-1 text-[13px] font-black text-jv-coral"
                  >
                    <XCircle class="size-4" :stroke-width="2.6" />
                    {{ u.letter }}
                  </span>
                </li>
              </ul>
            </div>
          </section>
        </div>
      </div>
    </div>
  </Teleport>
</template>

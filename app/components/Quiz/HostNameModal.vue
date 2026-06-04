<script setup>
import { ref, computed, onMounted, onUnmounted } from "vue";
import { User, RefreshCw, ArrowRight, Loader2, X } from "lucide-vue-next";
import { getRandomAvatarName, getAvatarUrlByName } from "~~/composables/avatar";
import NavigationLink from "@/components/common/NavigationLink.vue";

const props = defineProps({
  submitting: {
    type: Boolean,
    default: false,
  },
  initialName: {
    type: String,
    default: "",
  },
  allowClose: {
    type: Boolean,
    default: false,
  },
});

const emits = defineEmits(["submit", "close"]);

const username = ref((props.initialName || "").slice(0, 12));
const avatarName = ref(getRandomAvatarName());
const avatarUrl = computed(() => getAvatarUrlByName(avatarName.value));
const visible = ref(false);

const trimmedName = computed(() =>
  (username.value || "").trim().replace(/\s+/g, "")
);
const canSubmit = computed(
  () =>
    !props.submitting &&
    trimmedName.value.length > 0 &&
    trimmedName.value.length <= 12
);

const rerollAvatar = () => {
  avatarName.value = getRandomAvatarName();
};

const handleSubmit = () => {
  if (!canSubmit.value) return;
  emits("submit", {
    name: trimmedName.value,
    avatarName: avatarName.value,
  });
};

const handleClose = () => {
  if (!props.allowClose || props.submitting) return;
  emits("close");
};

const onKeydown = (e) => {
  if (e.key === "Escape") handleClose();
};

onMounted(() => {
  visible.value = true;
  if (process.client) {
    document.addEventListener("keydown", onKeydown);
    document.body.style.overflow = "hidden";
  }
});

onUnmounted(() => {
  if (process.client) {
    document.removeEventListener("keydown", onKeydown);
    document.body.style.overflow = "";
  }
});
</script>

<template>
  <Teleport to="body">
    <Transition name="jv-modal">
      <div
        v-if="visible"
        class="fixed inset-0 z-[1000] flex items-center justify-center px-4 py-6"
        role="dialog"
        aria-modal="true"
        aria-labelledby="hostNameModalLabel"
      >
        <button
          type="button"
          class="absolute inset-0 cursor-default bg-jv-ink/60"
          :disabled="!allowClose || submitting"
          aria-label="Close"
          @click="handleClose"
        ></button>

        <div
          class="jv-modal-card relative z-10 w-full max-w-[440px] rotate-1 jv-card border-2 border-jv-ink bg-jv-white px-6 py-7 shadow-brutal-lg sm:px-8 sm:py-9"
        >
          <span
            class="absolute left-1/2 -top-[8px] z-20 h-3 w-12 -translate-x-1/2 jv-card border-2 border-jv-ink bg-jv-slate shadow-brutal-sm"
            aria-hidden="true"
          ></span>

          <button
            v-if="allowClose"
            type="button"
            class="absolute right-3 top-3 grid size-9 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:right-4 sm:top-4"
            :disabled="submitting"
            aria-label="Close"
            @click="handleClose"
          >
            <X class="size-4" :stroke-width="2.6" />
          </button>

          <header class="mb-6 flex flex-col items-center gap-1.5">
            <div class="relative inline-block">
              <h1
                id="hostNameModalLabel"
                class="m-0 font-headings text-[28px] leading-none text-jv-ink sm:text-[34px]"
              >
                Pick Your Name
              </h1>
              <svg
                class="absolute -bottom-2 left-1/2 -translate-x-1/2 text-jv-yellow"
                width="120"
                height="14"
                viewBox="0 0 120 14"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                aria-hidden="true"
              >
                <path
                  d="M3 9 Q 20 1, 40 7 T 78 6 T 117 4"
                  stroke="currentColor"
                  stroke-width="2.5"
                  stroke-linecap="round"
                  fill="none"
                />
              </svg>
            </div>
            <p class="m-0 mt-1 font-body text-sm text-jv-ink/70 sm:text-base">
              Players will see this name on the leaderboard
            </p>
          </header>

          <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
            <div class="flex flex-col gap-1.5">
              <label
                for="hostName"
                class="px-0.5 font-body text-xs font-bold uppercase tracking-wide text-jv-ink sm:text-[13px]"
              >
                Your Player Profile
              </label>
              <div class="flex items-stretch gap-2.5">
                <NavigationLink
                  type="button"
                  class="relative size-[46px] shrink-0 overflow-hidden bg-jv-mint !p-0 shadow-none !jv-card !rounded-none hover:rotate-0 hover:scale-110"
                  :aria-label="`Generate new avatar (current: ${avatarName})`"
                  @click="rerollAvatar"
                >
                  <img
                    :src="avatarUrl"
                    :alt="avatarName"
                    class="absolute inset-0 size-full object-cover"
                  />
                  <span
                    class="absolute -bottom-1 -right-1 grid size-5 place-items-center rounded-full border-2 border-jv-canvas bg-jv-ink text-white"
                    aria-hidden="true"
                  >
                    <RefreshCw class="size-2.5" :stroke-width="2.5" />
                  </span>
                </NavigationLink>

                <div
                  class="flex flex-1 items-center gap-2.5 border-2 border-jv-ink bg-jv-white px-3 jv-card shadow-brutal-sm transition-all focus-within:translate-x-[1px] focus-within:translate-y-[1px] focus-within:shadow-none"
                >
                  <User
                    class="size-[18px] shrink-0 text-jv-ink/70"
                    :stroke-width="2.2"
                  />
                  <input
                    id="hostName"
                    v-model.trim="username"
                    type="text"
                    name="username"
                    maxlength="12"
                    placeholder="Pick a name"
                    autocomplete="off"
                    class="min-w-0 flex-1 border-0 bg-transparent font-body text-sm text-jv-ink outline-none placeholder:text-jv-ink/40 sm:text-base"
                  />
                </div>
              </div>
            </div>

            <NavigationLink
              type="submit"
              class="mt-1 w-full bg-jv-coral py-2.5 text-sm text-white sm:py-3 sm:text-base"
              :disabled="!canSubmit"
            >
              <template v-if="submitting">
                <Loader2 class="size-[18px] animate-spin" :stroke-width="2.4" />
                <span>Setting up…</span>
              </template>
              <template v-else>
                <span>Start Hosting</span>
                <ArrowRight class="size-[18px]" :stroke-width="2.4" />
              </template>
            </NavigationLink>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.jv-modal-enter-active,
.jv-modal-leave-active {
  transition: opacity 0.2s ease;
}
.jv-modal-enter-active .jv-modal-card,
.jv-modal-leave-active .jv-modal-card {
  transition: transform 0.25s cubic-bezier(0.34, 1.56, 0.64, 1),
    opacity 0.2s ease;
}
.jv-modal-enter-from,
.jv-modal-leave-to {
  opacity: 0;
}
.jv-modal-enter-from .jv-modal-card,
.jv-modal-leave-to .jv-modal-card {
  opacity: 0;
  transform: scale(0.92) rotate(-2deg);
}
</style>

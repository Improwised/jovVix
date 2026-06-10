<script setup lang="js">
import { AlertTriangle, X } from "lucide-vue-next";

const props = defineProps({
  modalTitle: {
    type: String,
    default: "This is a modal Title",
    required: true,
  },
  modalMessage: {
    type: String,
    default: "Hello this is a modal message",
    required: true,
  },
  modelPositiveMessage: {
    type: String,
    default: "Save",
    required: false,
  },
  modelNegativeMessage: {
    type: String,
    default: "Cancel",
    required: false,
  },
});

const emits = defineEmits(["confirmMessage"]);

const isSent = ref(false);
const visible = ref(false);

const handleClose = (confirm) => {
  if (isSent.value) return;
  isSent.value = true;
  visible.value = false;
  emits("confirmMessage", confirm);
};

const onKeydown = (e) => {
  if (e.key === "Escape") handleClose(false);
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
        :aria-labelledby="'confirmModalLabel'"
      >
        <button
          type="button"
          class="absolute inset-0 cursor-default bg-jv-ink/60"
          aria-label="Close"
          @click="handleClose(false)"
        ></button>

        <div
          class="jv-modal-card relative z-10 w-full max-w-[480px] -rotate-[0.4deg] jv-border-rough bg-jv-white p-5 shadow-brutal-lg sm:p-7"
        >
          <span
            class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
            aria-hidden="true"
          ></span>

          <button
            type="button"
            class="absolute right-3 top-3 grid size-9 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:right-4 sm:top-4"
            aria-label="Close"
            @click="handleClose(false)"
          >
            <X class="size-4" :stroke-width="2.6" />
          </button>

          <div class="flex items-start gap-3 pr-10 sm:gap-4 sm:pr-12">
            <span
              class="grid size-11 shrink-0 -rotate-[3deg] place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-yellow text-jv-ink shadow-brutal-sm sm:size-12"
              aria-hidden="true"
            >
              <AlertTriangle class="size-5 sm:size-6" :stroke-width="2.4" />
            </span>
            <div class="min-w-0 flex-1">
              <h2
                id="confirmModalLabel"
                class="font-headings text-[22px] leading-tight text-jv-ink sm:text-[28px]"
              >
                {{ props.modalTitle }}
              </h2>
              <p
                class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[15px]"
              >
                {{ props.modalMessage }}
              </p>
            </div>
          </div>

          <div
            class="mt-5 border-t-[2px] border-dashed border-jv-ink/25 sm:mt-6"
            aria-hidden="true"
          ></div>

          <div
            class="mt-4 flex flex-col-reverse items-stretch justify-end gap-3 sm:mt-5 sm:flex-row sm:items-center"
          >
            <button
              type="button"
              class="inline-flex h-10 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-6 sm:text-[15px]"
              @click="handleClose(false)"
            >
              {{ props.modelNegativeMessage }}
            </button>
            <button
              type="button"
              class="inline-flex h-10 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-coral px-5 font-body text-[14px] font-black text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-6 sm:text-[15px]"
              @click="handleClose(true)"
            >
              {{ props.modelPositiveMessage }}
            </button>
          </div>
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

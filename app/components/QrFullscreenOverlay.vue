<script setup>
import { X } from "lucide-vue-next";

const props = defineProps({
  joinUrl: {
    type: String,
    required: true,
  },
  code: {
    type: String,
    default: "",
  },
});

const isOpen = ref(false);
const overlayEl = ref(null);

const displayUrl = computed(() =>
  (props.joinUrl || "").replace(/^https?:\/\//, "")
);

async function open() {
  if (!import.meta.client) return;
  isOpen.value = true;
  await nextTick();
  const el = overlayEl.value;
  if (el && typeof el.requestFullscreen === "function") {
    try {
      await el.requestFullscreen();
    } catch (err) {
      // Fullscreen denied (e.g. permissions policy); keep the in-page overlay open as fallback.
      console.warn("QR fullscreen request rejected", err);
    }
  }
}

async function close() {
  if (!import.meta.client) {
    isOpen.value = false;
    return;
  }
  if (document.fullscreenElement) {
    try {
      await document.exitFullscreen();
    } catch (err) {
      console.warn("QR exit fullscreen failed", err);
    }
  }
  isOpen.value = false;
}

function handleFullscreenChange() {
  if (!document.fullscreenElement) {
    isOpen.value = false;
  }
}

function handleBackdropClick(event) {
  if (event.target === event.currentTarget) {
    close();
  }
}

onMounted(() => {
  if (import.meta.client) {
    document.addEventListener("fullscreenchange", handleFullscreenChange);
  }
});

onUnmounted(() => {
  if (import.meta.client) {
    document.removeEventListener("fullscreenchange", handleFullscreenChange);
  }
});

defineExpose({ open, close });
</script>

<template>
  <ClientOnly>
    <Teleport to="body">
      <div
        v-if="isOpen"
        ref="overlayEl"
        class="qr-fullscreen fixed inset-0 z-[60] flex flex-col items-center justify-center gap-6 bg-jv-canvas px-6 py-8 sm:gap-10 sm:px-10 sm:py-12"
        role="dialog"
        aria-modal="true"
        aria-label="Scan QR code"
        @click="handleBackdropClick"
      >
        <button
          type="button"
          aria-label="Close fullscreen"
          class="absolute right-4 top-4 grid size-12 -rotate-3 place-items-center rounded-[10px] border-[3px] border-jv-ink bg-jv-coral text-white shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:right-6 sm:top-6 sm:size-14"
          @click.stop="close"
        >
          <X class="size-6 sm:size-7" :stroke-width="2.6" />
        </button>

        <div
          class="flex max-h-full max-w-full flex-col items-center gap-4 sm:gap-6"
          @click.stop
        >
          <p
            class="font-body text-[12px] font-black uppercase tracking-[0.2em] text-jv-muted sm:text-[14px]"
          >
            Scan to join
          </p>
          <div
            class="qr-card-fullscreen grid place-items-center bg-jv-white p-4 shadow-brutal-lg jv-border-rough sm:p-6"
          >
            <QrCode
              :scan-u-r-l="joinUrl"
              :quiz-code="code"
              :size="720"
              class="size-full"
            />
          </div>
          <div
            class="flex flex-wrap items-center justify-center gap-3 text-center sm:gap-5"
          >
            <div
              class="flex flex-col items-center gap-1 jv-border-rough bg-jv-white px-4 py-2 shadow-brutal-sm sm:px-6 sm:py-3"
            >
              <span
                class="font-body text-[10px] font-black uppercase tracking-[0.14em] text-jv-muted sm:text-[12px]"
              >
                Quiz code
              </span>
              <span
                class="qr-code-text font-feature text-[28px] font-black leading-none text-jv-coral sm:text-[36px]"
              >
                {{ code }}
              </span>
            </div>
            <div
              class="flex flex-col items-center gap-1 jv-border-rough bg-jv-white px-4 py-2 shadow-brutal-sm sm:px-6 sm:py-3"
            >
              <span
                class="font-body text-[10px] font-black uppercase tracking-[0.14em] text-jv-muted sm:text-[12px]"
              >
                Join at
              </span>
              <span
                class="break-all font-body text-[18px] font-extrabold text-jv-ink sm:text-[22px]"
              >
                {{ displayUrl }}
              </span>
            </div>
          </div>
          <p
            class="mt-2 font-body text-[12px] font-bold text-jv-muted sm:text-[14px]"
          >
            Click anywhere outside or press Esc to exit
          </p>
        </div>
      </div>
    </Teleport>
  </ClientOnly>
</template>

<style scoped>
.qr-code-text {
  letter-spacing: 0.4rem;
}

@media (max-width: 768px) {
  .qr-code-text {
    letter-spacing: 0.22rem;
  }
}

.qr-card-fullscreen {
  width: min(70vh, 70vw);
  height: min(70vh, 70vw);
  max-width: 720px;
  max-height: 720px;
}

.qr-card-fullscreen :deep(svg) {
  width: 100%;
  height: 100%;
}

.qr-card-fullscreen :deep(.flex) {
  width: 100%;
  height: 100%;
}
</style>

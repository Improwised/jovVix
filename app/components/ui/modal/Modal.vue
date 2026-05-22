<script setup lang="ts">
import { onBeforeUnmount, watch } from "vue";
import { X } from "lucide-vue-next";
import type { HTMLAttributes } from "vue";
import { cn } from "@/lib/utils";

interface Props {
  modelValue: boolean;
  title?: string;
  description?: string;
  size?: "sm" | "md" | "lg" | "xl";
  closeOnBackdrop?: boolean;
  hideClose?: boolean;
  class?: HTMLAttributes["class"];
}

const props = withDefaults(defineProps<Props>(), {
  title: "",
  description: "",
  size: "md",
  closeOnBackdrop: true,
  hideClose: false,
  class: undefined,
});

const emit = defineEmits<{ "update:modelValue": [value: boolean] }>();

const widthClass = {
  sm: "max-w-[420px]",
  md: "max-w-[520px]",
  lg: "max-w-[720px]",
  xl: "max-w-[900px]",
};

function close() {
  emit("update:modelValue", false);
}

function onBackdropClick() {
  if (props.closeOnBackdrop) close();
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === "Escape") close();
}

watch(
  () => props.modelValue,
  (open) => {
    if (typeof document === "undefined") return;
    if (open) {
      document.body.style.overflow = "hidden";
      document.addEventListener("keydown", onKeydown);
    } else {
      document.body.style.overflow = "";
      document.removeEventListener("keydown", onKeydown);
    }
  }
);

onBeforeUnmount(() => {
  if (typeof document === "undefined") return;
  document.body.style.overflow = "";
  document.removeEventListener("keydown", onKeydown);
});
</script>

<template>
  <Teleport to="body">
    <Transition name="jv-modal">
      <div
        v-if="props.modelValue"
        class="fixed inset-0 z-[1000] flex items-center justify-center px-4 py-6"
        role="dialog"
        aria-modal="true"
      >
        <button
          type="button"
          class="absolute inset-0 cursor-default bg-jv-ink/60"
          aria-label="Close"
          @click="onBackdropClick"
        ></button>

        <div
          :class="
            cn(
              'jv-modal-card relative z-10 w-full -rotate-[0.4deg] jv-border-rough bg-jv-white p-5 shadow-brutal-lg sm:p-7',
              widthClass[props.size],
              props.class
            )
          "
        >
          <span
            class="absolute left-1/2 top-[-12px] z-10 h-4 w-16 -translate-x-1/2 rotate-[1deg] bg-jv-coral"
            aria-hidden="true"
          ></span>

          <button
            v-if="!props.hideClose"
            type="button"
            class="absolute right-3 top-3 grid size-9 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:rotate-[2deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:right-4 sm:top-4"
            aria-label="Close"
            @click="close"
          >
            <X class="size-4" :stroke-width="2.6" />
          </button>

          <header v-if="props.title || $slots.header" class="pr-10 sm:pr-12">
            <slot name="header">
              <h2
                class="font-headings text-[22px] leading-tight text-jv-ink sm:text-[28px]"
              >
                {{ props.title }}
              </h2>
              <p
                v-if="props.description"
                class="mt-2 font-body text-[14px] font-bold text-jv-muted sm:text-[15px]"
              >
                {{ props.description }}
              </p>
            </slot>
          </header>

          <div class="mt-4 sm:mt-5">
            <slot />
          </div>

          <footer
            v-if="$slots.footer"
            class="mt-5 flex flex-col-reverse items-stretch justify-end gap-3 sm:mt-6 sm:flex-row sm:items-center"
          >
            <slot name="footer" />
          </footer>
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

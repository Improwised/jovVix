<script setup>
import { Trash2 } from "lucide-vue-next";
import { Modal } from "@/components/ui/modal";

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: "Delete Item",
  },
  message: {
    type: String,
    default:
      "Once deleted, this item cannot be recovered. Do you want to proceed?",
  },
  confirmLabel: {
    type: String,
    default: "Delete",
  },
  cancelLabel: {
    type: String,
    default: "Cancel",
  },
});

const emit = defineEmits(["update:modelValue", "confirmDelete"]);

function close() {
  emit("update:modelValue", false);
}

function confirm() {
  emit("confirmDelete");
  close();
}
</script>

<template>
  <Modal
    :model-value="props.modelValue"
    :title="props.title"
    :description="props.message"
    size="sm"
    @update:model-value="emit('update:modelValue', $event)"
  >
    <template #footer>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center rounded-full border-[2px] border-jv-ink bg-jv-white px-5 font-body text-[14px] font-black text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-6 sm:text-[15px]"
        @click="close"
      >
        {{ props.cancelLabel }}
      </button>
      <button
        type="button"
        class="inline-flex h-10 items-center justify-center gap-2 rounded-full border-[2px] border-jv-ink bg-jv-coral px-5 font-body text-[14px] font-black text-white shadow-brutal-sm transition-transform hover:rotate-[1deg] active:translate-x-[2px] active:translate-y-[2px] active:shadow-none sm:h-11 sm:px-6 sm:text-[15px]"
        @click="confirm"
      >
        <Trash2 class="size-4" :stroke-width="2.4" />
        {{ props.confirmLabel }}
      </button>
    </template>
  </Modal>
</template>

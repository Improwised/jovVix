<script setup>
import { Pencil, Trash2 } from "lucide-vue-next";

const props = defineProps({
  user: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
});

const emits = defineEmits(["showEditForm", "deleteUserPermission"]);

// open the edit form to change this user's permission for the quiz
const showEditForm = () => {
  emits(
    "showEditForm",
    props.user.id,
    props.user.shared_to,
    props.user.permission
  );
};
</script>

<template>
  <div
    class="flex items-center justify-between gap-3 jv-border-rough border-[2px] border-jv-ink bg-jv-white px-3 py-3 shadow-brutal-sm"
  >
    <!-- User identity -->
    <div class="flex min-w-0 flex-1 items-center gap-3">
      <span class="relative shrink-0">
        <img
          src="https://api.dicebear.com/9.x/bottts/svg?seed=Jade"
          alt="User avatar"
          class="size-10 rounded-full border-[2px] border-jv-ink bg-jv-canvas"
          width="40"
          height="40"
        />
        <span
          class="absolute bottom-0 right-0 size-3 rounded-full border-2 border-jv-white bg-jv-accent-green"
          aria-hidden="true"
        ></span>
      </span>

      <div class="min-w-0 flex-1">
        <h4
          v-if="props.user.first_name?.Valid"
          class="truncate text-[15px] font-black text-jv-ink"
        >
          {{ props.user.first_name.String }}
          {{ props.user.last_name?.String }}
        </h4>
        <h4 v-else class="text-[15px] font-black text-jv-ink">Unknown</h4>
        <div class="truncate text-[13px] font-semibold text-jv-muted">
          {{ props.user.shared_to }}
        </div>
      </div>
    </div>

    <!-- Permission + actions -->
    <div class="flex shrink-0 items-center gap-2">
      <span
        class="rounded-full border-[2px] border-jv-ink bg-jv-yellow px-2.5 py-1 text-[12px] font-black uppercase tracking-[0.08em] text-jv-ink"
      >
        {{ props.user.permission }}
      </span>
      <button
        type="button"
        title="Edit Permission"
        aria-label="Edit permission"
        class="grid size-8 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-white text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-[3deg] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none"
        @click="showEditForm"
      >
        <Pencil class="size-4" :stroke-width="2.4" />
      </button>
      <button
        type="button"
        title="Delete Permission"
        aria-label="Delete permission"
        class="grid size-8 place-items-center rounded-[8px] border-[2px] border-jv-ink bg-jv-coral text-white shadow-brutal-sm transition-transform hover:rotate-[3deg] active:translate-x-[1px] active:translate-y-[1px] active:shadow-none"
        @click="emits('deleteUserPermission', props.user.id)"
      >
        <Trash2 class="size-4" :stroke-width="2.4" />
      </button>
    </div>
  </div>
</template>

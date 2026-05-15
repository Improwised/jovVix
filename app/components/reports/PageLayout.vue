<template>
  <div class="border-b-2 border-jv-ink/15">
    <div class="flex items-center gap-8">
      <button
        type="button"
        class="relative -mb-[2px] pb-3 text-[18px] sm:text-[20px] font-bold transition-colors"
        :class="
          props.currentTab === 'report'
            ? 'text-jv-ink'
            : 'text-jv-muted hover:text-jv-ink'
        "
        @click="changeComponent('report')"
      >
        Questions
        <span
          v-if="props.currentTab === 'report'"
          class="absolute inset-x-0 -bottom-[2px] h-1 bg-jv-yellow"
          aria-hidden="true"
        ></span>
      </button>
      <button
        type="button"
        class="relative -mb-[2px] pb-3 text-[18px] sm:text-[20px] font-bold transition-colors"
        :class="
          props.currentTab === 'participants'
            ? 'text-jv-ink'
            : 'text-jv-muted hover:text-jv-ink'
        "
        @click="changeComponent('participants')"
      >
        Participants
        <span
          v-if="props.currentTab === 'participants'"
          class="absolute inset-x-0 -bottom-[2px] h-1 bg-jv-yellow"
          aria-hidden="true"
        ></span>
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  currentTab: {
    default: "report",
    type: String,
    required: true,
  },
});

const emits = defineEmits(["changeTab"]);
const changeComponent = (tab) => {
  emits("changeTab", tab);
};
</script>

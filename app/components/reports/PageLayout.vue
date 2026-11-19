<template>
  <div class="border-b-2 border-jv-ink/15">
    <div class="flex items-center gap-8">
      <button
        type="button"
        class="relative -mb-[2px] pb-3 text-[18px] sm:text-[20px] font-bold transition-colors"
        :class="
          props.currentTab === props.firstValue
            ? 'text-jv-ink'
            : 'text-jv-muted hover:text-jv-ink'
        "
        @click="changeComponent(props.firstValue)"
      >
        {{ props.firstLabel }}
        <span
          v-if="props.currentTab === props.firstValue"
          class="absolute inset-x-0 -bottom-[2px] h-1 bg-jv-yellow"
          aria-hidden="true"
        ></span>
      </button>
      <button
        type="button"
        class="relative -mb-[2px] pb-3 text-[18px] sm:text-[20px] font-bold transition-colors"
        :class="
          props.currentTab === props.secondValue
            ? 'text-jv-ink'
            : 'text-jv-muted hover:text-jv-ink'
        "
        @click="changeComponent(props.secondValue)"
      >
        {{ props.secondLabel }}
        <span
          v-if="props.currentTab === props.secondValue"
          class="absolute inset-x-0 -bottom-[2px] h-1 bg-jv-yellow"
          aria-hidden="true"
        ></span>
      </button>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  firstValue: { type: String, default: "report" },
  secondValue: { type: String, default: "participants" },
  firstLabel: { type: String, default: "Questions" },
  secondLabel: { type: String, default: "Participants" },
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

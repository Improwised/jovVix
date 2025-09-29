<template>
  <div class="row">
    <div class="col-12">
      <h3 class="mb-2 fw-bold text-center">Quiz Analysis</h3>

      <div
        class="d-flex flex-column flex-md-row align-items-center position-relative"
      >
        <h3>
          {{
            currentTab == "report"
              ? `Total Questions ${props.totalQuestion}/${props.allQuestion}`
              : ""
          }}
        </h3>
        <!-- Dropdown -->
        <div class="ms-md-auto mb-2 mb-md-0 dropdown-container">
          <ReportsDownloadDropdown
            :current-tab="props.currentTab"
            :question-type-filter="props.questionTypeFilter"
            :user-filter="props.userFilter"
            @update:question-type-filter="
              $emit('update:questionTypeFilter', $event)
            "
            @update:user-filter="$emit('update:userFilter', $event)"
          />
        </div>

        <!-- Tabs -->
        <ul class="nav nav-tabs tabs-center">
          <li class="nav-item" @click="changeComponent('report')">
            <NuxtLink
              :class="{ active: props.currentTab === 'report' }"
              class="nav-link"
              :to="`/admin/reports/${activeQuizId}`"
            >
              Questions
            </NuxtLink>
          </li>
          <li class="nav-item" @click="changeComponent('participants')">
            <NuxtLink
              :class="{ active: props.currentTab === 'participants' }"
              class="nav-link"
              :to="`/admin/reports/${activeQuizId}`"
            >
              Participants
            </NuxtLink>
          </li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const activeQuizId = computed(() => route.params.id);

const props = defineProps({
  currentTab: {
    default: "report",
    type: String,
    required: true,
  },
  questionTypeFilter: {
    default: "all",
    type: String,
  },
  userFilter: {
    type: Object,
    default: () => ({
      isAsc: true,
      showTop10: false,
    }),
  },

  totalQuestion: {
    default: 0,
    type: Number,
  },
  allQuestion: {
    default: 0,
    type: Number,
  },
});

const emits = defineEmits([
  "changeTab",
  "update:questionTypeFilter",
  "update:userFilter",
]);
const changeComponent = (tab) => {
  emits("changeTab", tab);
};
</script>

<style scoped>
/* Make tabs truly centered on large screens */
@media (min-width: 830px) {
  .tabs-center {
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    width: auto; /* shrink to fit */
  }
}

/* Ensure dropdown is always clickable above tabs */
.dropdown-container {
  position: relative;
  z-index: 10;
}
</style>

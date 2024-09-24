<template>
  <div class="row">
    <div class="col-lg-12 col-md-12 col-12">
      <h3 class="mb-2 fw-bold">Quiz Analysis</h3>
      <ul class="nav nav-tabs">
        <li class="nav-item" @click="changeComponent('report')">
          <NuxtLink
            :class="{ active: props.currentTab === `report` }"
            class="nav-link"
            :to="`/admin/reports/${activeQuizId}`"
            >Questions</NuxtLink
          >
        </li>
        <li class="nav-item" @click="changeComponent('participants')">
          <NuxtLink
            :class="{ active: props.currentTab === `participants` }"
            class="nav-link"
            :to="`/admin/reports/${activeQuizId}`"
            >Participants</NuxtLink
          >
        </li>
      </ul>
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
});

const emits = defineEmits(["changeTab"]);
const changeComponent = (tab) => {
  emits("changeTab", tab);
};
</script>

<template>
  <div class="row">
    <div class="col-lg-12 col-md-12 col-12">
      <div class="sub-title d-flex justify-content-around">
        <h3 class="mb-2 fw-bold text-center fs-1">Quiz Analysis</h3>
          <button
          class="btn bg-light-primary btn-light btn-link mx-2"
          @click="downloadQuizReport"
          >
          Download report
        </button>
      </div>
      <ul class="nav nav-tabs justify-content-center">
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
const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

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

const downloadQuizReport = async () => {  
  try {
    const res = await fetch(
      `${apiUrl}/analytics_board/admin/${activeQuizId.value}/download`,
      {
        method: 'GET',
        headers: headers,
        mode: "cors",
        credentials: "include",
      }
    );

    if (!res.ok) {
      alert('Network response was not ok');
      return; 
      }

    const blob = await res.blob();

    const url = window.URL.createObjectURL(blob);

    const a = document.createElement('a');
    a.href = url;
    a.download = 'quiz_report.pdf';
    document.body.appendChild(a);
    a.click();
    a.remove();

    window.URL.revokeObjectURL(url);

} catch (err) {
  console.error('Download failed:', err);
  alert('Sorry, the file is not available right now.');
}

}

</script>

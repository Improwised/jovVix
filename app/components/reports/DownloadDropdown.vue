<template>
  <div class="d-flex gap-3 mb-2">
    <!-- If current tab = report -->
    <template v-if="currentTab === 'report'">
      <div class="dropdown">
        <select
          id="dropdown"
          class="form-select"
          :value="props.questionTypeFilter || 'all'"
          @change="emit('update:questionTypeFilter', $event.target.value)"
        >
          <option value="all">All</option>
          <option value="1">Single</option>
          <option value="2">Survey</option>
        </select>
      </div>
    </template>

    <!-- If current tab = participants -->
    <template v-else-if="currentTab === 'participants'">
      <div class="d-flex gap-3">
        <div class="dropdown">
          <select
            id="dropdown"
            v-model="userFilterType"
            class="form-select user-ranking-select"
            @change="handleUserFilterChange"
          >
            <option value="allUsers">All Users</option>
            <option value="top10">Top 10</option>
          </select>
        </div>

        <div v-if="userFilterType === 'allUsers'" class="dropdown">
          <select
            id="dropdown"
            v-model="orderType"
            class="form-select user-ranking-select"
            @change="handleOrderChange"
          >
            <option value="scoreAsc">Score ASC</option>
            <option value="scoreDesc">Score DESC</option>
          </select>
        </div>
      </div>
    </template>

    <!-- Common Download dropdown -->
    <div class="dropdown">
      <button
        id="dropdown"
        class="form-select d-flex align-items-center justify-content-between test-button"
        type="button"
        data-bs-toggle="dropdown"
        aria-expanded="false"
      >
        Download as
      </button>
      <ul class="dropdown-menu" aria-labelledby="downloadDropdown">
        <li>
          <button class="dropdown-item" @click="downloadReport('pdf')">
            PDF
          </button>
        </li>
        <li>
          <button class="dropdown-item" @click="downloadReport('csv')">
            CSV
          </button>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  questionTypeFilter: {
    type: String,
    default: "all",
  },
  currentTab: {
    type: String,
    default: "report",
  },
  userFilter: {
    type: Object,
    default: () => ({
      isAsc: true,
      showTop10: false,
    }),
  },

  isTop10: {
    type: Boolean,
    default: false,
  },
});

const emit = defineEmits(["update:questionTypeFilter", "update:userFilter"]);

const userFilterType = ref("allUsers");
const orderType = ref("scoreDesc");

function handleUserFilterChange() {
  if (userFilterType.value === "allUsers") {
    emit("update:userFilter", {
      ...props.userFilter,
      showTop10: false,
      // keep current order
      isAsc: orderType.value === "scoreAsc",
    });
  } else if (userFilterType.value === "top10") {
    emit("update:userFilter", {
      ...props.userFilter,
      showTop10: true,
    });
  }
}

function handleOrderChange() {
  emit("update:userFilter", {
    ...props.userFilter,
    isAsc: orderType.value === "scoreAsc",
    showTop10: false,
  });
}

const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const activeQuizId = computed(() => route.params.id);

const downloadReport = async (typeOfFile) => {
  if (typeOfFile === "pdf") {
    alert("Working on it.");
    return;
  }

  try {
    let api = "";

    if (props.currentTab === "report") {
      let questionType = "all";
      if (props.questionTypeFilter === "1") questionType = "single";
      if (props.questionTypeFilter === "2") questionType = "survey";

      api = `${apiUrl}/admin/reports/${activeQuizId.value}/download/analysis?question_type=${questionType}`;
    } else if (props.currentTab === "participants") {
      api = `${apiUrl}/admin/reports/${activeQuizId.value}/download/participants`;

      if (props.userFilter.showTop10) {
        api += `?top_10=true`;
      } else {
        api += `?top_10=false&order_by=${
          props.userFilter.isAsc ? "asc" : "desc"
        }`;
      }
    }

    if (!api) return;

    const response = await $fetch.raw(api, {
      method: "GET",
      headers: {
        ...headers,
        "Content-Type": typeOfFile === "csv" ? "text/csv" : "application/pdf",
      },
      credentials: "include",
      responseType: "blob",
    });

    const blob = new Blob([response._data]);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute(
      "download",
      `${
        props.currentTab == "report" ? "quiz" : props.currentTab
      }_report.${typeOfFile}`
    );
    document.body.appendChild(link);
    link.click();
    link.remove();
  } catch (err) {
    console.error("Download failed:", err);
  }
};
</script>

<style scoped>
#dropdown {
  background-color: var(--bs-light-primary);
  color: #212529;
  border: none;
  border-radius: 0.5rem;
  padding: 0.5rem 2rem 0.5rem 1rem;
  font-weight: 500;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml;charset=UTF-8,%3csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 16 16'%3e%3cpath fill='black' d='M1.5 5.5l6 6 6-6'/%3e%3c/svg%3e");
  background-repeat: no-repeat;
  background-position: right 0.75rem center;
  background-size: 1rem;
}
ul {
  padding: 0;
  border: 1px solid black;
  border-radius: 0;
}
li > button {
  background-color: var(--bs-light-primary);
}
li > button:hover {
  background-color: #182965;
  color: aliceblue;
}
</style>

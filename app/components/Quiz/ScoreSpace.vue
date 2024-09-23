<script setup>
import AnswerSubmissionChart from "../AnswerSubmissionChart.vue";
const url = useRuntimeConfig().public;
const props = defineProps({
  data: {
    default: () => {
      return {};
    },
    type: Object,
    required: true,
  },
  isAdmin: {
    default: false,
    type: Boolean,
    required: false,
  },
  userName: {
    required: false,
    type: String,
    default: "",
  },
  answer: {
    required: false,
    type: Number,
  },
  analysisTab: {
    type: String,
  },
});

const emits = defineEmits(["askSkipTimer", "changeAnalysisTab"]);
const timer = ref(null);
const time = ref(0);

function handleTimer() {
  clearInterval(timer.value);
  timer.value = setInterval(() => {
    time.value += 0.1;
    if (time.value == props.data.data.duration + 1) {
      clearInterval(timer.value);
      time.value = -1;
      timer.value = null;
    }
  }, 100);
}

handleTimer();

function handleSkipTimer(e) {
  e.preventDefault();
  emits("askSkipTimer");
}

const changeAnalysisTab = (tab) => emits("changeAnalysisTab", tab);
</script>

<template>
  <Frame page-title="Score Page" page-message="Rank Board">
    <v-progress-linear
      :striped="true"
      color="blue"
      :height="10"
      rounded="true"
      :model-value="(time * 100) / props.data.data.duration"
    ></v-progress-linear>
    <div class="card border-secondary mt-3">
      <div class="card-body">
        <h4 class="card-title">{{ props.data.data.question }}</h4>
        <img
          v-if="props.data?.data?.question_media === 'image'"
          :src="`${props.data?.data?.resource}`"
          :alt="`${props.data?.data?.resource}`"
          class="rounded img-thumbnail"
        />
        <CodeBlockComponent v-if="props.data?.data?.question_media === 'code'" :code="props.data?.data?.resource" />
        <div class="d-flex flex-column">
          <div
            v-for="(answer, key) in props.data.data.options"
            :key="key"
            class="border m-1 rounded p-1"
            :class="{
              'bg-success': answer.isAnswer,
              'bg-danger text-white': !answer.isAnswer && key == props.answer,
            }"
          >
            <img
              v-if="props.data?.data?.options_media === 'image'"
              :src="`${answer.value}`"
              :alt="`${answer.value}`"
              class="rounded img-thumbnail"
            />
            <CodeBlockComponent v-if="props.data?.data?.options_media === 'code'" :code="answer.value" />
            <div
              v-if="props.data?.data?.options_media === 'text'"
              class="form-check form-check-inline"
            >
              <label class="form-check-label">{{ answer.value }}</label>
            </div>
          </div>
        </div>
      </div>
    </div>
    <button
      v-if="isAdmin"
      type="button"
      class="btn text-white btn-primary mt-3"
      @click="handleSkipTimer"
    >
      Skip
    </button>
    <ul
      v-if="isAdmin"
      id="pills-tab"
      class="nav nav-tabs mt-3 mb-3 nav-justified"
      role="tablist"
    >
      <!-- Ranking tab -->
      <li class="nav-item">
        <a
          id="pills-ranking-tab"
          class="nav-link"
          :class="{ active: props.analysisTab === 'ranking' }"
          data-bs-toggle="pill"
          href="#pills-ranking"
          role="tab"
          aria-controls="pills-ranking"
          aria-selected="true"
          @click="changeAnalysisTab('ranking')"
          >Rankings</a
        >
      </li>

      <!-- Charts tab -->
      <li class="nav-item">
        <a
          id="pills-chart-tab"
          class="nav-link"
          :class="{ active: props.analysisTab === 'chart' }"
          data-bs-toggle=""
          href="#pills-chart"
          role="tab"
          aria-controls="pills-chart"
          aria-selected="true"
          @click="changeAnalysisTab('chart')"
          >Chart</a
        >
      </li>
    </ul>
    <!-- Ranking Table -->
    <div id="pills-tabContent" class="tab-content p-4">
      <div
        id="pills-ranking"
        class="tab-pane fade"
        :class="{ 'active show': props.analysisTab === 'ranking' || !isAdmin }"
        role="tabpanel"
        aria-labelledby="pills-ranking-tab"
      >
        <div class="table-responsive mt-2">
          <table class="table table-borderless align-middle">
            <thead class="table-light"></thead>
            <tbody class="table-group-divider">
              <tr
                v-for="(user, index) in props.data.data.rankList"
                :key="index"
              >
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  {{ user.rank }}
                </td>
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  {{ user.firstname }}
                  <span v-if="user.username === props.userName"
                    >&nbsp; ({{ user.username }})</span
                  >
                </td>
                <td :class="{ 'bg-primary': user.username === props.userName }">
                  {{ user.score }}
                </td>
              </tr>
            </tbody>
            <tfoot></tfoot>
          </table>
        </div>
      </div>
      <!-- Chart -->
      <div
        id="pills-chart"
        class="tab-pane fade"
        :class="{ 'active show': props.analysisTab === 'chart' }"
        role="tabpanel"
        aria-labelledby="pills-chart-tab"
      >
        <AnswerSubmissionChart
          v-if="isAdmin"
          :options="props.data.data.options"
          :responses="props.data.data.userResponses"
        />
      </div>
    </div>
  </Frame>
</template>

<style scoped>
.table-custom {
  background-color: #f9f9f9; /* Example light color */
}

.table td,
.table th {
  background-color: #f9f9f9; /* Ensure each cell has a white background */
}
</style>

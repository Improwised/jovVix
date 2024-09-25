<script setup>
const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
  userName: {
    type: String,
    required: true,
  },
});

const analysis = ref({});
const incorrectWidth = ref(0);
const unattemptedWidth = ref(0);

onMounted(() => {
  analysis.value = questionsAnalysis(props.data);
  unattemptedWidth.value =
    (analysis.value?.unAttemptedQuestions / analysis.value?.totalQuestions) *
    100;
  incorrectWidth.value =
    100 - analysis.value?.accuracy - unattemptedWidth.value;
});

const handleMouseEnter = (event) => {
  event.target.style.transform = "scale(1.05)"; // Scale up on hover
  event.target.style.transition = "transform 0.3s ease"; // Add transition
};

const handleMouseLeave = (event) => {
  event.target.style.transform = "scale(1)"; // Reset scale on leave
};
</script>

<template>
  <!-- userQuestionsAnalysis Modal -->

  <div
    :id="props.userName"
    class="modal fade"
    tabindex="-1"
    aria-labelledby="exampleModalLabel"
    aria-hidden="true"
  >
    <div
      class="modal-dialog modal-lg modal-dialog-scrollable modal-dialog-centered"
    >
      <div class="modal-content">
        <div class="modal-header">
          <h1 id="exampleModalLabel" class="modal-title fs-5">
            Questions Analysis
          </h1>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="modal"
            aria-label="Close"
          ></button>
        </div>
        <div class="modal-body">
          <QuizAnalysis :data="props.data" />
        </div>
      </div>
    </div>
  </div>

  <!-- Participants Analysis-->
  <div
    class="user-analytics-item"
    type="button"
    data-bs-toggle="modal"
    :data-bs-target="`#${props.userName}`"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
    <div class="user-stats-box">
      <div class="header">
        <div class="avatar-container">
          <img
            class="avatar"
            src="../../assets/images/avatar.png"
            alt="Avatar"
          />
          <div class="name">
            {{ props.data[0].firstname }}
            <span>({{ props.data[0].username }})</span>
          </div>
        </div>
        <div class="stats">
          <div class="stat-item">
            <span class="value">{{ analysis?.rank }}</span>
            <span class="label">Rank</span>
          </div>
          <div class="stat-item">
            <span class="value">{{ analysis?.accuracy }}%</span>
            <span class="label">Accuracy</span>
          </div>
          <div class="stat-item">
            <span class="value">{{ analysis?.totalScore }}</span>
            <span class="label">Score</span>
          </div>
        </div>
      </div>
      <div class="quiz-header mb-4">
        <div class="quiz-accuracy position-relative w-100">
          <div class="progress">
            <div
              class="progress-bar bg-success"
              role="progressbar"
              :style="{ width: analysis?.accuracy + '%' }"
              aria-valuenow="70"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
            <div
              class="progress-bar bg-danger"
              role="progressbar"
              :style="{ width: incorrectWidth + '%' }"
              aria-valuenow="20"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
            <div
              class="progress-bar bg-secondary"
              role="progressbar"
              :style="{ width: unattemptedWidth + '%' }"
              aria-valuenow="10"
              aria-valuemin="0"
              aria-valuemax="100"
            ></div>
          </div>
        </div>
      </div>
      <div class="d-flex justify-content-center">
        &#9989; {{ analysis?.correctAnwers }} &ensp; &#10060;
        {{ analysis?.wrongAnwers }} &ensp; &#x25CC;
        {{ analysis?.unAttemptedQuestions }} &ensp;
        <span v-if="analysis?.totalSurveyQuestions > 0"
          >&#128203; {{ analysis?.attemptedSurveyQuestions }} /
          {{ analysis?.totalSurveyQuestions }}</span
        >
      </div>
    </div>
  </div>
</template>

<style scoped>
.user-stats-box {
  border: 1px solid #ddd;
  padding: 10px;
  border-radius: 8px;
  width: 100%;
  max-width: 600px;
  background-color: white;
  margin: 0 auto;
  box-sizing: border-box;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
  /* Box shadow added */
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.avatar-container {
  display: flex;
  align-items: center;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
}

.name {
  font-size: 18px;
  font-weight: bold;
  margin-left: 10px;
}

.stats {
  display: flex;
  flex-wrap: nowrap;
  justify-content: flex-end;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-left: 20px;
}

.label {
  font-size: 12px;
  color: #888;
}

.value {
  font-size: 14px;
  font-weight: bold;
}

.divider {
  width: 100%;
  border: none;
  height: 1px;
  background-color: #000;
  margin: 10px 0;
}

.progress-bar {
  display: flex;
  height: 10px;
  border-radius: 5px;
  overflow: hidden;
  width: 100%;
}

.correct {
  background-color: #4caf50;
}

.incorrect {
  background-color: #f44336;
}

@media (max-width: 600px) {
  .header {
    flex-wrap: wrap;
    justify-content: center;
  }

  .avatar-container {
    margin-bottom: 10px;
  }

  .stats {
    justify-content: center;
    flex-wrap: nowrap;
  }

  .stat-item {
    margin-left: 10px;
    margin-top: 0;
  }
}

.user-analytics-item {
  padding: 10px;
  margin-bottom: 10px;
  border: none;
  /* Remove border */
  border-radius: 5px;
  cursor: pointer;
  transition: transform 0.3s ease;
  /* Add transition for scale */
}

.user-analytics-item:hover {
  transform: scale(1.05);
  /* Scale up on hover */
}
</style>

<script setup>
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
});
const emits = defineEmits(["askSkipTimer"]);
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
</script>

<template>
  <Frame page-title="Score page" page-message="rank-board">
    <v-progress-linear
      :striped="true"
      color="blue"
      :height="10"
      rounded="true"
      :model-value="(time * 100) / props.data.data.duration"
    ></v-progress-linear>
    <div class="card border-secondary">
      <div class="card-body">
        <h4 class="card-title">{{ props.data.data.question }}</h4>
        <div class="d-flex">
          <div
            v-for="(answer, key) in props.data.data.options"
            :key="key"
            class="flex-grow-1 border m-1 rounded p-1 card-text"
            :class="{ 'bg-success': answer.isAnswer }"
          >
            <div class="form-check form-check-inline">
              <label class="form-check-label">{{ answer.value }}</label>
            </div>
          </div>
        </div>
      </div>
    </div>

    <button
      v-if="isAdmin"
      type="button"
      class="btn btn-primary mt-3"
      @click="handleSkipTimer"
    >
      Next Question
    </button>

    <div class="table-responsive mt-5">
      <table
        class="table table-striped table-hover table-borderless table-light align-middle"
      >
        <thead class="table-light">
          <caption>
            Rankings
          </caption>
          <tr>
            <th>Rank</th>
            <th>User</th>
            <th>Score</th>
            <th>Response Time (ms)</th>
          </tr>
        </thead>
        <tbody class="table-group-divider">
          <tr
            v-for="(user, index) in props.data.data.rankList"
            :key="index"
            class="table-light"
          >
            <td scope="row">{{ user.rank }}</td>
            <td>{{ user.username }}</td>
            <td>{{ user.score }}</td>
            <td>{{ user.response_time }}</td>
          </tr>
        </tbody>
        <tfoot></tfoot>
      </table>
    </div>
  </Frame>
</template>

<script setup>
const headers = useRequestHeaders(["cookie"]);
const url = useRuntimeConfig().public;
const route = useRoute();

const played_quiz_id = computed(() => route.params.played_quiz_id);

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(`${url.api_url}/user_played_quizes/${played_quiz_id.value}`, {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});

const userAccuracy = computed(() => {
  if (quizList.value?.data) {
    let userAnswerAnalysis = [];
    quizList.value.data.filter((item) => {
      if (item.question_type != "survey" && item.is_attend) {
        //check if the answer is correct or not
        let correctIncorrectFlag = isCorrectAnswer(
          item.selected_answer.String,
          item.correct_answer
        );
        userAnswerAnalysis.push(correctIncorrectFlag);
      }
    });
    return (
      (userAnswerAnalysis.filter(Boolean).length / userAnswerAnalysis.length) *
      100
    ).toFixed(2);
  }

  return 0;
});

const userTotalScore = computed(() => {
  if (quizList.value?.data) {
    let userTotalScore = 0;
    quizList.value.data.map((item) => {
      userTotalScore += item.calculated_score;
    });
    return userTotalScore;
  }

  return 0;
});
</script>

<template>
  <ClientOnly>
    <div>
      <div v-if="quizError" class="alert alert-danger" role="alert">
        {{ quizError.data }}
      </div>
      <div v-else-if="quizPending" class="form-select w-full md:w-20rem">
        Pending...
      </div>
      <div>
        <h3 class="text-center">Accuracy: {{ userAccuracy }}%</h3>
        <h3 class="text-center">Total Score: {{ userTotalScore }}</h3>
        <Frame
          v-for="(item, index) in quizList?.data"
          :key="index"
          :page-title="'Q' + (index + 1) + '. ' + item.question"
        >
          <ul style="list-style-type: none; padding-left: 0">
            <li
              v-for="(option, key) in item.options"
              :key="key"
              style="display: flex; align-items: center; padding-left: 20px"
            >
              <span
                v-if="item.correct_answer.includes(key)"
                style="margin-right: 10px"
                >&#10004;</span
              >
              <span
                v-if="
                  item.selected_answer.String.includes(key) &&
                  !item.correct_answer.includes(key)
                "
                style="margin-right: 10px"
                >&#10006;</span
              >
              <span>{{ key }}: {{ option }}</span>
            </li>
          </ul>
          <div
            style="
              display: flex;
              flex: 1;
              margin-top: 10px;
              border-top: 1px solid #ccc;
            "
          >
            <div
              v-if="item.response_time > 0"
              style="flex: 1; padding: 10px; border-right: 1px solid #ccc"
            >
              Response Time:
              {{ (item.response_time / 1000).toFixed(2) }} seconds
            </div>
            <div
              v-else
              style="flex: 1; padding: 10px; border-right: 1px solid #ccc"
            >
              Response Time: -
            </div>
            <div style="flex: 1; padding: 10px">
              {{ item.is_attend ? "Attempted" : "Not Attempted" }}
            </div>
          </div>
        </Frame>
      </div>
    </div>
  </ClientOnly>
</template>

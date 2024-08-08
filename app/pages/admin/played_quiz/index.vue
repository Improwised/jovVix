<template>
  <div class="container">
    <h3 class="text-center">Played Quiz List</h3>
    <UtilsQuizListWaiting v-if="quizPending" />

    <div v-else-if="quizError" class="alert alert-danger" role="alert">
      {{ quizError.data }}
    </div>
    <div v-else class="row">
      <div
        v-for="(details, index) in quizList.data"
        :key="index"
        class="card-body col-md-3 mt-3"
      >
        <PlayedQuizListCard :details="details" />
      </div>
      <!-- list loader -->

      <!-- quiz details -->
      <!-- show quiz list -->
      <!-- <div v-else>
                <div class="card text-center">
                    <div v-if="quizList.data == null || quizList.data.length < 1"
                        class="no-quiz-list d-flex flex-column align-items-center mt-4 mb-2">
                        <h2>No Quiz Played By You !</h2>
                    </div>
                    <div v-else class="row">
                        <div v-for="(details, index) in quizList.data" :key="index" class="card-body col-md-4">
                            <PlayedQuizListCard :details="details" />
                        </div>
                    </div>
                </div>
            </div> -->
    </div>
  </div>
</template>

<script setup>
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(url.api_url + "/user_played_quizes", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
});
</script>

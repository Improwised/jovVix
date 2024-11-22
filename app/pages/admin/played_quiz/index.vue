<template>
  <div class="container p-0">
    <div class="d-flex flex-column justify-content-center">
      <UtilsQuizListWaiting v-if="quizPending" />
      <div v-else-if="quizError">{{ quizError.message }}</div>
      <div v-else>
        <div
          v-if="quizList?.data.length < 1"
          class="no-quiz-list d-flex flex-column align-items-center"
        >
          <h1>No Quiz Played By You ! !</h1>
        </div>
        <div v-else>
          <nav class="navbar pb-4">
            <div class="container-fluid p-0">
              <h1 class="mb-0">Played Quiz List</h1>
              <input
                v-model="titleInput"
                type="text"
                placeholder="Search quiz"
                class="border rounded p-2"
              />
            </div>
          </nav>
          <div class="d-flex flex-column gap-3">
            <div class="row">
              <div
                v-for="(details, index) in quizList?.data?.data"
                :key="index"
                class="col-md-6 mb-5"
              >
                <QuizListCard :details="details" :is-played-quiz="true" />
              </div>
            </div>
          </div>
          <div class="d-flex align-items-center justify-content-center">
            <Pagination
              :page="page"
              :num-of-records="quizList?.data?.count / 10"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import debounce from "lodash/debounce";
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const page = computed(() => Number(route.query.page) || 1);
const titleInput = ref(route.query.title);
const title = computed(() => route.query.title || "");

const {
  data: quizList,
  pending: quizPending,
  error: quizError,
} = useFetch(url.apiUrl + "/user_played_quizes", {
  method: "GET",
  headers: headers,
  mode: "cors",
  credentials: "include",
  query: {
    page,
    title,
  },
});

const debouncedNavigateTo = debounce((query) => {
  navigateTo({
    path: route.path,
    query: query,
  });
}, 500);

watch(titleInput, (newValue) => {
  debouncedNavigateTo({
    ...route.query,
    title: newValue,
  });
});
</script>

<style scoped>
.max-width {
  max-width: 922px;
}
</style>

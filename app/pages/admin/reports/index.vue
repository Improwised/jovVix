<template>
  <div v-if="quizListPending">Loading....</div>
  <div v-else-if="quizListError" class="text-danger alert alert-danger mt-3">
    Error while fetching data:
    <span>
      {{ quizListError }}
    </span>
  </div>
  <div v-else class="row mt-3">
    <div class="card">
      <div class="card-header bg-white py-4 row">
        <div class="col-sm-12 col-md-4">
          <h3>Quiz Reports</h3>
        </div>
        <div
          class="col-sm-12 col-md-8 d-flex align-items-center justify-content-md-end flex-wrap gap-2"
        >
          <input
            v-model="nameInput"
            type="text"
            placeholder="Search quiz"
            class="border rounded p-2 mx-2"
          />
          <input
            v-model="date"
            type="datetime-local"
            placeholder="Select date"
            class="border rounded p-2 mx-2"
          />
        </div>
      </div>
      <div class="table-responsive">
        <table class="table text-nowrap mb-0">
          <thead class="table-light">
            <tr>
              <th role="button" @click="sortEventHandler('title')">
                Quiz Name
                <font-awesome-icon
                  v-if="orderBy === 'title' && order === 'asc'"
                  icon="arrow-up-short-wide"
                  class="bx bx-sort-up"
                />
                <font-awesome-icon
                  v-else-if="orderBy === 'title' && order === 'desc'"
                  icon="arrow-up-wide-short"
                  class="bx bx-sort-down"
                />
                <font-awesome-icon v-else icon="sort" class="bx bx-sort" />
              </th>
              <th role="button" @click="sortEventHandler('participants')">
                Total Participants
                <font-awesome-icon
                  v-if="orderBy === 'participants' && order === 'asc'"
                  icon="arrow-up-short-wide"
                  class="bx bx-sort-up"
                />
                <font-awesome-icon
                  v-else-if="orderBy === 'participants' && order === 'desc'"
                  icon="arrow-up-wide-short"
                  class="bx bx-sort-down"
                />
                <font-awesome-icon v-else icon="sort" class="bx bx-sort" />
              </th>
              <th role="button" @click="sortEventHandler('activated_from')">
                Starts At
                <font-awesome-icon
                  v-if="orderBy === 'activated_from' && order === 'asc'"
                  icon="arrow-up-short-wide"
                  class="bx bx-sort-up"
                />
                <font-awesome-icon
                  v-else-if="orderBy === 'activated_from' && order === 'desc'"
                  icon="arrow-up-wide-short"
                  class="bx bx-sort-down"
                />
                <font-awesome-icon v-else icon="sort" class="bx bx-sort" />
              </th>
              <th role="button" @click="sortEventHandler('activated_to')">
                Ends At
                <font-awesome-icon
                  v-if="orderBy === 'activated_to' && order === 'asc'"
                  icon="arrow-up-short-wide"
                  class="bx bx-sort-up"
                />
                <font-awesome-icon
                  v-else-if="orderBy === 'activated_to' && order === 'desc'"
                  icon="arrow-up-wide-short"
                  class="bx bx-sort-down"
                />
                <font-awesome-icon
                  v-else
                  role="button"
                  icon="sort"
                  class="bx bx-sort"
                />
              </th>
              <th role="button" @click="sortEventHandler('questions')">
                Total Questions
                <font-awesome-icon
                  v-if="orderBy === 'questions' && order === 'asc'"
                  icon="arrow-up-short-wide"
                  class="bx bx-sort-up"
                />
                <font-awesome-icon
                  v-else-if="orderBy === 'questions' && order === 'desc'"
                  icon="arrow-up-wide-short"
                  class="bx bx-sort-down"
                />
                <font-awesome-icon v-else icon="sort" class="bx bx-sort" />
              </th>
              <th>Accuracy</th>
            </tr>
          </thead>
          <tbody>
            <div v-if="quizList?.data?.Count <= 0">No quiz found....</div>
            <tr
              v-for="(quiz, index) in quizList?.data.Data"
              v-else
              :key="index"
              role="button"
              @click="viewReport(quiz.id)"
            >
              <td>
                <div class="ms-3 lh-1">
                  <p class="mb-1 h4 font-weight-bold text-primary">
                    {{ decodeURI(quiz.title) }}
                  </p>
                  <p class="text-secondary">{{ quiz.description.String }}</p>
                </div>
              </td>
              <td>{{ quiz.participants }}</td>
              <td>{{ quiz.activated_from.Time }}</td>
              <td>{{ quiz.activated_to.Time }}</td>
              <td>{{ quiz.questions }}</td>
              <td>
                {{
                  (
                    (quiz.correct_answers /
                      (quiz.participants * quiz.questions)) *
                    100
                  ).toFixed(2)
                }}%
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="d-flex align-items-center justify-content-center">
        <Pagination :page="page" :num-of-records="quizList?.data?.Count / 10" />
      </div>
    </div>
  </div>
</template>
<script setup>
import { format, formatISO } from "date-fns";
import debounce from "lodash/debounce";

definePageMeta({
  layout: "default",
});

const { apiUrl } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const route = useRoute();
const page = computed(() => Number(route.query.page) || 1);
const order = computed(() => route.query.order || "desc");
const orderBy = computed(() => route.query.orderBy || "activated_to");
const date = ref(route.query.date);
const name = computed(() => route.query.name || "");
const nameInput = ref(route.query.name);

const {
  data: quizList,
  error: quizListError,
  pending: quizListPending,
} = useFetch(`${apiUrl}/admin/reports/list`, {
  transform: (quizList) => {
    quizList.data.Data = quizList.data?.Data?.map((quiz) => {
      quiz.activated_from.Time = format(
        formatISO(quiz.activated_from.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      quiz.activated_to.Time = format(
        formatISO(quiz.activated_to.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      return quiz;
    });
    return quizList;
  },
  credentials: "include",
  headers: headers,
  query: {
    page,
    orderBy,
    order,
    name,
    date,
  },
});

const viewReport = (activeQuizId) => {
  navigateTo(`/admin/reports/${activeQuizId}`);
};

const sortEventHandler = (columnName) => {
  let ordercol = order.value;
  if (columnName === orderBy.value) {
    ordercol === "asc" ? (ordercol = "desc") : (ordercol = "asc");
  }
  navigateTo({
    path: route.path,
    query: {
      ...route.query,
      orderBy: columnName,
      order: ordercol,
    },
  });
};

const debouncedNavigateTo = debounce((query) => {
  navigateTo({
    path: route.path,
    query: query,
  });
}, 500);

watch(nameInput, (newName) => {
  debouncedNavigateTo({
    ...route.query,
    name: newName,
  });
});
</script>

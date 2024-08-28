<template>
  <div v-if="quizListPending">loading....</div>
  <div v-else-if="quizListError">{{ quizListError }}</div>
  <div v-else class="row mt-6">
    <div class="card">
      <div class="card-header bg-white py-4 row">
        <h4 class="mb-0 col-sm-12 col-md-4">Quiz Reports</h4>

        <input
          type="text"
          v-model="nameInput"
          class="col-sm-5 col-md-3 border rounded p-2 mx-2"
          placeholder="Search quiz"
        />
        <input
          type="datetime-local"
          v-model="date"
          placeholder="Select date"
          class="col-sm-5 col-md-3 border rounded p-2 mx-2"
        />
      </div>
      <div class="table-responsive">
        <table class="table text-nowrap mb-0">
          <thead class="table-light">
            <tr>
              <th>
                Quiz Name
                <font-awesome-icon
                  icon="arrow-up-short-wide"
                  v-if="orderBy === 'title' && order === 'asc'"
                  class="bx bx-sort-up"
                  @click="sortEventHandler('title')"
                />
                <font-awesome-icon
                  icon="arrow-up-wide-short"
                  v-else-if="orderBy === 'title' && order === 'desc'"
                  class="bx bx-sort-down"
                  @click="sortEventHandler('title')"
                />
                <font-awesome-icon
                  icon="sort"
                  v-else
                  class="bx bx-sort"
                  @click="sortEventHandler('title')"
                />
              </th>
              <th>
                Total Participants
                <font-awesome-icon
                  icon="arrow-up-short-wide"
                  v-if="orderBy === 'participants' && order === 'asc'"
                  class="bx bx-sort-up"
                  @click="sortEventHandler('participants')"
                />
                <font-awesome-icon
                  icon="arrow-up-wide-short"
                  v-else-if="orderBy === 'participants' && order === 'desc'"
                  class="bx bx-sort-down"
                  @click="sortEventHandler('participants')"
                />
                <font-awesome-icon
                  icon="sort"
                  v-else
                  class="bx bx-sort"
                  @click="sortEventHandler('participants')"
                />
              </th>
              <th>
                Starts At
                <font-awesome-icon
                  icon="arrow-up-short-wide"
                  v-if="orderBy === 'activated_from' && order === 'asc'"
                  class="bx bx-sort-up"
                  @click="sortEventHandler('activated_from')"
                />
                <font-awesome-icon
                  icon="arrow-up-wide-short"
                  v-else-if="orderBy === 'activated_from' && order === 'desc'"
                  class="bx bx-sort-down"
                  @click="sortEventHandler('activated_from')"
                />
                <font-awesome-icon
                  icon="sort"
                  v-else
                  class="bx bx-sort"
                  @click="sortEventHandler('activated_from')"
                />
              </th>
              <th>
                Ends At
                <font-awesome-icon
                  icon="arrow-up-short-wide"
                  v-if="orderBy === 'activated_to' && order === 'asc'"
                  class="bx bx-sort-up"
                  @click="sortEventHandler('activated_to')"
                />
                <font-awesome-icon
                  icon="arrow-up-wide-short"
                  v-else-if="orderBy === 'activated_to' && order === 'desc'"
                  class="bx bx-sort-down"
                  @click="sortEventHandler('activated_to')"
                />
                <font-awesome-icon
                  icon="sort"
                  v-else
                  class="bx bx-sort"
                  @click="sortEventHandler('activated_to')"
                />
              </th>
              <th>
                Total Questions
                <font-awesome-icon
                  icon="arrow-up-short-wide"
                  v-if="orderBy === 'questions' && order === 'asc'"
                  class="bx bx-sort-up"
                  @click="sortEventHandler('questions')"
                />
                <font-awesome-icon
                  icon="arrow-up-wide-short"
                  v-else-if="orderBy === 'questions' && order === 'desc'"
                  class="bx bx-sort-down"
                  @click="sortEventHandler('questions')"
                />
                <font-awesome-icon
                  icon="sort"
                  v-else
                  class="bx bx-sort"
                  @click="sortEventHandler('questions')"
                />
              </th>
              <th>Accuracy</th>
            </tr>
          </thead>
          <tbody>
            <div v-if="quizList?.data?.Count <= 0">No quiz found....</div>
            <tr
              v-else
              v-for="(quiz, index) in quizList?.data.Data"
              :key="index"
            >
              <td>
                <div class="ms-3 lh-1">
                  <h5 class="mb-1">
                    <NuxtLink
                      :to="`/admin/reports/${quiz.id}`"
                      class="text-inherit"
                      >{{ decodeURI(quiz.title) }}
                    </NuxtLink>
                  </h5>
                  <p>{{ quiz.description.String }}</p>
                </div>
              </td>
              <td>{{ quiz.participants }}</td>
              <td>{{ quiz.activated_from.Time }}</td>
              <td>{{ quiz.activated_to.Time }}</td>
              <td>{{ quiz.questions }}</td>
              <td>{{ (quiz.correct_answers / (quiz.participants * quiz.questions) * 100).toFixed(2) }}%</td>
            </tr>
          </tbody>
        </table>
      </div>
      <Pagination :page="page" :numOfRecords="quizList?.data?.Count" />
    </div>
  </div>
</template>
<script setup>
import { format, formatISO } from "date-fns";
import debounce from "lodash/debounce";

definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
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
} = useFetch(`${api_url}/admin/reports/list`, {
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

<template>
  <div v-if="quizListPending">loading....</div>
  <div v-else-if="quizListError">{{ quizListError }}</div>
  <!-- <div v-else>{{ quizList || "lol" }}</div> -->
  <div v-else class="row mt-6">
    <div class="card">
      <div class="card-header bg-white py-4">
        <h4 class="mb-0">Quiz Reports</h4>
      </div>
      <div class="table-responsive">
        <<table class="table text-nowrap mb-0 table-hover">
          <thead class="table-light">
            <tr>
              <th>Quiz Name</th>
              <th>Total Participants</th>
              <th>Starts At</th>
              <th>Ends At</th>
              <th>Total Questions</th>
              <th>Accuracy</th>
            </tr>
          </thead>
          <tbody>
            <tr
              v-for="(quiz, index) in paginatedQuizList"
              :key="index"
              class="clickable-row"
              @click="navigateToQuiz(quiz.id)"
            >
              <td>
                <div class="ms-3 lh-1">
                  <h5 class="mb-1">
                    <NuxtLink
                      :to="`/admin/reports/${quiz.id}`"
                      class="text-inherit"
                      >{{ decodeURI(quiz.title) }}</NuxtLink
                    >
                  </h5>
                  <p>{{ quiz.description.String }}</p>
                </div>
              </td>
              <td>{{ quiz.participants }}</td>
              <td>{{ quiz.activated_from.Time }}</td>
              <td>{{ quiz.activated_to.Time }}</td>
              <td>{{ quiz.questions }}</td>
              <td>
                {{
                  (quiz.correct_answers /
                    (quiz.participants * quiz.questions)) *
                  100
                }}%
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <!-- Pagination Controls -->
      <div class="mt-3">
        <button
          :disabled="currentPage === 1"
          @click="currentPage--"
          class="btn btn-primary me-2"
        >
          Previous
        </button>
        <button
          :disabled="currentPage >= totalPages"
          @click="currentPage++"
          class="btn btn-primary"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { format, formatISO } from "date-fns";
import { useRouter } from 'vue-router';


definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const router = useRouter();


const itemsPerPage = 8;
const currentPage = ref(1);



const {
  data: quizList,
  error: quizListError,
  pending: quizListPending,
} = useFetch(`${api_url}/admin/reports/list`, {
  transform: (quizList) => {
    return quizList.data
      .map((quiz) => {
        quiz.activated_from.Time = format(
          formatISO(quiz.activated_from.Time),
          "dd-MMM-yyyy HH:mm:ss"
        );
        quiz.activated_to.Time = format(
          formatISO(quiz.activated_to.Time),
          "dd-MMM-yyyy HH:mm:ss"
        );
        return quiz;
      })
      .sort((a, b) => new Date(b.activated_from.Time) - new Date(a.activated_from.Time)); // Sorting by activated_from.Time in descending order
  },
  credentials: "include",
  headers: headers,
});

// Computed property to get the paginated list
const paginatedQuizList = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage;
  const end = start + itemsPerPage;
  return quizList.value.slice(start, end);
});

// Computed property to get the total number of pages
const totalPages = computed(() => {
  return Math.ceil(quizList.value.length / itemsPerPage);
});

// Function to navigate to the quiz details page
const navigateToQuiz = (quizId) => {
  router.push(`/admin/reports/${quizId}`);
}
</script>

<style scoped>
/* Applying light purple border on hover for table rows */
.table-hover tbody tr {
  transition: border 0.2s ease-in-out;
}

.table-hover tbody tr:hover {
  border: 2px solid #3d3d7e; /* Light purple border */
  cursor: pointer;
}
</style>
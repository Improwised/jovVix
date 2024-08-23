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
        <table class="table text-nowrap mb-0">
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
            <div v-if="!quizList?.length">No quiz found....</div>
            <tr v-for="(quiz, index) in quizList" :key="index">
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
    </div>
  </div>
</template>
<script setup>
import { format, formatISO } from "date-fns";

definePageMeta({
  layout: "default",
});

const { api_url } = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);

const {
  data: quizList,
  error: quizListError,
  pending: quizListPending,
} = useFetch(`${api_url}/admin/reports/list`, {
  transform: (quizList) =>
    quizList.data?.map((quiz) => {
      quiz.activated_from.Time = format(
        formatISO(quiz.activated_from.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      quiz.activated_to.Time = format(
        formatISO(quiz.activated_to.Time),
        "dd-MMM-yyyy HH:mm:ss"
      );
      return quiz;
    }),
  credentials: "include",
  headers: headers,
});
</script>

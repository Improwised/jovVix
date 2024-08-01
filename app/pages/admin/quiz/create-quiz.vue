<script setup>
// core dependencies
import { useToast } from "vue-toastification";

// define nuxt configs
const toast = useToast();
const app = useNuxtApp();
useSystemEnv();

// define props and emits
let file = ref(0);
let title = ref("");
const url = useRuntimeConfig().public;
let quizId = ref();

// core
async function submit(e) {
  e.preventDefault();
  const formData = new FormData();

  const description = document.getElementById("description");
  const attachment = document.getElementById("attachment");

  formData.append(description.name, description.value);
  formData.append(attachment.name, attachment.files[0]);
  const { data, error } = await useFetch(
    encodeURI(url.api_url + "/admin/quizzes/" + title.value + "/upload"),
    {
      method: "POST",
      body: formData,
      mode: "cors",
      credentials: "include",
    }
  );

  if (error.value?.data) {
    toast.error(error.value.data.data);
    return;
  }

  quizId.value = data.value.data;
  toast.success(app.$CsvUploadSuccess);
}
</script>

<template>
  <Frame
    page-title="Create Quiz"
    page-message="Create New Quiz By Uploading CSV"
  >
    <form @submit="submit">
      <div class="mb-3">
        <div class="mb-3">
          <label for="title" class="form-label">Quiz Title</label>
          <input
            id="title"
            v-model="title"
            type="text"
            class="form-control"
            name="title"
            aria-describedby="helpId"
            placeholder=""
            required
          />
          <small v-if="title == ''" id="helpId" class="form-text text-danger"
            >Required</small
          >
        </div>
        <div class="mb-3">
          <label for="description" class="form-label">Quiz Description</label>
          <input
            id="description"
            type="text"
            class="form-control"
            name="description"
            aria-describedby="helpId"
            placeholder=""
          />
          <!-- <small id="helpId" class="form-text text-muted">Help text</small> -->
        </div>

        <label for="attachment" class="form-label">Choose File</label>
        <input
          id="attachment"
          type="file"
          class="form-control"
          name="attachment"
          placeholder="upload"
          aria-describedby="fileHelpId"
          accept=".csv"
          @change="(e) => (file = e.target.files.length)"
        />
        <div v-if="file == 0" id="fileHelpId" class="form-text text-danger">
          Required
        </div>
      </div>
      <button type="submit" class="btn text-white btn-primary me-2">
        Create Quiz
      </button>
      <a class="btn btn-primary me-2" href="/files/demo.csv" download="demo.csv"
        >Download Sample</a
      >
      <UtilsStartQuiz v-if="quizId" :quiz-id="quizId" />
    </form>
  </Frame>
</template>

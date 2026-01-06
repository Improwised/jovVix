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
const requestPending = ref(false);
const imageRequestPending = ref(false);
const requiredImage = ref([]);

const uploadQuizAndQuestions = async (e) => {
  e.preventDefault();
  requestPending.value = true;
  const formData = new FormData();

  const description = document.getElementById("description");
  const attachment = document.getElementById("attachment");

  formData.append(description.name, description.value);
  formData.append(attachment.name, attachment.files[0]);
  try {
    await $fetch(encodeURI(`${url.apiUrl}/quizzes/${title.value}/upload`), {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: formData,
      mode: "cors",
      credentials: "include",
        onResponse({ response }) {
          if (response.status === 202) {
            quizId.value = response._data?.data;
            toast.success(app.$CsvUploadSuccess);
          }
        },
    });
  } catch (error) {
      const parsed = error?.data?.data ?? error?.data?.message ?? error?.message ?? JSON.stringify(error)
      toast.error(parsed || "error while creating quiz")
      requestPending.value = false
      return
  }

  try {
    await $fetch(
      encodeURI(`${url.apiUrl}/quizzes/${quizId.value}/questions?media=image`),
      {
        method: "GET",
        headers: {
          Accept: "application/json",
        },
        mode: "cors",
        credentials: "include",
        onResponse({ response }) {
          if (response.status != 200) {
            requestPending.value = false;
            toast.error("error while create quiz");
            return;
          }
          if (response.status == 200) {
            const data = response._data.data?.data;
            if (data) {
              requiredImage.value = data;
            }
            requestPending.value = false;
          }
        },
      }
    );
  } catch (error) {
    toast.error(error.message);
    requestPending.value = false;
    return;
  }
};

const imageFileUpload = async (e) => {
  if (e.target.files.length === 0) {
    toast.error("Please select a file to upload.");
    return;
  }

  // Validate file
  const file = e.target.files[0];

  const validImageTypes = [
    "image/jpeg",
    "image/png",
    "image/gif",
    "image/webp",
    "image/heic",
    "image/heif",
  ];
  if (!validImageTypes.includes(file.type)) {
    toast.error(
      "Please upload a valid image file (JPEG, PNG, GIF, WEBP, HEIC, HEIF)."
    );
    return;
  }

  // 1 MB max
  if (file.size > 1000000) {
    toast.error("Please upload an image less than 1 MB.");
    return;
  }

  const imageForm = new FormData();
  imageForm.append("image-attachment", e.target.files[0], e.target.name);

  imageRequestPending.value = true;
  try {
    await $fetch(encodeURI(`${url.apiUrl}/images?quiz_id=${quizId.value}`), {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: imageForm,
      mode: "cors",
      credentials: "include",
      onResponse({ response }) {
        if (response.status != 200) {
          imageRequestPending.value = false;
          toast.error("error upload image");
          return;
        }
        if (response.status == 200) {
          toast.success(response._data?.data);
          imageRequestPending.value = false;
        }
      },
    });
  } catch (error) {
    toast.error(error.message);
    imageRequestPending.value = false;
    return;
  }
};
</script>

<template>
  <!-- ImageUpload Modal -->
  <div v-if="requiredImage.length > 0">
    <div id="imageUpload" class="modal fade" tabindex="-1" aria-labelledby="exampleModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-xl modal-dialog-scrollable modal-dialog-centered">
        <div class="modal-content">
          <div class="modal-header">
            <h1 id="exampleModalLabel" class="modal-title fs-5">
              Questions Analysis
            </h1>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <table class="table table-responsive">
              <thead>
                <tr>
                  <th scope="col">Question</th>
                  <th scope="col">QuestionImage</th>
                  <th scope="col">OptionsImage</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="value in requiredImage" :key="value.question_id">
                  <td>{{ value.question }}</td>
                  <td>
                    <v-file-input v-if="value.question_media == 'image'" id="image-attachment-question"
                      prepend-icon="mdi-camera" type="file" class="form-control" :name="value.question_id"
                      label="Question" accept="image/*" @change="imageFileUpload">
                    </v-file-input>
                  </td>
                  <td v-if="value.options_media == 'image'">
                    <v-file-input v-for="index in 5" id="image-attachment-option" :key="index"
                      :name="index + '_' + value.question_id" :label="'Option ' + index" prepend-icon="mdi-camera"
                      type="file" class="form-control mb-2" accept="image/*" @change="imageFileUpload">
                    </v-file-input>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
          <div class="modal-footer text-white">
            <button v-if="imageRequestPending" type="button" class="btn btn-secondary">
              Pending...
            </button>
            <button v-else type="button" class="btn btn-secondary" data-bs-dismiss="modal">
              Close
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Frame page-title="Create Quiz" page-message="Create New Quiz By Uploading CSV">
    <form @submit="uploadQuizAndQuestions">
      <div class="mb-3">
        <div class="mb-3">
          <label for="title" class="form-label">Quiz Title
            <small v-if="title == ''" id="helpId" class="form-text text-danger">*</small>
          </label>
          <input id="title" v-model="title" type="text" class="form-control" name="title" aria-describedby="helpId"
            placeholder="" required />

        </div>
        <div class="mb-3">
          <label for="description" class="form-label">Quiz Description</label>
          <input id="description" type="text" class="form-control" name="description" aria-describedby="helpId"
            placeholder="" required />
          <!-- <small id="helpId" class="form-text text-muted">Help text</small> -->
        </div>
        <div class="mb-3">
          <label for="attachment" class="form-label">Choose File
            <small v-if="file == 0" id="fileHelpId" class="form-text text-danger">*</small>
          </label>
          <input id="attachment" type="file" class="form-control" name="attachment" placeholder="upload"
            aria-describedby="fileHelpId" accept=".csv" @change="(e) => (file = e.target.files.length)" />

        </div>
      </div>
      <div class="d-flex p-2">
        <button v-if="requestPending" class="btn text-white btn-primary me-2">
          Pending...
        </button>
        <button v-else type="submit" class="btn text-white btn-primary me-2">
          Create Quiz
        </button>
        <a class="btn btn-primary me-2" href="/files/demo.csv" download="demo.csv">Download Sample</a>
        <div v-if="requiredImage.length > 0" data-bs-toggle="modal" :data-bs-target="`#imageUpload`" type="button"
          class="btn text-white btn-primary me-2">
          Upload Images
        </div>
        <UtilsStartQuiz v-if="quizId && !requestPending" :quiz-id="quizId" />
      </div>
    </form>
  </Frame>
</template>

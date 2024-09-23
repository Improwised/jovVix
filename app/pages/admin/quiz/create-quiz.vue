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
const requiredImage = ref([])
const imageForm = new FormData();

const uploadQuizAndQuestions = async (e) => {
  e.preventDefault();
  requestPending.value = true;
  const formData = new FormData();

  const description = document.getElementById("description");
  const attachment = document.getElementById("attachment");

  formData.append(description.name, description.value);
  formData.append(attachment.name, attachment.files[0]);
  try {
    await $fetch(
      encodeURI(`${url.api_url}/admin/quizzes/${title.value}/upload`),
      {
        method: "POST",
        headers: {
          Accept: "application/json",
        },
        body: formData,
        mode: "cors",
        credentials: "include",
        onResponse({ response }) {
          if (response.status != 202) {
            requestPending.value = false;
            toast.error("error while create quiz");
            return;
          }
          if (response.status == 202) {
            quizId.value = response._data?.data;
            toast.success(app.$CsvUploadSuccess);
          }
        },
      }
    );
  } catch (error) {
    toast.error(error.message);
    requestPending.value = false;
    return;
  }

  try {
      await $fetch(encodeURI(`${url.api_url}/admin/quizzes/${quizId.value}?media=image`), {
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
            requiredImage.value = response._data.data;
            requestPending.value = false;
          }
        },
      });
    } catch (error) {
      toast.error(error.message);
      requestPending.value = false;
      return;
    }
};

const imageFileAppend = (e) => {
  imageForm.append("image-attachment", e.target.files[0], e.target.name);
};

const imageUlpoads3 = async() => {
  const imageAttachment = document.getElementById("image-attachment");
  if (imageAttachment.files.length !== 0) {
    imageRequestPending.value = true;
    try {
      await $fetch(encodeURI(`${url.api_url}/images?quiz_id=${quizId.value}`), {
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
            toast.error("error while create quiz");
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
  }
};
</script>

<template>

    <!-- ImageUpload Modal -->
     <div v-if="requiredImage.length > 0">
      <div
        id="imageUpload"
        class="modal fade"
        tabindex="-1"
        aria-labelledby="exampleModalLabel"
        aria-hidden="true"
      >
        <div
          class="modal-dialog modal-xl modal-dialog-scrollable modal-dialog-centered"
        >
          <div class="modal-content">
            <div class="modal-header">
              <h1 id="exampleModalLabel" class="modal-title fs-5">
                Questions Analysis
              </h1>
              <button
                type="button"
                class="btn-close"
                data-bs-dismiss="modal"
                aria-label="Close"
              ></button>
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
                    <tr v-for="(value) in requiredImage" :key="value.id">
                      <td>{{ value.question }}</td>
                      <td>
                        <v-file-input
                          v-if="value.question_media == 'image'"
                          prepend-icon="mdi-camera"
                          id="image-attachment"
                          type="file"
                          class="form-control"
                          :name="value.id"
                          label="Question"
                          accept="image/*"
                          @change="imageFileAppend"
                        >
                        </v-file-input>
                      </td>
                      <td>
                        <v-file-input
                          v-if="value.options_media == 'image'"
                          v-for="index in 5"
                          prepend-icon="mdi-camera"
                          id="image-attachment"
                          type="file"
                          class="form-control mb-2"
                          :name="index + '_' + value.id"
                          :label="'Option ' + index"
                          accept="image/*"
                          @change="imageFileAppend"
                        >
                        </v-file-input>
                      </td>
                    </tr>
                  </tbody>
              </table>
            </div>
            <div class="modal-footer text-white">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
              <button v-if="imageRequestPending" type="button" class="btn btn-primary">Pending...</button>
              <button v-else @click="imageUlpoads3" type="button" class="btn btn-primary">Save changes</button>
            </div>
          </div>
        </div>
      </div>
     </div>
    

  <Frame
    page-title="Create Quiz"
    page-message="Create New Quiz By Uploading CSV"
  >
    <form @submit="uploadQuizAndQuestions">
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
            required
          />
          <!-- <small id="helpId" class="form-text text-muted">Help text</small> -->
        </div>
        <div class="mb-3">
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
      </div>
      <div class="d-flex p-2">
        <button v-if="requestPending" class="btn text-white btn-primary me-2">
          Pending...
        </button>
        <button v-else type="submit" class="btn text-white btn-primary me-2">
          Create Quiz
        </button>
        <a
          class="btn btn-primary me-2"
          href="/files/demo.csv"
          download="demo.csv"
          >Download Sample</a
        >
        <div
          v-if="requiredImage.length > 0"
          data-bs-toggle="modal"
          :data-bs-target="`#imageUpload`" 
          type="button"
          class="btn text-white btn-primary me-2"
        >
          Upload Images
        </div>
        <UtilsStartQuiz v-if="quizId" :quiz-id="quizId" />
      </div>
    </form>
  </Frame>
</template>

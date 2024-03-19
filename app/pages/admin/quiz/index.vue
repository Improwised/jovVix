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
let urls = useState("urls");

// if urls not found
if (!urls.value.api_url) {
  toast.warning("please reload the page");
}

// core
async function submit(e) {
  urls = useState("urls");
  e.preventDefault();
  const formData = new FormData();

  const description = document.getElementById("description");
  const attachment = document.getElementById("attachment");

  formData.append("title", title.value);
  formData.append(description.name, description.value);
  formData.append(attachment.name, attachment.files[0]);

  const { error } = await useFetch(urls.value.api_url + "/admin/uploads", {
    method: "POST",
    body: formData,
    mode: "cors",
    credentials: "include",
  });

  if (error.value) {
    toast.error(error.value);
    return;
  }

  toast.success(app.$CsvUploadSuccess);
}
</script>

<template>
  <Frame
    page-title="Create quiz"
    page-message="create new quiz by uploading csv"
  >
    <form @submit="submit">
      <div class="mb-3">
        <div class="mb-3">
          <label for="title" class="form-label">Quiz title</label>
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
          <label for="description" class="form-label">Quiz description</label>
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

        <label for="attachment" class="form-label">Choose file</label>
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
      <a class="btn btn-primary me-2" href="/files/demo.csv" download="demo.csv"
        >Download sample</a
      >
      <button type="submit" class="btn btn-primary">Submit</button>
    </form>
  </Frame>
</template>

<template>
  <div class="row m-2">
    <div class="col-lg-12">
      <input
        v-model="editableQuestion.question"
        class="form-control font-bold my-2"
        placeholder="Edit Question"
      />
    </div>
    <div v-if="editableQuestion.question_media == 'image'" class="col-lg-12">
      <v-file-input
        id="image-attachment-question"
        prepend-icon="mdi-camera"
        type="file"
        class="col-12"
        :name="editableQuestion.question_id"
        label="Image Upload"
        accept="image/*"
        show-size
        variant="outlined"
        @change="onImageChange"
      >
      </v-file-input>
    </div>
    <div
      v-if="editableQuestion.question_media === 'image'"
      class="d-flex align-items-center justify-content-center"
    >
      <img
        v-if="editableQuestion.resource"
        :src="editableQuestion.resource"
        :alt="editableQuestion.resource"
        class="rounded img-thumbnail"
      />
    </div>
    <CodeBlockComponent
      v-if="editableQuestion.question_media === 'code'"
      :code="editableQuestion.resource"
      :read-only="false"
      @code-change="changeCode"
    />
  </div>
  <div class="row d-flex align-items-stretch m-2">
    <div
      v-for="(option, order) in props.question.options"
      :key="order"
      class="col-lg-6 col-md-12"
    >
      <div class="option-box wrong-option">
        <div class="d-flex mb-2">
          <div
            v-if="props.question.options_media == 'image'"
            class="container flex-column"
          >
            <v-file-input
              id="image-attachment-option"
              :key="index"
              :name="order + '_' + props.question.question_id"
              :label="'Option ' + order"
              prepend-icon="mdi-camera"
              type="file"
              class="form-control mb-2"
              accept="image/*"
              @change="onImageChange"
            >
            </v-file-input>
            <div class="d-flex justify-content-center mb-2">
              <img
                :src="editableOptions[order]"
                :alt="option"
                class="rounded img-thumbnail"
              />
            </div>
          </div>

          <CodeBlockComponent
            v-if="editableQuestion.options_media === 'code'"
            :code="option"
            :read-only="false"
            :option-order="Number(order)"
            @code-change="changeCode"
          />
          <input
            v-if="editableQuestion.options_media === 'text'"
            v-model="editableOptions[order]"
            class="form-control"
            placeholder="Edit Option"
          />
          <input
            v-if="props.question.question_type_id == 1"
            :id="order"
            v-model="picked"
            type="radio"
            :value="order"
            class="ml-1"
          />
          <!-- <input v-else  type="checkbox" :id="order"  :value="order" v-model="checkedNames" class="ml-1"/> -->
        </div>
      </div>
    </div>
  </div>
  <div class="text-right">
    <button class="btn btn-primary" @click="updateQuestion">
      Save Changes
    </button>
  </div>
</template>

<script setup>
import { useToast } from "vue-toastification";
const url = useRuntimeConfig().public;
const headers = useRequestHeaders(["cookie"]);
const toast = useToast();

const props = defineProps({
  question: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
  quizId: {
    type: String,
    required: false,
    default: "",
  },
  questionId: {
    type: String,
    required: false,
    default: "",
  },
});

const editableQuestion = ref({ ...props.question });
const editableOptions = ref({ ...props.question.options });
const correctAnswers = ref(props.question.correct_answer);
const checkedNames = ref(JSON.parse(correctAnswers.value));
const picked = ref(checkedNames.value[0]);

const changeCode = (data, order) => {
  if (order) {
    editableOptions.value[order] = data;
    return;
  }
  editableQuestion.value.resource = data;
};

const onImageChange = async (e) => {
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

  // 2 MB max
  if (file.size > 2000000) {
    toast.error("Please upload an image less than 2 MB.");
    return;
  }

  if (e.target.id === "image-attachment-option" && file) {
    const order = e.target.name[0];
    const reader = new FileReader();
    reader.onload = (e) => {
      editableOptions.value[order] = e.target.result;
    };
    reader.readAsDataURL(file);
  } else if (file) {
    const reader = new FileReader();
    reader.onload = (e) => {
      editableQuestion.value.resource = e.target.result;
    };
    reader.readAsDataURL(file);
  }

  const imageForm = new FormData();
  imageForm.append("image-attachment", e.target.files[0], e.target.name);

  try {
    await $fetch(encodeURI(`${url.api_url}/images?quiz_id=${props.quizId}`), {
      method: "POST",
      headers: {
        Accept: "application/json",
      },
      body: imageForm,
      mode: "cors",
      credentials: "include",
      onResponse({ response }) {
        if (response.status != 200) {
          toast.error("error upload image");
          return;
        }
        if (response.status == 200) {
          toast.success(response._data?.data);
        }
      },
    });
  } catch (error) {
    toast.error(error.message);
    imageRequestPending.value = false;
    return;
  }
};

const updateQuestion = async () => {
  try {
    if (props.question.options_media === "image") {
      Object.entries(editableOptions.value).forEach(([key]) => {
        editableOptions.value[
          key
        ] = `${props.quizId}/${key}_${props.questionId}`;
      });
    }

    const payload = {
      question: editableQuestion.value.question,
      type: props.question.question_type_id,
      options: editableOptions.value,
      answers:
        props.question.question_type_id === 1
          ? [Number(picked.value)]
          : checkedNames.value.map(Number),
      points: props.question.points,
      duration_in_seconds: props.question.duration_in_seconds,
      question_media: props.question.question_media,
      options_media: props.question.options_media,
      resource:
        props.question.question_media === "image"
          ? props.quizId + "/" + props.questionId
          : editableQuestion.value.resource,
    };

    await $fetch(
      `${url.api_url}/quizzes/${props.quizId}/questions/${props.questionId}`,
      {
        method: "PUT",
        headers: headers,
        body: payload,
        credentials: "include",
      }
    );
    toast.success("Question updated successfully!");
  } catch (error) {
    console.error("Failed to update the question", error);
    toast.error("Failed to update the question.");
  }
};
</script>

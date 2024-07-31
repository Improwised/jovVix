<script setup>
import { useToast } from "vue-toastification";
let urls = useRuntimeConfig().public;
const router = useRouter();
const toast = useToast();
const props = defineProps({
  quizId: {
    default: () => {
      return "";
    },
    type: String,
    required: true,
  },
});
console.log();

async function handleStartDemo() {
  try {
    const { data, error } = await useFetch(
      encodeURI(
        urls.api_url + "/admin/quizzes/" + props.quizId + "/demo_session"
      ),
      {
        method: "POST",
        mode: "cors",
        credentials: "include",
      }
    );

    if (error.value?.data) {
      console.log(urls.value.api_url);
      toast.error(error.value.data.data);
      return;
    }

    router.push("/admin/arrange/" + data.value.data);
  } catch (error) {
    console.error(error);
  }
}
</script>
<template>
  <button
    type="button"
    class="btn text-white btn-primary me-0"
    @click="handleStartDemo"
  >
    Start Quiz
  </button>
</template>

<script setup>
let urls = useState("urls");
const router = useRouter();
const props = defineProps({
  quizId: {
    default: () => {
      return "";
    },
    type: String,
    required: true,
  },
});

async function handleStartDemo() {
  const { data, error } = await useFetch(
    encodeURI(
      urls.value.api_url + "/admin/quizzes/" + props.quizId + "/demo_session"
    ),
    {
      method: "POST",
      mode: "cors",
      credentials: "include",
    }
  );

  if (error.value?.data) {
    toast.error(error.value.data.data);
    return;
  }

  router.push("/admin/arrange/" + data.value.data);
}
</script>
<template>
  <button type="button" class="btn btn-primary me-0" @click="handleStartDemo">
    Start Quiz
  </button>
</template>

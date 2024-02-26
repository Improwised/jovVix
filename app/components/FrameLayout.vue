<script setup>
import { useRouter } from "nuxt/app";
import { useToast } from "vue-toastification";

// Define props using defineProps
const props = defineProps({
  pageTitle: {
    type: String,
    required: true,
    default: "title",
  },
  pageWelcomeMessage: {
    type: String,
    required: false,
    default: null,
  },
  maxWidth: {
    type: String,
    required: false,
    default: "700px",
  },
});

const router = useRouter();
const toast = useToast();

const errorQueryParam = router.currentRoute.value.query.error;

if (errorQueryParam) {
  toast.error(errorQueryParam);

  onMounted(() => {
    if (process.client) {
      setTimeout(() => {
        const updatedQuery = { ...router.currentRoute.value.query };
        delete updatedQuery.error;
        router.replace({ query: updatedQuery });
      }, 3000);
    }
  });
}
</script>

<template>
  <div class="d-flex justify-content-center">
    <div
      class="border p-2 m-0 m-sm-5 p-sm-5 rounded"
      :style="{ width: props.maxWidth }"
    >
      <h1>{{ pageTitle }}</h1>
      <h6 v-if="props.pageWelcomeMessage">{{ pageWelcomeMessage }}</h6>

      <hr class="m-2" />
      <slot></slot>
    </div>
  </div>
</template>

<style scoped>
.max-width {
  width: 700px;
}
</style>

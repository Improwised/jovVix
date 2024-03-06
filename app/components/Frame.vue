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
  pageMessage: {
    type: String,
    required: false,
    default: null,
  },
});

const router = useRouter();
const toast = useToast();

onMounted(() => {
  if (process.client) {
    const errorQueryParam = router.currentRoute.value.query?.error;
    if (errorQueryParam) {
      toast.error(errorQueryParam);
    }
  }
});
</script>

<template>
  <div class="d-flex justify-content-center">
    <div class="border p-2 m-0 m-sm-5 p-sm-5 max-width rounded">
      <h1>{{ pageTitle }}</h1>
      <h6 v-if="props.pageMessage">{{ pageMessage }}</h6>

      <hr class="m-2" />
      <slot></slot>
    </div>
  </div>
</template>

<style>
.max-width {
  width: 700px;
}
</style>

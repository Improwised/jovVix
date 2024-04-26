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

const toastError = () => {
  const errorQueryParam = router.currentRoute.value.query?.error;
  if (errorQueryParam) {
    toast.error(errorQueryParam);
  }
};
onMounted(() => {
  if (process.client) {
    toastError();
  }
});

watch(
  () => router.currentRoute.value.query,
  () => toastError()
);
</script>

<template>
  <div class="d-flex justify-content-center container">
    <div class="border p-2 m-0 m-sm-5 p-sm-5 max-width rounded">
      <div class="d-flex flex-row">
        <div class="flex-grow-1">
          <h1>{{ pageTitle }}</h1>
          <h6 v-if="props.pageMessage">{{ pageMessage }}</h6>
        </div>
        <div>
          <slot name="sub-title"></slot>
        </div>
      </div>

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

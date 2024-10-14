<script setup>
import { useRouter } from "nuxt/app";
import { useToast } from "vue-toastification";
import { useMusicStore } from "~~/store/music";
const musicStore = useMusicStore();
const { getMusic, setMusic } = musicStore;

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
  musicComponent: {
    type: Boolean,
    required: false,
    default: false,
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

const music = computed(() => {
  return getMusic();
});
</script>

<template>
  <div class="d-flex justify-content-center container p-0 bg-transparent">
    <div class="border p-2 m-0 m-sm-5 p-sm-5 max-width rounded bg-white shadow">
      <div class="d-flex flex-row">
        <div class="flex-grow-1">
          <h1 class="join-page-title">{{ pageTitle }}</h1>
          <h6 v-if="props.pageMessage">{{ pageMessage }}</h6>
        </div>
        <div>
          <slot name="sub-title"></slot>
        </div>
      </div>

      <hr class="m-2" />
      <slot></slot>
      <div v-if="props.musicComponent" class="d-flex justify-content-end">
        <button v-if="music" class="" @click="setMusic(false)">
          <font-awesome-icon :icon="['fas', 'volume-high']" />
        </button>
        <button v-else class="" @click="setMusic(true)">
          <font-awesome-icon :icon="['fas', 'volume-xmark']" />
        </button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.max-width {
  width: 800px;
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.join-page-title {
  color: #663399;
}

@media (max-width: 576px) {
  .max-width {
    width: 100%;
    padding: 1rem;
    margin: 0.5rem;
    box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
  }
}

.shadow {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}
</style>

<script setup>
// define props and emits
const props = defineProps({
  details: {
    default: () => {
      return [];
    },
    type: Object,
    required: true,
  },
});
</script>
<template>
  <div class="card text-center">
    <div class="card-body">
      <h5 class="card-title">{{ decodeURI(props.details?.title) }}</h5>
      <p class="card-text">{{ props.details?.description }}</p>
      <p class="card-text text-muted">
        {{ props.details?.total_questions }} Questions
      </p>
      <p class="card-text text-muted">
        <small class="text-muted">
          {{ useGetTime(props.details?.created_at) }}
        </small>
      </p>
      <button
        class="btn btn-primary btn-sm"
        @click="navigateTo(`/admin/${props.details?.id}`)"
      >
        Details
      </button>
    </div>
  </div>
</template>
<style scoped>
.randomImg {
  object-fit: cover;
}

img {
  height: auto !important;
  max-height: 120px !important;
}
</style>

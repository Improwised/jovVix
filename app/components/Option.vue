<template>
  <div class="d-flex align-items-center flex-wrap flex-grow-1">
    <button
      class="btn btn-icon btn-white border border-2 rounded-circle btn-dashed ms-2 mr-2"
    >
      {{ String.fromCharCode(64 + Number(props.order)) }}
    </button>
    <div
      v-if="props.optionsMedia === 'image'"
      class="d-flex flex-grow-1 justify-content-center"
    >
      <img
        :src="`${props.option}`"
        :alt="`${props.option}`"
        class="rounded img-thumbnail"
      />
    </div>
    <div
      v-if="props.optionsMedia === 'text'"
      :class="{ 'text-success': props.isCorrect }"
      class="mx-3 font-weight-bold"
    >
      {{ props.option }}
    </div>
    <div v-if="props.optionsMedia === 'code'" class="code-block-container">
      <CodeBlockComponent :code="props?.option" />
    </div>
  </div>
  <span
    v-if="props.isAdminAnalysis"
    :class="{ 'bg-success': props.isCorrect }"
    class="badge bg-secondary text-white mx-3 rounded-pill"
    ><font-awesome-icon icon="fa-solid fa-user" class="mx-2" />
    {{ props.selected }}</span
  >
</template>
<script setup>
const props = defineProps({
  order: {
    type: Number,
    required: false,
    default: 0,
  },
  selected: {
    type: Number,
    required: false,
    default: 0,
  },
  option: {
    type: String,
    required: false,
    default: "",
  },
  isCorrect: {
    type: Boolean,
    required: false,
    default: false,
  },
  optionsMedia: {
    type: String,
    required: true,
    default: "",
  },
  isAdminAnalysis: {
    type: Boolean,
    required: false,
    default: false,
  },
});
</script>

<style scoped>
img {
  height: 150px;
  max-width: 180px;
}
.code-block-container {
  flex: 1;
}
</style>

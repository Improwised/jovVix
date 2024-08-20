<template>
  <div class="d-flex align-items-center flex-wrap">
    <button
      class="btn btn-icon btn-white border border-2 rounded-circle btn-dashed ms-2"
    >
      {{ String.fromCharCode(64 + Number(props.order)) }}
    </button>
    <div class="mx-3 text-bold">{{ props.option }}</div>
  </div>
  <span
    :class="{ 'bg-success': props.isCorrect }"
    class="badge bg-primary text-white mx-3 rounded-pill"
    ><font-awesome-icon icon="fa-solid fa-user" class="mx-2" />
    {{ props.selected }}</span
  >
</template>
<script setup>
const props = defineProps(["order", "selected", "option", "icon", "isCorrect"]);
</script>

<template>
  <div class="d-flex align-items-center flex-wrap">
    <button
      class="btn btn-icon btn-white border border-2 rounded-circle btn-dashed ms-2"
    >
      {{ String.fromCharCode(64 + Number(props.order)) }}
    </button>
    <div v-if="props.optionsMedia === 'image'" class="d-flex align-items-center justify-content-center">
      <img
        :src="`${props.option}`"
        :alt="`${props.option}`"
        class="rounded img-thumbnail"
      />
    </div>
    <div v-if="props.optionsMedia === 'text'" :class="{'text-success': props.isCorrect }" class="mx-3 font-weight-bold">{{ props.option }}</div>
    <div v-if="props.optionsMedia === 'code'" class="d-flex align-items-center justify-content-center" >
      <CodeBlockComponent :code="props?.option"/>
    </div>
  </div>
  <span
    :class="{ 'bg-success': props.isCorrect }"
    class="badge bg-secondary text-white mx-3 rounded-pill"
    ><font-awesome-icon icon="fa-solid fa-user" class="mx-2" />
    {{ props.selected }}</span
  >
</template>
<script setup>
const props = defineProps(["order", "selected", "option", "icon", "isCorrect","optionsMedia"]);
</script>

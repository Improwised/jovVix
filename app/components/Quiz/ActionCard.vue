<template>
  <div class="action-card" @click="$emit('cardClick')">
    <slot name="image" />
    <div>
      <h3 class="fs-3 fw-bold mt-2">{{ props.actionTitle }}</h3>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  actionTitle: {
    default: "",
    type: String,
    required: true,
  },
});

defineEmits(["cardClick"]);
</script>

<style scoped>
.action-card {
  max-width: 300px;
  min-height: 200px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  max-width: 500px;
  height: 300px;
  padding: 35px;
  border: 1px solid rgba(255, 255, 255, 0.25);
  border-radius: 20px;
  box-shadow: 0 0 10px 1px rgba(0, 0, 0, 0.25);
  backdrop-filter: blur(15px);
  cursor: pointer;
  margin: 20px;
}
</style>

<template>
  <div class="full-page-popup">
    <div class="popup-content">
      <!-- Your popup content goes here -->
      <button @click="closePopup">Close Popup</button>
    </div>
  </div>
</template>

<script setup>
import { defineEmits } from "vue";

const emits = defineEmits(["close"]);

const closePopup = () => {
  emits("close");
};
</script>

<style scoped>
.full-page-popup {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5); /* semi-transparent background */
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: center;
}

.popup-content {
  background-color: white;
  padding: 20px;
  border-radius: 10px;
  width: 80%;
  max-width: 600px; /* Adjust width as needed */
  max-height: 80%;
  overflow-y: auto; /* Enable scrolling if content exceeds height */
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.2);
}
</style>

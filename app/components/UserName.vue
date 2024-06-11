<script setup>
const props = defineProps({
  userName: {
    type: String,
    required: true,
  },
});
</script>

<template>
  <div class="container-fluid">
    <div class="row justify-content-center">
      <div class="col-6 col-md-4 mt-5 d-flex justify-content-center">
        <div
          class="border border-1 px-7 py-2 py-md-4 border-radius d-flex align-items-center"
        >
          <img
            src="../assets/images/avatar.png"
            alt="User Avatar"
            height="70"
            width="70"
            class="me-3"
          />
          <h5 class="text-center text-sm fs-5 mb-0">{{ props.userName }}</h5>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.border-radius {
  border-radius: 2rem !important;
}
</style>

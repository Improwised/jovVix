<template>
  <div class="d-flex flex-column align-items-center">
    <qrcode-vue :value="qrUrl" :size="250" level="H" render-as="svg" />
    <div class="d-flex align-items-center justify-content-center mt-1">
      <div>{{ props.scanURL }}</div>
      <font-awesome-icon
        id="URL-input"
        icon="fa-solid fa-copy"
        style="color: #0c6efd"
        class="ml-1"
        role="button"
      />
    </div>
  </div>
</template>
<script setup>
import QrcodeVue from "qrcode.vue";
import usecopyToClipboard from "~~/composables/copy_to_clipboard";

const props = defineProps({
  scanURL: {
    type: String,
    required: true,
    default: "",
  },
  quizCode: {
    type: String,
    required: true,
    default: "",
  },
});

const qrUrl = computed(() =>
  props.quizCode ? `${props.scanURL}?code=${props.quizCode}` : props.scanURL
);

onMounted(() => {
  const copyBtn = document.getElementById("URL-input");
  if (process.client && copyBtn) {
    copyBtn.addEventListener("click", () => {
      usecopyToClipboard(qrUrl.value);
    });
  }
});
</script>

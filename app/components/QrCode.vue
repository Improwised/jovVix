<template>
  <div class="d-flex flex-column align-items-center">
    <qrcode-vue :value="qrUrl" :size="props.size" level="H" render-as="svg" />
  </div>
</template>
<script setup>
import QrcodeVue from "qrcode.vue";

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
  size: {
    type: Number,
    required: false,
    default: 250,
  },
});

const qrUrl = computed(() =>
  props.quizCode ? `${props.scanURL}?code=${props.quizCode}` : props.scanURL
);
</script>

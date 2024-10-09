<template>
  <ClientOnly>
    <CodeEditor
      v-model="code"
      :wrap="true"
      font-size="10px"
      width="98%"
      :read-only="props.readOnly"
      :header="false"
      theme="atom-one-dark"
      :value="props.code"
    ></CodeEditor>
  </ClientOnly>
</template>

<script setup>
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import hljs from "highlight.js";
import CodeEditor from "simple-code-editor";
const emits = defineEmits(["codeChange"]);

const props = defineProps({
  code: {
    default: "",
    type: String,
    required: true,
  },
  readOnly: {
    type: Boolean,
    required: false,
    default: true,
  },
  optionOrder: {
    type: Number,
    required: false,
    default: 0,
  },
});

const code = ref(props.code);
watch(
  code,
  () => {
    emits("codeChange", code.value, props.optionOrder);
  },
  { deep: true }
);
</script>

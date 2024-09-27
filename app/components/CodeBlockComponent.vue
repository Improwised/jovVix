<template>
  <ClientOnly>
    <CodeEditor
      :wrap="true"
      font-size="10px"
      width="98%"
      height="150px"
      :read-only="true"
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

const props = defineProps({
  code: {
    default: "",
    type: String,
    required: true,
  },
});
</script>

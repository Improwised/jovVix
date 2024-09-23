<template>
    <ClientOnly>
      <CodeEditor :wrap="true" font-size="10px" width="100%" height="100%" :read-only="true" :header="false" theme="atom-one-dark" :value="props.code"></CodeEditor>
    </ClientOnly>
  </template>
  
  <script setup>
  import hljs from 'highlight.js';
  import CodeEditor from 'simple-code-editor';
  
  const props = defineProps({
    code: {
      default: "",
      type: String,
      required: true,
    },
  });
  </script>
<template>
  <ErrorTemplate
    v-bind="{
      status,
      statusText,
      statusCode: status,
      statusMessage: statusText,
      description,
      stack,
    }"
  />
</template>

<script setup>
// Mirrors Nuxt 3's built-in error page so 404/500 screens stay unchanged
// after the Nuxt 4 upgrade.
import Error404 from "@/components/error/Error404.vue";
import Error500 from "@/components/error/Error500.vue";

const props = defineProps({
  error: {
    type: Object,
    required: true,
  },
});

const escapeHtml = (text) =>
  text.replace(
    /[&<>"']/g,
    (c) =>
      ({ "&": "&amp;", "<": "&lt;", ">": "&gt;", '"': "&quot;", "'": "&#39;" })[
        c
      ]
  );

const _error = props.error;
const stacktrace =
  import.meta.dev && _error.stack
    ? _error.stack
        .split("\n")
        .splice(1)
        .map((line) => {
          const text = line
            .replace("webpack:/", "")
            .replace(".vue", ".js")
            .trim();
          return {
            text,
            internal:
              (line.includes("node_modules") && !line.includes(".cache")) ||
              line.includes("internal") ||
              line.includes("new Promise"),
          };
        })
        .map(
          (i) =>
            `<span class="stack${i.internal ? " internal" : ""}">${escapeHtml(
              i.text
            )}</span>`
        )
        .join("\n")
    : "";
const status = Number(_error.status ?? _error.statusCode ?? 500);
const is404 = status === 404;
const statusText =
  _error.statusText ??
  _error.statusMessage ??
  (is404 ? "Page Not Found" : "Internal Server Error");
const description = _error.message || _error.toString();
const stack =
  import.meta.dev && !is404
    ? _error.description || `<pre>${stacktrace}</pre>`
    : undefined;
const ErrorTemplate = is404 ? Error404 : Error500;
</script>

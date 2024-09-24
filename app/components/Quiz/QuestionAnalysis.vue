<script setup>
const props = defineProps({
  data: {
    type: Object,
    required: true,
  },
});

const questionsAnalysis = computed(() => {
  // to remove rank object from data
  const filteredData = props.data?.filter(
    (item) => !item.hasOwnProperty("rank")
  );
  return filteredData;
});
</script>

<template>
  <Frame
    v-for="(item, index) in questionsAnalysis"
    :key="index"
    :page-title="'Q' + (index + 1) + '. ' + item.question"
    class="mb-2"
  >
    <div
      v-if="item?.question_media === 'image'"
      class="d-flex align-items-center justify-content-center"
    >
      <img
        :src="`${item?.resource}`"
        :alt="`${item?.resource}`"
        class="rounded img-thumbnail"
      />
    </div>
    <div
      v-if="item?.question_media === 'code'"
      class="d-flex align-items-center justify-content-center"
    >
      <CodeBlockComponent :code="item?.resource" />
    </div>
    <ul style="list-style-type: none; padding-left: 0">
      <li
        v-for="(option, key) in item.options"
        :key="key"
        style="display: flex; align-items: center; padding-left: 20px"
      >
        <span
          v-if="item.correct_answer.includes(key)"
          style="margin-right: 10px"
          >&#10004;</span
        >
        <span
          v-if="
            item.selected_answer.String.includes(key) &&
            !item.correct_answer.includes(key)
          "
          style="margin-right: 10px"
        >
          &#10006;
        </span>
        <span v-if="item?.options_media === 'text'"
          >{{ key }}: {{ option }}</span
        >
        <div
          v-if="item?.options_media === 'image'"
          class="d-flex align-items-center justify-content-center"
        >
          <span>{{ key }}:</span>
          <img
            :src="`${option}`"
            :alt="`${option}`"
            class="rounded img-thumbnail"
          />
        </div>
        <div
          v-if="item?.options_media === 'code'"
          class="mt-2 d-flex align-items-center justify-content-center"
        >
          <span class="mr-1">{{ key }}:</span>
          <CodeBlockComponent :code="option" />
        </div>
      </li>
    </ul>
    <div
      style="
        display: flex;
        flex: 1;
        margin-top: 10px;
        border-top: 1px solid #ccc;
      "
    >
      <div
        v-if="item.response_time > 0"
        style="flex: 1; padding: 10px; border-right: 1px solid #ccc"
      >
        Response Time:
        {{ (item.response_time / 1000).toFixed(2) }} seconds
      </div>
      <div v-else style="flex: 1; padding: 10px; border-right: 1px solid #ccc">
        Response Time: -
      </div>
      <div style="flex: 1; padding: 10px">
        {{ item.is_attend ? "Attempted" : "Not Attempted" }}
      </div>
    </div>
  </Frame>
</template>

<template>
  <NuxtLoadingIndicator />
  <div class="bg-jv-canvas text-jv-ink font-body">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <RunningQuizBanner />
  </div>
</template>
<script setup>
import RunningQuizBanner from "@/components/common/RunningQuizBanner.vue";
</script>

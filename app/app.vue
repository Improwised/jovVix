<template>
  <NuxtLoadingIndicator />
  <div class="bg-jv-canvas text-jv-ink font-body">
    <NuxtLayout>
      <NuxtPage />
    </NuxtLayout>
    <RunningQuizBanner />
    <Notivue v-slot="item">
      <Notification :item="item" :theme="pastelTheme">
        <NotificationProgress :item="item" :theme="pastelTheme" />
      </Notification>
    </Notivue>
  </div>
</template>
<script setup>
import RunningQuizBanner from "@/components/common/RunningQuizBanner.vue";
</script>

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

// Global SEO defaults — inherited by every page and overridden per-page where set.
// Keeps Open Graph / Twitter cards present site-wide so shared links always render
// a preview, even on pages that only declare their own title/description.
const { baseUrl } = useRuntimeConfig().public;
const defaultShareImage = new URL("/og-image.png", baseUrl).href;

useSeoMeta({
  ogSiteName: "jovVix",
  ogType: "website",
  ogLocale: "en_US",
  ogImage: defaultShareImage,
  twitterCard: "summary_large_image",
  twitterImage: defaultShareImage,
});
</script>

<script setup>
import { ArrowLeft } from "lucide-vue-next";
const route = useRoute();
const report = ref(null);
definePageMeta({ layout: "empty", hideSidebar: false });
useSeoMeta({ title: "Room Report - jovVix", robots: "noindex, nofollow" });
</script>
<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 text-jv-ink sm:gap-10 sm:px-6 md:px-8 md:py-6"
  >
    <header class="min-w-0">
      <h1
        class="font-headings text-[38px] leading-none min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
      >
        Room Report
      </h1>
      <div
        class="ml-1 mt-1 h-3 w-40 rounded-full border-b-[3px] border-jv-yellow sm:ml-2 sm:w-48"
        aria-hidden="true"
      />
      <p
        class="mt-4 max-w-2xl text-[15px] font-semibold text-jv-muted sm:text-[16px]"
      >
        Review combined scores across completed quizzes, or choose an individual
        quiz to see its results.
      </p>
    </header>
    <NuxtLink
      :to="
        report?.is_host
          ? '/admin/reports?type=rooms'
          : `/room/${route.params.id}`
      "
      class="inline-flex w-fit items-center gap-2 rounded-full border-2 border-jv-ink bg-jv-white px-5 py-2 font-bold text-jv-ink shadow-brutal-sm transition-transform hover:-rotate-1 active:shadow-none"
      ><ArrowLeft class="size-4" />{{
        report?.is_host ? "Back to reports" : "Back to room"
      }}</NuxtLink
    >
    <RoomStandings
      :room-id="String(route.params.id)"
      detailed
      @report-loaded="report = $event"
    />
  </main>
</template>

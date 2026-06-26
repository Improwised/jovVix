<template>
  <div class="min-h-screen bg-jv-canvas px-4 py-8 sm:px-6 sm:py-10">
    <div class="mx-auto max-w-6xl">
      <header class="mb-8">
        <h1
          class="font-headings text-[28px] leading-tight text-jv-ink sm:text-[36px]"
        >
          Quizzes
        </h1>
        <p class="mt-2 font-body text-sm text-jv-muted sm:text-base">
          Pick what you'd like to do today.
        </p>
      </header>

      <div
        class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 lg:gap-8"
      >
        <ActionCard
          action-title="Create New Quiz"
          @card-click="navigateTo('/admin/quiz/list-quiz?create=1')"
        >
          <template #image>
            <img
              src="@/assets/images/add-file.webp"
              alt=""
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="My Quizzes"
          @card-click="navigateTo('/admin/quiz/list-quiz')"
        >
          <template #image>
            <img
              src="@/assets/images/list.webp"
              alt=""
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="Shared Quizzes"
          @card-click="navigateTo('/admin/quiz/shared-quiz')"
        >
          <template #image>
            <img
              src="@/assets/images/shared-file.webp"
              alt=""
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
        <ActionCard
          action-title="Shared With Me"
          @card-click="navigateTo('/admin/quiz/share-with-me')"
        >
          <template #image>
            <img
              src="@/assets/images/share-with-me.webp"
              alt=""
              class="max-h-full max-w-full object-contain"
            />
          </template>
        </ActionCard>
      </div>
    </div>
  </div>
</template>

<script setup>
import ActionCard from "~~/components/Quiz/ActionCard.vue";

definePageMeta({
  layout: "empty",
});

useSeoMeta({
  title: "My Quizzes - jovVix",
  description: "Create, edit, and organize your interactive quizzes on jovVix.",
  robots: "noindex, nofollow",
});
</script>

<template>
  <div class="hidden items-start justify-end gap-6 lg:flex">
    <NavigationLink
      url="/join"
      url-name="Enter Code"
      class="bg-jv-coral text-white font-[500] py-2"
    />
    <NavigationLink
      url="/admin/quiz/create-quiz"
      url-name="Create Quiz"
      class="bg-jv-coral text-white font-[500] py-2"
    />
  </div>
</template>

<script setup>
import NavigationLink from "@/components/common/NavigationLink.vue";
</script>

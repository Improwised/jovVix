<template>
  <div class="flex flex-col gap-4">
    <div
      v-for="n in 6"
      :key="n"
      class="jv-border-rough flex items-center gap-4 border-2 border-jv-ink bg-jv-white p-4 shadow-brutal-sm sm:p-5"
    >
      <Skeleton class="size-12 shrink-0 rounded-full" />
      <div class="flex flex-1 flex-col gap-2">
        <Skeleton class="h-4 w-3/5" />
        <Skeleton class="h-3 w-2/5" />
        <Skeleton class="h-3 w-1/3" />
      </div>
      <Skeleton class="h-9 w-24 shrink-0" />
    </div>
  </div>
</template>

<script setup>
import { Skeleton } from "@/components/ui/skeleton";
</script>

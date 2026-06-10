<template>
  <Frame page-title="Loading">
    <div class="flex flex-col gap-3">
      <Skeleton class="h-6 w-2/3" />
      <Skeleton class="h-4 w-full" />
      <Skeleton class="h-4 w-5/6" />
      <Skeleton class="h-4 w-4/5" />
      <Skeleton class="h-4 w-3/5" />
      <Skeleton class="mt-2 h-4 w-full" />
      <Skeleton class="h-4 w-11/12" />
    </div>
  </Frame>
</template>

<script setup>
import { Skeleton } from "@/components/ui/skeleton";
</script>

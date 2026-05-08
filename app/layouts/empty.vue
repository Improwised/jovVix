<template>
  <div class="lg:flex lg:h-screen lg:overflow-hidden">
    <Sidebar v-if="!hideSidebar" />
    <div class="lg:flex-1 lg:overflow-y-auto">
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import Sidebar from "@/components/common/Sidebar.vue";

const route = useRoute();
const hideSidebar = computed(() => route.meta.hideSidebar === true);
</script>

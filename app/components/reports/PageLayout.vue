<template>
  <div class="row">
    <div class="col-lg-12 col-md-12 col-12">
      <h3 class="mb-2 fw-bold">Quiz Analysis</h3>
      <ul class="nav nav-tabs">
        <li class="nav-item">
          <NuxtLink
            :class="{ active: route.path === `/admin/reports/${id}` }"
            class="nav-link"
            :to="`/admin/reports/${id}`"
            >Questions</NuxtLink
          >
        </li>
        <li class="nav-item">
          <NuxtLink
            :class="{ active: route.path === `/admin/reports/${id}/participants` }"
            class="nav-link"
            :to="`/admin/reports/${id}/participants`"
            >Participants</NuxtLink
          >
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
const route = useRoute();
const id = route.params.id;
</script>

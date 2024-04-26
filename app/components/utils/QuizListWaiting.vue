<template>
  <div class="d-flex flex-column gap-3">
    <div v-for="n in 6" :key="n">
      <v-skeleton-loader
        class="bg-light p-2 gap-md-5"
        :loading="true"
        type="avatar, paragraph, button"
      >
      </v-skeleton-loader>
    </div>
  </div>
</template>

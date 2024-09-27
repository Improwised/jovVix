<template>
  <v-card class="winner-card text-light">
    <v-img height="150" :src="`${avatar}&scale=20`" cover></v-img>
    <v-card-item>
      <div class="text-center">
        <span class="badge fs-5 bg-light m-1 text-dark"
          >Rank {{ props.winner.rank }}</span
        >
        <h2 class="mb-0 text-h4 text-light">
          {{ props.winner.firstname.toUpperCase() }}
        </h2>
        <div>{{ props.winner.username }}</div>
        <v-row class="mt-6">
          <v-col>
            <h2 class="mb-0 text-h4 text-light">{{ props.winner.score }}</h2>
            <div>Score</div>
          </v-col>
        </v-row>
      </div>
    </v-card-item>
  </v-card>
</template>

<script setup>
const props = defineProps({
  winner: {
    type: Object,
    required: true,
    default: () => {
      return {};
    },
  },
});

const avatar = computed(() => {
  return getAvatarUrlByName(props.winner?.img_key);
});
</script>

<style>
.winner-card {
  background: #ffffff2b !important;
  border-radius: 100% !important;
}
</style>

<template>
  <div class="winner-card">
    <img
      v-if="props.winner.rank == 1"
      src="@/assets/images/medal/1.webp"
      class="bg-image"
      alt="medal"
    />
    <img
      v-if="props.winner.rank == 2"
      src="@/assets/images/medal/2.webp"
      class="bg-image"
      alt="medal"
    />
    <img
      v-if="props.winner.rank == 3"
      src="@/assets/images/medal/3.webp"
      class="bg-image"
      alt="medal"
    />
    <img :src="avatar + '&scale=60'" class="avatar-image overlay" />
    <div class="winner-details overlay">
      <div class="card-body text-center py-3">
        <!-- <span class="badge fs-5 bg-light m-1 text-dark">Rank {{ props.winner.rank }}</span> -->
        <h2 class="mb-0 text-h4 text-light">
          {{ props.winner.firstname.toUpperCase() }}
        </h2>
        <div>{{ props.winner.username }}</div>
        <h2 class="mb-0 text-h4 text-light">{{ props.winner.score }}</h2>
      </div>
    </div>
  </div>
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
  position: relative;
  width: 100%;
  max-width: 300px;
  height: 75%;
  max-width: 500px;
}

.bg-image {
  height: 100%;
  width: 100%;
}

.avatar-image {
  display: block;
  width: 35%;
  height: 40%;
}

.overlay {
  position: absolute;
  top: 15%;
  color: #f1f1f1;
  width: 100%;
  opacity: 1;
  color: white;
  text-align: center;
}

.winner-details {
  background-color: rgba(0, 0, 0, 0);
  border-radius: 10px;
  top: 64%;
  width: 80%;
  left: 10%;
}
</style>

<template>
  <article
    class="winner-card"
    :aria-label="`${getOrdinal(props.winner.rank)} place winner: ${
      props.winner.firstname
    } ${props.winner.username} with score ${props.winner.score}`"
    role="article"
  >
    <img
      v-if="props.winner.rank == 1"
      src="@/assets/images/medal/1.webp"
      class="bg-image"
      alt="First place gold medal background"
      role="img"
    />
    <img
      v-if="props.winner.rank == 2"
      src="@/assets/images/medal/2.webp"
      class="bg-image"
      alt="Second place silver medal background"
      role="img"
    />
    <img
      v-if="props.winner.rank == 3"
      src="@/assets/images/medal/3.webp"
      class="bg-image"
      alt="Third place bronze medal background"
      role="img"
    />
    <img
      :src="avatar + '&scale=60'"
      class="avatar-image overlay"
      :alt="`Avatar for ${props.winner.firstname} ${props.winner.username}`"
      role="img"
    />
    <div class="winner-details overlay">
      <div class="card-body text-center py-3">
        <h2
          :id="`winner-name-${props.winner.rank}`"
          class="mb-0 text-h4 text-light"
        >
          {{ props.winner.firstname.toUpperCase() }}
        </h2>
        <div
          class="text-light"
          :aria-labelledby="`winner-name-${props.winner.rank}`"
        >
          {{ props.winner.username }}
        </div>
        <div
          class="mb-0 text-h4 text-light"
          role="text"
          :aria-label="`Score: ${props.winner.score} points`"
        >
          {{ props.winner.score }}
        </div>
      </div>
    </div>
  </article>
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

const getOrdinal = (rank) => {
  const suffixes = ["th", "st", "nd", "rd"];
  const v = rank % 100;
  return rank + (suffixes[(v - 20) % 10] || suffixes[v] || suffixes[0]);
};
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

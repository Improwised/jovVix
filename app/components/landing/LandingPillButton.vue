<template>
  <NuxtLink
    :to="to"
    class="landing-pill-button"
    :style="buttonStyle"
    :aria-label="ariaLabel || undefined"
  >
    <span class="landing-pill-button__text">
      <slot />
    </span>
  </NuxtLink>
</template>

<script setup>
import { computed } from "vue";

const props = defineProps({
  to: {
    type: String,
    required: true,
  },
  width: {
    type: String,
    default: "148px",
  },
  ariaLabel: {
    type: String,
    default: "",
  },
});

const buttonStyle = computed(() => ({
  width: props.width,
}));
</script>

<style scoped>
.landing-pill-button {
  align-items: center;
  background: var(--jv-coral);
  border: 2px solid var(--jv-ink);
  border-radius: 135px 120px 120px 135px;
  box-shadow: 2px 2px 0 var(--jv-ink);
  color: #ffffff;
  display: inline-flex;
  font-family: var(--font-body);
  font-size: 20px;
  font-weight: 600;
  height: 48px;
  justify-content: center;
  line-height: 28px;
  max-width: 100%;
  padding: 8px 24px;
  text-align: center;
  text-decoration: none;
  transform: matrix(1, 0.01, -0.04, 1, 0, 0);
}

.landing-pill-button:hover,
.landing-pill-button:focus {
  background: var(--jv-coral);
  color: #ffffff;
  text-decoration: none;
}

.landing-pill-button__text {
  transform: matrix(1, 0.01, -0.04, 1, 0, 0);
  white-space: nowrap;
}
</style>

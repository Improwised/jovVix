<template>
  <header class="sticky top-0 bg-jv-canvas">
    <div
      class="absolute inset-0 pointer-events-none jv-grid"
      aria-hidden="true"
    ></div>
    <div
      class="container mx-auto max-w-7xl flex items-center justify-between gap-4 px-6 sm:px-8 py-2.5 lg:py-5"
    >
      <NuxtLink
        to="/"
        class="flex items-center gap-2 text-[20px] sm:text-[22px] md:text-[24px] font-black tracking-[-0.5px] text-jv-ink no-underline w-[120px] lg:w-[150px] h-auto"
      >
        <img src="@/assets/images/jovvix-logo.svg" />
      </NuxtLink>

      <ul
        class="hidden lg:flex items-center font-headings gap-8 text-[18px] text-jv-ink"
      >
        <li v-for="link in navLinks" :key="link.label" class="text-center">
          <NuxtLink
            :href="link.href"
            class="hover:text-jv-coral transition-colors duration-150"
          >
            {{ link.label }}
          </NuxtLink>
        </li>
      </ul>

      <NavigationLink
        url="/join"
        url-name="Play a Quiz"
        class="bg-jv-yellow-2"
      />
    </div>
  </header>
</template>

<script setup>
import NavigationLink from "@/components/common/NavigationLink.vue";
const navLinks = [
  { href: "https://jovvix.com/", label: "Home" },
  { href: "https://jovvix.com/docs/user-guide", label: "User Guide" },
  { href: "https://jovvix.com/docs/developer-guide", label: "Developer Guide" },
  {
    href: "https://jovvix.com/docs/deployment-guide",
    label: "Deployment Guide",
  },
];
</script>

<template>
  <!-- Desktop sidebar (lg and up) -->
  <aside
    class="hidden lg:flex bg-jv-ivory min-w-56 h-screen flex-col justify-between border-r-[4px] border-jv-ink px-8 py-5 overflow-y-hidden"
  >
    <div class="flex flex-col gap-4">
      <NuxtLink
        to="/"
        class="flex items-center gap-[6px] text-[30px] font-black tracking-[-0.75px] text-jv-ink no-underline"
      >
        <span
          class="grid size-9 rotate-[-6deg] place-items-center rounded-[135px_120px_120px_135px] border-2 border-jv-ink bg-jv-yellow text-jv-ink shadow-[2px_2px_0_#2D2D2D]"
        >
          <Zap class="size-5 rotate-[-6deg]" :stroke-width="2.4" />
        </span>
        <span>JovVix</span>
      </NuxtLink>

      <NavigationLink url="/admin/quiz/create-quiz" url-name="Create Quiz">
        <Plus class="size-[18px]" />
      </NavigationLink>

      <NavigationLink
        url="/"
        url-name="Home"
        class="flex justify-baseline py-1.5 gap-3 rounded-full border-2 border-dashed border-jv-ink/40 bg-jv-yellow/20 shadow-none hover:rotate-[2deg]"
      >
        <Home class="size-5 text-jv-coral" :stroke-width="2.4" />
      </NavigationLink>
    </div>

    <div class="flex flex-col gap-4">
      <NavigationLink url="/account/login" url-name="Sign In"></NavigationLink>
      <NavigationLink
        url="/account/register"
        url-name="Sign Up"
      ></NavigationLink>
    </div>
  </aside>

  <!-- Mobile/Tablet header (md and below) -->
  <header
    class="lg:hidden sticky top-0 z-40 border-b-[3px] border-jv-ink bg-jv-ivory/95 px-4 sm:px-6 py-3 backdrop-blur"
  >
    <div class="flex items-center justify-between gap-3">
      <NuxtLink
        to="/"
        class="flex items-center gap-[6px] text-[22px] sm:text-[26px] font-black tracking-[-0.75px] text-jv-ink no-underline"
      >
        <span
          class="grid size-8 sm:size-9 rotate-[-6deg] place-items-center rounded-[135px_120px_120px_135px] border-2 border-jv-ink bg-jv-yellow text-jv-ink shadow-[2px_2px_0_#2D2D2D]"
        >
          <Zap class="size-4 sm:size-5 rotate-[-6deg]" :stroke-width="2.4" />
        </span>
        <span>JovVix</span>
      </NuxtLink>

      <button
        type="button"
        class="grid size-10 place-items-center rounded-[10px] border-[3px] border-jv-ink bg-jv-yellow text-jv-ink shadow-brutal-sm active:shadow-none active:translate-x-[2px] active:translate-y-[2px] transition-transform"
        :aria-expanded="open"
        aria-label="Toggle navigation menu"
        aria-controls="mobile-nav"
        @click="open = !open"
      >
        <X v-if="open" class="size-5" :stroke-width="2.4" />
        <Menu v-else class="size-5" :stroke-width="2.4" />
      </button>
    </div>

    <nav
      v-show="open"
      id="mobile-nav"
      class="mt-4 grid grid-cols-1 min-[420px]:grid-cols-2 sm:grid-cols-3 gap-3"
    >
      <NavigationLink
        url="/"
        url-name="Home"
        class="!py-2 !px-4 !text-sm w-full"
        @click="open = false"
      >
        <Home class="size-4 text-jv-coral" :stroke-width="2.4" />
      </NavigationLink>
      <NavigationLink
        url="/join"
        url-name="Enter Code"
        class="!py-2 !px-4 !text-sm w-full bg-jv-coral !text-white"
        @click="open = false"
      />
      <NavigationLink
        url="/admin/quiz/create-quiz"
        url-name="Create Quiz"
        class="!py-2 !px-4 !text-sm w-full"
        @click="open = false"
      >
        <Plus class="size-4" />
      </NavigationLink>
      <NavigationLink
        url="/account/login"
        url-name="Sign In"
        class="!py-2 !px-4 !text-sm w-full"
        @click="open = false"
      />
      <NavigationLink
        url="/account/register"
        url-name="Sign Up"
        class="!py-2 !px-4 !text-sm w-full"
        @click="open = false"
      />
    </nav>
  </header>
</template>

<script setup>
import { ref } from "vue";
import { Home, Plus, Zap, Menu, X } from "lucide-vue-next";
import NavigationLink from "@/components/common/NavigationLink.vue";

const open = ref(false);
</script>

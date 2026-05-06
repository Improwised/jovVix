<template>
  <aside
    class="fixed inset-y-0 left-0 z-40 hidden w-56 flex-col border-r-[4px] border-jv-ink bg-[#F0EBD8] px-5 py-5 shadow-[3px_3px_0_#2D2D2D] lg:flex"
  >
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

    <NuxtLink
      to="/admin/quiz/create-quiz"
      class="mt-7 flex h-14 items-center justify-center gap-2 rounded-[9px] border-2 border-jv-ink !bg-jv-yellow px-5 text-[20px] font-normal !text-jv-ink no-underline shadow-[3px_3px_0_#2D2D2D]"
    >
      <Plus class="size-[18px]" />
      <span>Create Quiz</span>
    </NuxtLink>

    <nav class="mt-8">
      <NuxtLink
        to="/"
        class="flex h-11 rotate-[-1deg] items-center gap-3 rounded-full border-2 border-dashed border-jv-ink/40 !bg-jv-yellow/20 px-4 text-[20px] font-normal !text-jv-ink no-underline"
      >
        <Home class="size-5 text-jv-coral" :stroke-width="2.4" />
        <span>Home</span>
      </NuxtLink>
    </nav>

    <div class="mt-auto space-y-8 gap-4 flex flex-col">
      <NuxtLink
        to="/account/login"
        class="flex h-14 items-center justify-center rounded-[9px] border-2 border-jv-ink !bg-jv-yellow px-5 text-[20px] font-normal !text-jv-ink no-underline shadow-[3px_3px_0_#2D2D2D]"
      >
        Sign In
      </NuxtLink>
      <NuxtLink
        to="/account/register"
        class="flex h-14 items-center justify-center rounded-[9px] border-2 border-jv-ink !bg-jv-yellow px-5 text-[20px] font-normal !text-jv-ink no-underline shadow-[3px_3px_0_#2D2D2D]"
      >
        Sign Up
      </NuxtLink>
    </div>
  </aside>

  <header
    class="sticky top-0 z-40 border-b-[3px] border-jv-ink bg-[#eee9d7]/95 px-4 py-3 backdrop-blur lg:hidden"
  >
    <div class="flex items-center justify-between gap-3">
      <NuxtLink
        to="/"
        class="flex items-center gap-2 text-2xl font-black text-jv-ink no-underline"
      >
        <span
          class="grid size-9 place-items-center rounded-full border-[3px] border-jv-ink bg-jv-yellow text-jv-ink"
        >
          <Zap class="size-5" :stroke-width="2.6" />
        </span>
        <span>JovVix</span>
      </NuxtLink>
      <NuxtLink
        to="/admin/quiz/create-quiz"
        class="grid size-10 place-items-center rounded-[10px] border-[3px] border-jv-ink !bg-jv-yellow !text-jv-ink shadow-brutal-sm"
        aria-label="Create Quiz"
      >
        <Plus class="size-5" />
      </NuxtLink>
    </div>

    <div
      class="mt-3 grid grid-cols-2 gap-x-5 gap-y-3 text-sm font-bold min-[420px]:grid-cols-4"
    >
      <NuxtLink
        to="/"
        class="flex h-10 items-center justify-center gap-1 rounded-full border-2 border-dashed border-[#b8aa75] !bg-jv-yellow-soft/50 !text-jv-ink no-underline"
      >
        <Home class="size-4 text-jv-coral" />
        Home
      </NuxtLink>
      <NuxtLink
        to="/join"
        class="flex h-10 items-center justify-center rounded-full border-[3px] border-jv-ink !bg-jv-coral !text-white no-underline shadow-brutal-sm"
      >
        Enter Code
      </NuxtLink>
      <NuxtLink
        to="/account/login"
        class="flex h-10 items-center justify-center rounded-[10px] border-[3px] border-jv-ink !bg-jv-yellow !text-jv-ink no-underline shadow-brutal-sm"
      >
        Sign In
      </NuxtLink>
      <NuxtLink
        to="/account/register"
        class="flex h-10 items-center justify-center rounded-[10px] border-[3px] border-jv-ink !bg-jv-yellow !text-jv-ink no-underline shadow-brutal-sm"
      >
        Sign Up
      </NuxtLink>
    </div>
  </header>
</template>

<script setup>
import { Home, Plus, Zap } from "lucide-vue-next";
</script>

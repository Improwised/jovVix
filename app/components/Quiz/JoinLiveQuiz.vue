<template>
  <div
    class="relative bg-jv-yellow/85 mx-auto w-full max-w-[1386px] jv-border-rough shadow-brutal-lg p-6 sm:p-8 md:p-10 text-center -rotate-1 mt-4"
  >
    <div
      class="mx-auto inline-flex size-12 sm:size-14 items-center justify-center bg-jv-coral text-white jv-border-rough shadow-brutal-sm rotate-3"
      aria-hidden="true"
    >
      <Gamepad2 class="size-6 sm:size-7" />
    </div>

    <h2
      class="mt-4 sm:mt-5 font-headings text-jv-ink text-[24px] min-[420px]:text-[28px] sm:text-[34px] md:text-[40px] leading-[1.15]"
    >
      Join a Live Quiz
    </h2>
    <p
      class="mt-2 sm:mt-3 font-body text-jv-ink/80 text-[14px] sm:text-[16px] md:text-[18px] leading-[1.55] max-w-md mx-auto px-2"
    >
      Enter your event invitation code below to jump straight into the action
      and compete with friends!
    </p>

    <form
      action="https://app.jovvix.com/join"
      method="get"
      class="mt-6 flex flex-col sm:flex-row items-stretch gap-3 sm:gap-0 bg-white jv-border-rough shadow-brutal p-2 max-w-lg mx-auto"
    >
      <div class="flex items-center flex-1 gap-2 px-3 font-body text-jv-ink">
        <Hash size="18" class="text-jv-muted shrink-0" />
        <input
          type="text"
          name="code"
          inputmode="numeric"
          autocomplete="off"
          placeholder="Enter 6-digit code"
          class="w-full bg-transparent outline-none text-[16px] py-2 placeholder:text-jv-muted/70"
        />
      </div>
      <Button type="submit" size="sm" :tilt="2" variant="accent">
        Join Game
      </Button>
    </form>
  </div>
</template>

<script setup>
import { Hash, Gamepad2 } from "lucide-vue-next";
import { Button } from "@/components/ui/button";
</script>

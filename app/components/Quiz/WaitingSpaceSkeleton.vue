<script setup>
import { Info, Keyboard, Users } from "lucide-vue-next";
import { Skeleton } from "@/components/ui/skeleton";
</script>

<template>
  <main
    class="flex min-h-screen flex-col gap-8 bg-jv-canvas px-4 py-5 text-jv-ink sm:gap-10 sm:px-6 md:px-8 md:py-6"
    aria-busy="true"
    aria-live="polite"
  >
    <section class="mx-auto flex w-full max-w-[1220px] flex-1 flex-col">
      <header class="mb-7 text-center sm:mb-9">
        <h1
          class="font-headings text-[38px] leading-none text-jv-ink min-[420px]:text-[44px] sm:text-[52px] md:text-[56px]"
        >
          Quiz Lobby
        </h1>
        <p
          class="mx-auto mt-2 max-w-[680px] text-[15px] leading-[1.5] text-jv-muted sm:text-[18px] md:text-[20px]"
        >
          Share the code with players. The quiz will start when you're ready!
        </p>
      </header>

      <div
        class="grid flex-1 gap-6 md:gap-7 xl:grid-cols-[minmax(0,1fr)_minmax(340px,0.95fr)] xl:items-stretch"
      >
        <div
          class="relative flex min-h-0 rotate-[-0.4deg] flex-col overflow-hidden bg-jv-white shadow-brutal-sm jv-border-rough sm:shadow-brutal-lg xl:min-h-[640px]"
        >
          <span
            class="absolute left-1/2 top-[-2px] z-10 h-4 w-12 -translate-x-1/2 rotate-[-1deg] bg-jv-coral"
            aria-hidden="true"
          ></span>

          <div class="bg-jv-yellow px-4 pb-7 pt-6 sm:px-8 sm:pb-9 md:px-10">
            <h2
              class="font-headings text-[30px] leading-tight text-jv-ink min-[420px]:text-[34px] sm:text-[42px] lg:text-[48px]"
            >
              Ready, Steady, Go!
            </h2>

            <div
              class="mt-4 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-brutal-sm sm:flex-row sm:items-center sm:justify-between sm:p-4"
            >
              <div
                class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
              >
                <p
                  class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                >
                  Join at
                </p>
                <Skeleton
                  class="h-[20px] w-[200px] max-w-full rounded-[4px] bg-jv-ink/10 sm:h-[24px] sm:w-[260px]"
                />
              </div>
              <Skeleton
                class="h-11 w-full shrink-0 rotate-[-1deg] rounded-[999px] border-[3px] border-jv-ink bg-jv-coral/40 shadow-brutal-sm sm:h-12 sm:w-[120px]"
              />
            </div>
          </div>

          <div class="flex flex-1 flex-col px-4 py-5 sm:px-8 sm:py-6 md:px-10">
            <div
              class="-mt-11 flex flex-col gap-3 jv-border-rough bg-jv-white p-3 shadow-brutal-sm sm:-mt-14 sm:flex-row sm:items-center sm:justify-between sm:p-4 sm:shadow-brutal"
            >
              <div class="flex min-w-0 items-center gap-3 sm:gap-4">
                <span
                  class="grid size-10 shrink-0 place-items-center border-[2px] border-jv-ink bg-jv-mint text-jv-ink sm:size-11"
                >
                  <Keyboard class="size-5" :stroke-width="2.2" />
                </span>
                <div
                  class="flex min-w-0 flex-wrap items-baseline gap-x-2 gap-y-1"
                >
                  <p
                    class="shrink-0 font-body text-[12px] font-black uppercase tracking-[0.08em] text-jv-muted"
                  >
                    Quiz code
                  </p>
                  <Skeleton
                    class="h-[24px] w-[160px] max-w-full rounded-[4px] bg-jv-ink/10 sm:h-[28px] sm:w-[180px]"
                  />
                </div>
              </div>
              <Skeleton
                class="h-11 w-full shrink-0 rotate-[-1deg] rounded-[999px] border-[3px] border-jv-ink bg-jv-coral/40 shadow-brutal-sm sm:h-12 sm:w-[120px]"
              />
            </div>

            <div
              class="my-5 flex items-center justify-center gap-3 text-center font-body text-[11px] font-black uppercase tracking-[0.08em] text-jv-muted min-[420px]:text-[12px] sm:my-6 sm:gap-4 sm:text-[13px]"
            >
              <span
                class="h-px min-w-8 flex-1 border-t-2 border-dashed border-jv-ink/40 sm:max-w-24"
              ></span>
              <span class="shrink-0">Or scan QR code</span>
              <span
                class="h-px min-w-8 flex-1 border-t-2 border-dashed border-jv-ink/40 sm:max-w-24"
              ></span>
            </div>

            <div class="flex justify-center">
              <div
                class="grid size-[176px] place-items-center bg-jv-white p-3 shadow-brutal-sm jv-border-rough min-[420px]:size-[196px] sm:size-[220px] sm:p-4 md:size-[240px]"
              >
                <Skeleton class="size-full rounded-none bg-jv-ink/10" />
              </div>
            </div>

            <div class="mt-auto pt-5 sm:pt-6">
              <Skeleton
                class="mx-auto h-[52px] w-full max-w-[390px] rotate-[-1deg] rounded-[999px] border-[3px] border-jv-ink bg-jv-mint/50 shadow-brutal-sm sm:h-16"
              />
              <p
                class="mt-4 flex items-center justify-center gap-2 text-center font-body text-[12px] leading-[1.4] text-jv-muted sm:text-[13px]"
              >
                <Info class="size-4" :stroke-width="2.3" />
                <span>Host can start the quiz at any time</span>
              </p>
            </div>
          </div>
        </div>

        <aside
          class="relative flex min-h-0 rotate-[0.4deg] flex-col bg-jv-white px-4 py-6 shadow-brutal-sm jv-border-rough sm:px-8 sm:py-8 sm:shadow-brutal-lg md:px-10 xl:min-h-[640px]"
        >
          <span
            class="absolute left-1/2 top-[-10px] h-4 w-12 -translate-x-1/2 rotate-[1deg] bg-jv-salmon"
            aria-hidden="true"
          ></span>

          <div
            class="flex items-start justify-between gap-3 sm:items-center sm:gap-4"
          >
            <div class="flex min-w-0 items-center gap-3 sm:gap-4">
              <span
                class="grid size-10 shrink-0 rotate-[-3deg] place-items-center jv-border-rough bg-jv-mint sm:size-12"
              >
                <Users class="size-5 sm:size-6" :stroke-width="2.4" />
              </span>
              <h2
                class="min-w-0 break-words font-headings text-[28px] leading-tight text-jv-ink min-[420px]:text-[32px] sm:text-[42px] xl:text-[48px]"
              >
                Participants
              </h2>
            </div>
            <Skeleton
              class="size-11 shrink-0 rotate-[3deg] rounded-none border-[3px] border-jv-ink bg-jv-yellow/60 shadow-brutal-sm sm:size-14"
            />
          </div>

          <div
            class="my-5 border-t-2 border-dashed border-jv-ink/20 sm:my-6"
          ></div>

          <div class="flex flex-col gap-3 sm:gap-4">
            <div
              v-for="i in 3"
              :key="i"
              class="flex min-h-14 items-center justify-between gap-3 jv-border-rough bg-jv-white px-3 py-2 shadow-[2px_2px_0_#2D2D2D] sm:min-h-16"
            >
              <div class="flex min-w-0 items-center gap-3">
                <Skeleton
                  class="size-9 shrink-0 rotate-[-2deg] rounded-none border-[2px] border-jv-ink bg-jv-ink/10 sm:size-10"
                />
                <Skeleton
                  class="h-[18px] w-[140px] rounded-[4px] bg-jv-ink/10 sm:h-[22px] sm:w-[180px]"
                />
              </div>
              <Skeleton
                class="size-8 shrink-0 rotate-[2deg] rounded-none border-[2px] border-jv-ink bg-jv-mint/40 sm:size-9"
              />
            </div>
          </div>

          <div class="mt-auto pt-6 sm:pt-8">
            <div class="border-t-2 border-jv-ink/15 pt-5 text-center sm:pt-6">
              <div class="mb-3 flex justify-center gap-2 sm:mb-4">
                <span
                  class="size-2.5 rounded-full bg-jv-coral sm:size-3"
                ></span>
                <span
                  class="size-2.5 rounded-full bg-jv-yellow sm:size-3"
                ></span>
                <span class="size-2.5 rounded-full bg-jv-mint sm:size-3"></span>
              </div>
              <p
                class="font-body text-[18px] text-jv-muted sm:text-[22px] md:text-[24px]"
              >
                Waiting for more players...
              </p>
            </div>
          </div>
        </aside>
      </div>
    </section>
    <span class="sr-only">Loading quiz lobby…</span>
  </main>
</template>

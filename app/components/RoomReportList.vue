<script setup>
import {
  ChevronUp,
  ChevronDown,
  ChevronRight,
  ChevronsUpDown,
} from "lucide-vue-next";
const props = defineProps({ date: { type: String, default: "" } });
const sortKey = ref("finished");
const sortDirection = ref("desc");
const columns = [
  { key: "code", label: "Room name" },
  { key: "participants", label: "Total participants", center: true },
  { key: "started", label: "Started at" },
  { key: "finished", label: "Ended at" },
  { key: "rounds", label: "Completed quizzes", center: true },
  { key: "ended", label: "Status", center: true },
];
const sortBy = (key) => {
  sortDirection.value =
    sortKey.value === key && sortDirection.value === "asc" ? "desc" : "asc";
  sortKey.value = key;
};
const { apiUrl } = useRuntimeConfig().public;
const { data, error, status, refresh } = await useFetch(
  `${apiUrl}/room-reports`,
  { credentials: "include", headers: useRequestHeaders(["cookie"]) }
);
const datePart = (value, time = false) => {
  if (!value) return "—";
  const date = new Date(value);
  return time
    ? date.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
    : date.toLocaleDateString("en-GB", {
        day: "2-digit",
        month: "short",
        year: "numeric",
      });
};
const rooms = computed(() => {
  const after = props.date ? new Date(props.date).getTime() : null;
  return (data.value?.data || [])
    .filter(
      (room) => after === null || new Date(room.started).getTime() >= after
    )
    .slice()
    .sort((a, b) => {
      const key = sortKey.value;
      const sortVal = (room) => {
        if (key === "started" || key === "finished") {
          // Open rooms have no finished time — sort by started instead
          const val = room[key] || (key === "finished" ? room.started : null);
          return val ? new Date(val).getTime() : 0;
        }
        return room[key];
      };
      const x = sortVal(a);
      const y = sortVal(b);
      return (
        (x < y ? -1 : x > y ? 1 : 0) * (sortDirection.value === "asc" ? 1 : -1)
      );
    });
});
</script>
<template>
  <section class="flex flex-col gap-3" aria-label="Room reports">
    <div
      v-if="status === 'pending'"
      role="status"
      class="jv-border-rough bg-jv-white p-5 text-center text-lg font-semibold text-jv-muted shadow-brutal-sm"
    >
      Loading room reports…
    </div>
    <div
      v-else-if="error"
      role="alert"
      class="jv-border-rough bg-jv-white p-5 text-lg font-semibold text-jv-coral shadow-brutal-sm"
    >
      Unable to load room reports.
      <button type="button" class="cursor-pointer underline" @click="refresh">
        Try again
      </button>
    </div>
    <template v-else>
      <div
        class="hidden grid-cols-[2fr_1fr_1fr_1fr_1fr_1fr] gap-4 px-5 py-2 text-xs font-black uppercase tracking-wider text-jv-muted lg:grid"
      >
        <button
          v-for="column in columns"
          :key="column.key"
          type="button"
          class="flex cursor-pointer items-center gap-1 text-left uppercase transition-colors hover:text-jv-ink"
          :class="{ 'justify-center text-center': column.center }"
          :aria-label="`Sort by ${column.label}`"
          @click="sortBy(column.key)"
        >
          {{ column.label
          }}<ChevronUp
            v-if="sortKey === column.key && sortDirection === 'asc'"
            class="size-4 shrink-0"
          /><ChevronDown
            v-else-if="sortKey === column.key"
            class="size-4 shrink-0"
          /><ChevronsUpDown v-else class="size-4 shrink-0 opacity-40" />
        </button>
      </div>
      <p
        v-if="!rooms.length"
        class="jv-border-rough bg-jv-white p-5 text-center text-lg font-semibold text-jv-muted shadow-brutal-sm"
      >
        {{
          props.date
            ? "No rooms match the selected date."
            : "No room reports yet."
        }}
      </p>
      <NuxtLink
        v-for="room in rooms"
        :key="room.id"
        :to="`/room/report/${room.id}`"
        class="relative grid grid-cols-3 items-start gap-x-2 gap-y-2 jv-border-rough bg-jv-white pl-3 pr-9 py-2.5 sm:pl-5 sm:pr-10 lg:items-center lg:gap-x-4 lg:py-4 text-jv-ink shadow-brutal-sm transition-all hover:-translate-y-0.5 hover:shadow-brutal-md lg:grid-cols-[2fr_1fr_1fr_1fr_1fr_1fr]"
      >
        <ChevronRight
          class="pointer-events-none absolute right-3 top-1/2 size-5 -translate-y-1/2 text-jv-muted"
          aria-hidden="true"
        />
        <div class="col-span-full min-w-0 lg:col-span-1 lg:pr-4">
          <h3
            class="font-headings text-[20px] font-bold leading-tight truncate lg:text-[24px]"
          >
            Room {{ room.code }}
          </h3>
        </div>
        <div
          class="text-[16px] font-bold text-jv-ink lg:text-center lg:text-[22px] lg:font-black"
        >
          <span
            class="mb-1 block text-[11px] font-bold uppercase leading-tight text-jv-muted lg:hidden"
            >Participants:</span
          >{{ room.participants ?? "—" }}
        </div>
        <div
          v-for="(value, key) in {
            Started: room.started,
            Ended: room.finished,
          }"
          :key="key"
          class="hidden lg:block"
        >
          <div
            class="inline-block text-[15px] font-bold leading-snug text-jv-ink lg:block"
          >
            {{ datePart(value) }}
            <div v-if="value" class="text-[13px] font-semibold text-jv-muted">
              {{ datePart(value, true) }}
            </div>
          </div>
        </div>
        <div
          class="text-[16px] font-bold text-jv-ink lg:text-center lg:text-[22px] lg:font-black"
        >
          <span
            class="mb-1 block text-[11px] font-bold uppercase leading-tight text-jv-muted lg:hidden"
            >Quizzes:</span
          >{{ room.rounds }}
        </div>
        <div class="lg:text-center">
          <span
            class="mb-1 block text-[11px] font-bold uppercase leading-tight text-jv-muted lg:hidden"
            >Status:</span
          >
          <span
            class="inline-flex rounded-full border-2 border-jv-ink px-1.5 py-0.5 text-xs font-black lg:px-3 lg:py-1"
            :class="room.ended ? 'bg-jv-lavender' : 'bg-jv-mint'"
            >{{ room.ended ? "Ended" : "Open" }}</span
          >
        </div>
      </NuxtLink>
    </template>
  </section>
</template>

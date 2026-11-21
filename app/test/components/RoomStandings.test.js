import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import RoomStandings from "~/components/RoomStandings.vue";

const ROOM = "11111111-0000-4000-8000-000000000001";
const SESSION = "6f1c2a3e-0000-4000-8000-000000000001";

const vue = await vi.hoisted(() => import("vue"));
mockNuxtImport("useFetch", () => () => ({
  data: vue.ref({
    data: {
      is_host: true,
      ended: true,
      standings: [],
      quizzes: [{ id: SESSION, title: "Maps & Rivers" }],
    },
  }),
  error: vue.ref(null),
  status: vue.ref("success"),
  refresh: () => {},
}));
mockNuxtImport("useRoute", () => () => ({
  path: `/room/report/${ROOM}`,
  query: { tab: "quizzes" },
}));

// GitHub #38: the host's report link carries the quiz title, which names the PDF.
describe("RoomStandings", () => {
  it("links the host to each quiz report with its title", async () => {
    const Parent = defineComponent({
      render: () =>
        h(Suspense, null, {
          default: () => h(RoomStandings, { roomId: ROOM, detailed: true }),
        }),
    });
    const wrapper = mount(Parent, {
      global: {
        stubs: {
          NuxtLink: { props: ["to"], template: '<a :href="to"><slot /></a>' },
        },
      },
    });
    await flushPromises();
    const link = wrapper
      .findAll("a")
      .find((a) => a.text().includes("View analysis"));
    expect(link.attributes("href")).toBe(
      `/admin/reports/${SESSION}?room=${ROOM}&title=Maps%20%26%20Rivers`
    );
  });
});

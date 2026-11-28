import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import ReportsPage from "~/pages/admin/reports/index.vue";

const vue = await vi.hoisted(() => import("vue"));
mockNuxtImport("useFetch", () => () => ({
  data: vue.ref({
    data: {
      count: 1,
      data: [
        {
          id: "6f1c2a3e-0000-4000-8000-000000000001",
          title: "Maps%20and%20Rivers",
          description: { String: "", Valid: false },
          activated_from: { Time: "20-Nov-2026 10:00:00" },
          activated_to: { Time: "20-Nov-2026 10:30:00" },
          questions: 3,
          participants: 4,
          correct_answers: 5,
        },
      ],
    },
  }),
  error: vue.ref(null),
  pending: vue.ref(false),
}));
const nav = vi.hoisted(() => vi.fn());
mockNuxtImport("navigateTo", () => nav);

// GitHub #38: the report page names the PDF after the quiz, so the list passes the title along.
describe("Reports list", () => {
  it("opens a report with its decoded quiz title", async () => {
    const Parent = defineComponent({
      render: () => h(Suspense, null, { default: () => h(ReportsPage) }),
    });
    const wrapper = mount(Parent, {
      global: { stubs: { Pagination: true, RoomReportList: true } },
    });
    await flushPromises();
    nav.mockClear();

    const link = wrapper.find('a[href^="/admin/reports/"]');
    await link.trigger("click");

    expect(nav).toHaveBeenCalledWith({
      path: "/admin/reports/6f1c2a3e-0000-4000-8000-000000000001",
      query: { title: "Maps and Rivers" },
    });
  });
});

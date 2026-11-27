import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import ReportsList from "~/pages/admin/reports/index.vue";
import ReportDetail from "~/pages/admin/reports/[id]/index.vue";
import QuestionPage from "~/pages/admin/quiz/list-quiz/[quiz_id]/[question_id].vue";

// Every useFetch on these pages fails as an expired login does.
const vue = await vi.hoisted(() => import("vue"));
mockNuxtImport("useFetch", () => () => ({
  data: vue.ref(null),
  error: vue.ref({ statusCode: 401, data: { code: 401 } }),
  pending: vue.ref(false),
  status: vue.ref("error"),
  refresh: () => {},
}));
const nav = vi.hoisted(() => vi.fn());
mockNuxtImport("navigateTo", () => nav);

// GitHub #77: registered-only pages that printed the raw fetch error now send the person to sign in.
describe("pages with an expired session", () => {
  beforeEach(() => nav.mockReset());

  it.each([
    ["reports list", ReportsList],
    ["report detail", ReportDetail],
    ["question editor", QuestionPage],
  ])("%s goes to the login page", async (_name, Page) => {
    vi.stubGlobal(
      "fetch",
      vi.fn(async () => ({ ok: false, status: 401, json: async () => ({}) }))
    );
    const Parent = defineComponent({
      render: () => h(Suspense, null, { default: () => h(Page) }),
    });
    mount(Parent, {
      global: {
        mocks: { $route: { query: {} } },
        stubs: { NuxtLink: true, NuxtImg: true },
      },
    });
    await flushPromises();
    expect(nav).toHaveBeenCalledWith(
      expect.stringMatching(/^\/account\/login\?returnTo=/),
      { replace: true }
    );
  });
});

import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, beforeEach, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import ReportPage from "~/pages/admin/reports/[id]/index.vue";

const SESSION = "6f1c2a3e-0000-4000-8000-000000000001";

// useFetch answers synchronously with the analysis the API would send, run through the page's transform.
const vue = await vi.hoisted(() => import("vue"));
const state = vi.hoisted(() => ({ analysis: null, error: null }));
mockNuxtImport("useFetch", () => (_url, opts) => ({
  data: vue.ref(state.analysis ? opts.transform(state.analysis()) : null),
  error: vue.ref(state.error),
  pending: vue.ref(false),
}));
mockNuxtImport("useRoute", () => () => ({
  params: { id: SESSION },
  query: { title: "Geography" },
}));

const pdf = vi.hoisted(() => ({ build: vi.fn(), save: vi.fn() }));
mockNuxtImport("buildReportPdf", () => pdf.build);

const analysis = () => ({
  data: [
    {
      question: "Capital of France?",
      question_type: "single answer",
      options: { 1: "Paris", 2: "Lyon" },
      correct_answer: [1],
      selected_answers: { ann: [1], bob: [2] },
      avg_response_time: 2000,
    },
  ],
});
const analytics = [
  {
    username: "ann",
    firstname: "Ann",
    is_attend: true,
    question_type: "single answer",
    selected_answer: { String: "[1]" },
    correct_answer: "[1]",
    calculated_score: 900,
  },
  {
    username: "bob",
    firstname: "Bob",
    is_attend: true,
    question_type: "single answer",
    selected_answer: { String: "[2]" },
    correct_answer: "[1]",
    calculated_score: 0,
  },
];
const ranks = [
  { username: "ann", rank: 1, score: 900, response_time: 1, img_key: "" },
  { username: "bob", rank: 2, score: 0, response_time: 2, img_key: "" },
];

const mountPage = async () => {
  vi.stubGlobal(
    "fetch",
    vi.fn(async (url) => ({
      ok: true,
      json: async () => ({
        data: String(url).includes("final_score") ? ranks : analytics,
      }),
    }))
  );
  const Parent = defineComponent({
    render: () => h(Suspense, null, { default: () => h(ReportPage) }),
  });
  const wrapper = mount(Parent, {
    global: {
      mocks: { $route: { query: { title: "Geography" } } },
      stubs: {
        NuxtLink: true,
        NuxtImg: true,
        PageLayout: true,
        QuizUserAnalyticsSpace: true,
        QuizQuestionDetailModal: true,
        CodeBlockComponent: true,
      },
    },
  });
  await flushPromises();
  return wrapper;
};

const downloadButton = (wrapper) =>
  wrapper.find('button[aria-label="Download PDF Report"]');

beforeEach(() => {
  state.error = null;
  pdf.build.mockReset();
  pdf.save.mockReset();
  pdf.build.mockResolvedValue({ save: pdf.save });
});

describe("Report page: Download PDF", () => {
  it("builds the PDF from the loaded report and saves it under the quiz title", async () => {
    state.analysis = analysis;
    const wrapper = await mountPage();
    const button = downloadButton(wrapper);
    expect(button.attributes("disabled")).toBeUndefined();

    await button.trigger("click");
    await flushPromises();

    expect(pdf.build).toHaveBeenCalledTimes(1);
    const input = pdf.build.mock.calls[0][0];
    expect(input.title).toBe("Geography");
    expect(input.questions[0].question).toBe("Capital of France?");
    expect(input.participants).toEqual([
      {
        name: "Ann",
        rank: 1,
        score: 900,
        correct: 1,
        incorrect: 0,
        skipped: 0,
      },
      { name: "Bob", rank: 2, score: 0, correct: 0, incorrect: 1, skipped: 0 },
    ]);
    expect(pdf.save).toHaveBeenCalledWith(
      expect.stringMatching(/^geography-report-\d{4}-\d{2}-\d{2}\.pdf$/)
    );
  });

  // GitHub #44: the button shows progress and takes no second click while the PDF builds.
  it("shows a busy state while the PDF builds", async () => {
    state.analysis = analysis;
    let finish;
    pdf.build.mockReturnValue(
      new Promise((resolve) => (finish = () => resolve({ save: pdf.save })))
    );
    const wrapper = await mountPage();
    await downloadButton(wrapper).trigger("click");
    await flushPromises();

    const busy = downloadButton(wrapper);
    expect(busy.text()).toContain("Preparing PDF");
    expect(busy.attributes("disabled")).toBeDefined();
    expect(busy.attributes("aria-busy")).toBe("true");
    await busy.trigger("click");
    expect(pdf.build).toHaveBeenCalledTimes(1);

    finish();
    await flushPromises();
    expect(pdf.save).toHaveBeenCalledTimes(1);
    expect(downloadButton(wrapper).text()).toContain("Download PDF");
    expect(downloadButton(wrapper).attributes("disabled")).toBeUndefined();
  });

  it("is disabled with a note when the report has no data", async () => {
    state.analysis = () => ({ data: [] });
    const wrapper = await mountPage();
    expect(downloadButton(wrapper).attributes("disabled")).toBeDefined();
    expect(wrapper.text()).toContain("No report data to download yet.");
    await downloadButton(wrapper).trigger("click");
    expect(pdf.build).not.toHaveBeenCalled();
  });

  // A report the user may not read shows only the page's error, not the "no data" note.
  it("shows no download note when the report fails to load", async () => {
    state.analysis = null;
    state.error = { statusCode: 403, data: { code: 403 } };
    const wrapper = await mountPage();
    expect(downloadButton(wrapper).attributes("disabled")).toBeDefined();
    expect(wrapper.text()).not.toContain("No report data to download yet.");
  });
});

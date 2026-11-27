import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import QuizList from "~/pages/admin/quiz/list-quiz/index.vue";

// Every list the page fetches answers with the quizzes this test sets.
const vue = await vi.hoisted(() => import("vue"));
const own = vi.hoisted(() => ({ list: [] }));
mockNuxtImport("useFetch", () => (url) => ({
  data: vue.ref({ data: String(url).endsWith("/quizzes") ? own.list : [] }),
  error: vue.ref(null),
  pending: vue.ref(false),
  status: vue.ref("success"),
  refresh: () => {},
}));
vi.mock("notivue", () => ({
  usePush: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn() }),
}));

const mountPage = async () => {
  const Parent = defineComponent({
    render: () => h(Suspense, null, { default: () => h(QuizList) }),
  });
  const wrapper = mount(Parent, {
    global: {
      stubs: {
        NuxtLink: { props: ["to"], template: '<a :href="to"><slot /></a>' },
        NuxtImg: true,
        AdminQuizListCard: true,
        GenerateWithAiModal: true,
        ShareQuizModal: true,
        ClientOnly: { template: "<div><slot /></div>" },
      },
    },
  });
  await flushPromises();
  return wrapper;
};

// A host with no quizzes yet is new: point them at a ready-made quiz, not at
// AI generation, which needs an AI provider key first.
describe("quiz list for a host with no quizzes", () => {
  it("offers a ready-made quiz and no AI button", async () => {
    own.list = [];
    const wrapper = await mountPage();
    expect(wrapper.text()).not.toContain("Generate with AI");
    const link = wrapper
      .findAll("a")
      .find((a) => a.text().includes("Host a ready-made quiz"));
    expect(link?.attributes("href")).toBe("/#public-quiz");
  });

  it("keeps the AI button for a host with a quiz", async () => {
    own.list = [
      { id: "q1", title: "Trivia", created_at: "2026-11-01T00:00:00Z" },
    ];
    const wrapper = await mountPage();
    expect(wrapper.text()).toContain("Generate with AI");
  });
});

import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, beforeAll, beforeEach, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import CombinedReport from "~/pages/admin/reports/combined/[quizId].vue";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";

const route = vi.hoisted(() => ({
  path: "/admin/reports/combined/quiz-1",
  params: { quizId: "quiz-1" },
  query: {},
}));
mockNuxtImport("useRoute", () => () => route);

const navigate = vi.hoisted(() => vi.fn());
mockNuxtImport("navigateTo", () => navigate);

// Real refs keep the page reactive without network access.
const fetchHolder = vi.hoisted(() => ({ state: null, url: null, opts: null }));
mockNuxtImport("useFetch", () => async (url, opts) => {
  const { ref } = await import("vue");
  fetchHolder.url = url;
  fetchHolder.opts = opts;
  fetchHolder.state = {
    data: ref(null),
    error: ref(null),
    pending: ref(false),
  };
  return fetchHolder.state;
});

const fixture = {
  quiz_id: "quiz-1",
  title: "Space%20Quiz",
  plays: 5,
  participants: 4,
  sessions: [
    {
      id: "s1",
      activated_from: "2026-11-01T10:00:00Z",
      participants: 2,
      selected: true,
    },
    {
      id: "s2",
      activated_from: "2026-11-02T10:00:00Z",
      participants: 3,
      selected: true,
    },
    {
      id: "s3",
      activated_from: "2026-11-03T10:00:00Z",
      participants: 1,
      selected: false,
    },
  ],
  leaderboard: [
    {
      rank: 1,
      firstname: "Chen",
      username: "chen",
      session_id: "s2",
      score: 70,
      correct: 1,
      questions: 2,
      plays: 1,
    },
    {
      rank: 1,
      firstname: "Asha",
      username: "asha",
      session_id: "s2",
      score: 70,
      correct: 2,
      questions: 2,
      plays: 2,
    },
    {
      rank: 2,
      firstname: "Dana",
      username: "dana",
      session_id: "s2",
      score: 30,
      correct: 1,
      questions: 2,
      plays: 1,
    },
    {
      rank: 2,
      firstname: "Bilal",
      username: "bilal",
      session_id: "s1",
      score: 30,
      correct: 1,
      questions: 2,
      plays: 1,
    },
  ],
};

let wrapper;
let page;
const Parent = defineComponent({
  render: () => h(Suspense, null, { default: () => h(CombinedReport) }),
});

beforeAll(async () => {
  wrapper = mount(Parent, {
    global: { stubs: { NuxtLink: true, NuxtImg: true } },
  });
  await flushPromises();
  page = wrapper.findComponent(CombinedReport);
});

beforeEach(async () => {
  navigate.mockClear();
  fetchHolder.state.data.value = { data: structuredClone(fixture) };
  await flushPromises();
});

describe("Combined report page", () => {
  it("requests the quiz's combined report", () => {
    expect(fetchHolder.url()).toMatch(
      /\/admin\/reports\/quizzes\/quiz-1\/combined$/
    );
  });

  it("shows combined participants, sessions and the merged leaderboard", () => {
    const text = wrapper.text();
    expect(text).toContain("Space Quiz");
    expect(text).toContain("5 plays across 2 of 3 sessions");
    expect(page.vm.stats).toEqual([
      { label: "Participants", value: 4 },
      { label: "Sessions", value: 2 },
      { label: "Accuracy", value: "63%" },
    ]);
    const table = wrapper.findComponent(ScoreBoardTable);
    expect(table.props("scoreboardData").map((r) => r.firstname)).toEqual([
      "Chen",
      "Asha",
      "Dana",
      "Bilal",
    ]);
    // Top scorers get the rank badges; the merged board shares rank 1.
    expect(table.findAll('[aria-label="First place"]')).toHaveLength(2);
    expect(text).toContain("best of 2 plays");
  });

  it("ticks the selected sessions and regenerates with the new choice", async () => {
    const boxes = wrapper.findAll('input[type="checkbox"]');
    expect(boxes.map((b) => b.element.checked)).toEqual([true, true, false]);
    await boxes[2].setValue(true);
    await wrapper.find("form").trigger("submit");
    expect(navigate).toHaveBeenCalledWith({
      path: route.path,
      query: { sessions: "s1,s2,s3" },
    });
  });

  it("needs two sessions before it generates", async () => {
    const boxes = wrapper.findAll('input[type="checkbox"]');
    await boxes[1].setValue(false);
    expect(wrapper.find('button[type="submit"]').element.disabled).toBe(true);
  });

  it("selects and clears every session with buttons that stay visible", async () => {
    const button = (label) =>
      wrapper.findAll('button[type="button"]').find((b) => b.text() === label);
    const boxes = () =>
      wrapper.findAll('input[type="checkbox"]').map((b) => b.element.checked);
    const submit = () => wrapper.find('button[type="submit"]').element;

    await button("Select all").trigger("click");
    expect(boxes()).toEqual([true, true, true]);
    expect(submit().disabled).toBe(false);
    // Both buttons stay on the page when every session is ticked.
    expect(button("Clear all").exists()).toBe(true);

    await button("Clear all").trigger("click");
    expect(boxes()).toEqual([false, false, false]);
    expect(submit().disabled).toBe(true);
    expect(button("Select all").exists()).toBe(true);
  });

  it("scrolls a long session list inside its own box", () => {
    const list = wrapper.find('[data-testid="session-list"]');
    expect(list.classes()).toEqual(
      expect.arrayContaining(["max-h-72", "overflow-y-auto"])
    );
    expect(list.findAll('input[type="checkbox"]')).toHaveLength(3);
  });

  it("shows the API's reason when the report fails", async () => {
    fetchHolder.state.data.value = null;
    fetchHolder.state.error.value = {
      data: {
        status: "fail",
        data: "no ended sessions with players for this quiz",
      },
    };
    await flushPromises();
    expect(wrapper.find('[role="alert"]').text()).toContain(
      "no ended sessions with players for this quiz"
    );
    fetchHolder.state.error.value = null;
  });
});

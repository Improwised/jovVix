import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, beforeAll, beforeEach, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import SoloResults from "~/pages/solo/results/[attempt].vue";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";

const route = vi.hoisted(() => ({ params: { attempt: "att-1" }, query: {} }));
mockNuxtImport("useRoute", () => () => route);

// The page awaits useFetch; hand it real refs so the page stays reactive
// without any network access.
const fetchHolder = vi.hoisted(() => ({ state: null }));
mockNuxtImport("useFetch", () => async () => {
  const { ref } = await import("vue");
  fetchHolder.state = {
    data: ref(null),
    error: ref(null),
    pending: ref(false),
    refresh: () => {},
  };
  return fetchHolder.state;
});

const fixture = {
  attempt: {
    quiz_id: "quiz-1",
    title: "Space Quiz",
    score: 120,
    correct: 4,
    total: 5,
  },
  players: [
    { name: "Alice", score: 150, correct: 5, rank: 1, is_you: false },
    { name: "You", score: 120, correct: 4, rank: 2, is_you: true },
    { name: "Bob", score: 90, correct: 3, rank: 3, is_you: false },
  ],
};

let wrapper;
let page;
// The page has async setup (await useFetch), so mount it under Suspense.
const Parent = defineComponent({
  render: () => h(Suspense, null, { default: () => h(SoloResults) }),
});

beforeAll(async () => {
  wrapper = mount(Parent, {
    global: {
      stubs: {
        ClientOnly: { template: "<div><slot /></div>" },
        NuxtImg: true,
        NuxtLink: true,
        NavigationLink: true,
        ScoreBoardTable: true,
        QuizStatisticsBadges: true,
      },
    },
  });
  await flushPromises();
  page = wrapper.findComponent(SoloResults);
});

beforeEach(async () => {
  fetchHolder.state.data.value = { data: fixture };
  await flushPromises();
});

describe("Solo results page", () => {
  it("builds ScoreBoardTable rows with firstname and a unique you marker", () => {
    const rows = page.vm.leaderboardRows;
    expect(rows.map((r) => r.firstname)).toEqual(["Alice", "You", "Bob"]);
    const you = rows.find((r) => r.is_you);
    const others = rows.filter((r) => !r.is_you);
    expect(you.username).toBe(page.vm.YOU);
    expect(others.map((r) => r.username)).toEqual(["Alice", "Bob"]);
  });

  it("shows top 3, an ellipsis, then your rank when you are below the top 3", async () => {
    fetchHolder.state.data.value = {
      data: {
        attempt: { quiz_id: "q", title: "T", score: 30, correct: 1, total: 10 },
        players: [
          { name: "A", score: 100, correct: 10, rank: 1, is_you: false },
          { name: "B", score: 90, correct: 9, rank: 2, is_you: false },
          { name: "C", score: 80, correct: 8, rank: 3, is_you: false },
          { name: "D", score: 70, correct: 7, rank: 4, is_you: false },
          { name: "E", score: 60, correct: 6, rank: 5, is_you: false },
          { name: "You", score: 30, correct: 1, rank: 6, is_you: true },
        ],
      },
    };
    await flushPromises();
    const rows = page.vm.leaderboardRows;
    expect(rows).toHaveLength(5);
    expect(rows.slice(0, 3).map((r) => r.firstname)).toEqual(["A", "B", "C"]);
    expect(rows[3]).toEqual({ divider: true });
    expect(rows[4].firstname).toBe("You");
    expect(rows[4].rank).toBe(6);
    // The count pill shows all players, not just the sliced rows.
    expect(wrapper.findComponent(ScoreBoardTable).props("totalPlayers")).toBe(
      6
    );
  });

  it("keeps you in the board when a rank tie leaves you outside the top 3", async () => {
    fetchHolder.state.data.value = {
      data: {
        attempt: { quiz_id: "q", title: "T", score: 90, correct: 9, total: 10 },
        players: [
          { name: "A", score: 100, correct: 10, rank: 1, is_you: false },
          { name: "B", score: 90, correct: 9, rank: 2, is_you: false },
          { name: "C", score: 90, correct: 9, rank: 2, is_you: false },
          { name: "You", score: 90, correct: 9, rank: 2, is_you: true },
        ],
      },
    };
    await flushPromises();
    const rows = page.vm.leaderboardRows;
    expect(rows).toHaveLength(5);
    expect(rows.slice(0, 3).map((r) => r.firstname)).toEqual(["A", "B", "C"]);
    expect(rows[3]).toEqual({ divider: true });
    expect(rows[4].firstname).toBe("You");
    expect(rows[4].rank).toBe(2);
  });

  it("passes the you marker as user-name so the row highlights", () => {
    const table = wrapper.findComponent(ScoreBoardTable);
    expect(table.exists()).toBe(true);
    expect(table.props("userName")).toBe(page.vm.YOU);
  });

  it("derives quiz statistics badges from the attempt", () => {
    expect(page.vm.userStatistics).toEqual({
      accuracy: 80,
      totalScore: 120,
      correctAnwers: 4,
      wrongAnwers: 1,
      unAttemptedQuestions: 0,
    });
  });

  it("counts skipped questions when the API reports total questions", async () => {
    fetchHolder.state.data.value = {
      data: {
        attempt: {
          quiz_id: "q",
          title: "T",
          score: 60,
          correct: 3,
          total: 3,
          total_questions: 6,
        },
        players: [
          { name: "You", score: 60, correct: 3, rank: 1, is_you: true },
        ],
      },
    };
    await flushPromises();
    expect(page.vm.userStatistics).toEqual({
      accuracy: 50,
      totalScore: 60,
      correctAnwers: 3,
      wrongAnwers: 0,
      unAttemptedQuestions: 3,
    });
    expect(wrapper.text()).toContain("3 / 6 correct");
  });

  it("renders hero score, rank badge and results copy", () => {
    expect(wrapper.text()).toContain("Great game!");
    expect(wrapper.text()).toContain("Solo complete · Space Quiz");
    expect(wrapper.text()).toContain("2nd Place");
    expect(wrapper.text()).toContain("4 / 5 correct");
    expect(wrapper.text()).toContain("Personal best: 120 points");
  });

  it("shows this attempt's score in the hero, keeping personal best separate", async () => {
    fetchHolder.state.data.value = {
      data: {
        attempt: { quiz_id: "q", title: "T", score: 50, correct: 2, total: 5 },
        players: [
          { name: "A", score: 200, correct: 5, rank: 1, is_you: false },
          { name: "You", score: 120, correct: 4, rank: 2, is_you: true },
          { name: "B", score: 90, correct: 3, rank: 3, is_you: false },
        ],
      },
    };
    await flushPromises();
    const hero = wrapper.find("header");
    expect(hero.text()).toContain("Score");
    expect(hero.text()).toContain("2nd Place");
    // This attempt scored 50 (not the best 120).
    expect(hero.text()).toContain("50");
    expect(hero.text()).toContain("Personal best: 120 points");
  });
});

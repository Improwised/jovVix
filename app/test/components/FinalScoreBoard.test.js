import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { mount } from "@vue/test-utils";
import { describe, it, expect, beforeAll, vi } from "vitest";
import FinalScoreBoard from "~/components/FinalScoreBoard.vue";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";

const route = vi.hoisted(() => ({
  query: { aqi: "quiz123", winner_ui: "false" },
}));
mockNuxtImport("useRoute", () => () => route);

const props = {
  userURL: "/scoreboard",
  isAdmin: true,
};
const mountComponent = () =>
  mount(FinalScoreBoard, {
    props,
    global: {
      stubs: {
        NuxtImg: true,
        NuxtLink: true,
        ScoreBoardTable: true,
        QuizStatisticsBadges: true,
        QuizAnalysis: true,
        QuizAnalysisLoginGate: true,
      },
    },
  });

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time.
let wrapper;
beforeAll(() => {
  wrapper = mountComponent();
});

describe("FinalScoreBoard Test", () => {
  beforeAll(async () => {
    wrapper.vm.requestPending = false;
    wrapper.vm.scoreboardData = [
      { id: 1, firstname: "User 1", username: "User 1", rank: 1, score: 1000 },
      { id: 2, firstname: "User 2", username: "User 2", rank: 2, score: 900 },
    ];

    await wrapper.vm.$nextTick();
  });

  it("fetches scoreboard data for admin", async () => {
    expect(wrapper.vm.scoreboardData).toHaveLength(2);
  });

  it("renders ScoreBoardTable component for admin", async () => {
    expect(wrapper.findComponent(ScoreBoardTable).exists()).toBe(true);
  });

  it("shows player names without account usernames on the podium", async () => {
    route.query.winner_ui = "true";
    const podium = mount(FinalScoreBoard, {
      props,
      global: {
        stubs: {
          ClientOnly: { template: "<div><slot /></div>" },
          NuxtImg: true,
          NuxtLink: true,
          WinnerConfetti: true,
        },
      },
    });
    podium.vm.requestPending = false;
    podium.vm.scoreboardData = [
      { firstname: "Smit", username: "smit_one" },
      { firstname: "Smit", username: "smit_two" },
      { firstname: "Alex" },
    ];
    await podium.vm.$nextTick();
    expect(podium.text()).not.toContain("@smit_one");
    expect(podium.text()).not.toContain("@smit_two");
    expect(podium.text()).toContain("Alex");
    expect(podium.text()).not.toContain("@undefined");
    podium.unmount();
    route.query.winner_ui = "false";
  });
});

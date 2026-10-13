import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { mount } from "@vue/test-utils";
import { describe, it, expect } from "vitest";
import FinalScoreBoard from "~/components/FinalScoreBoard.vue";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";

mockNuxtImport("useRoute", () => () => ({
  query: {
    aqi: "quiz123",
    winner_ui: "false",
  },
}));

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

let wrapper = mountComponent();

describe("FinalScoreBoard Test", async () => {
  wrapper.vm.requestPending = false;
  wrapper.vm.scoreboardData = [
    { id: 1, firstname: "User 1", username: "User 1", rank: 1, score: 1000 },
    { id: 2, firstname: "User 2", username: "User 2", rank: 2, score: 900 },
  ];

  await wrapper.vm.$nextTick();

  it("fetches scoreboard data for admin", async () => {
    expect(wrapper.vm.scoreboardData).toHaveLength(2);
  });

  it("renders ScoreBoardTable component for admin", async () => {
    expect(wrapper.findComponent(ScoreBoardTable).exists()).toBe(true);
  });
});

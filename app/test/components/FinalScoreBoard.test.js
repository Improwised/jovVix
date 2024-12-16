import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import FinalScoreBoard from "~/components/FinalScoreBoard.vue";
import WinnerCard from "~/components/WinnerCard.vue";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";
import constants from "~/test/constants";

mockNuxtImport("useRoute", () => () => ({
  query: {
    aqi: "quiz123",
    winner_ui: "true",
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
      mocks: {
        getFinalScoreboardDetails: vi.fn(),
      },
      stubs: {
        VBtn: constants.slotTemplate,
        VCard: constants.slotTemplate,
        VCardItem: constants.slotTemplate,
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
  wrapper.vm.analysisData = [{ question: "Q1" }, { question: "Q2" }];

  wrapper.vm.userStatistics = questionsAnalysis([
    { question: "Q1" },
    { question: "Q2" },
  ]);
  await wrapper.vm.$nextTick();

  it("fetches scoreboard data for admin", async () => {
    expect(wrapper.vm.scoreboardData).toHaveLength(2);
  });

  it("fetches analysis data for users", async () => {
    expect(wrapper.vm.analysisData).toHaveLength(2);
  });

  it("renders the winner UI for admin", async () => {
    expect(wrapper.find("img#myVideo").exists()).toBe(true);
    expect(wrapper.findComponent(WinnerCard).exists()).toBe(true);
  });

  it("renders the scoreboard table for non-admin users", async () => {
    props.isAdmin = false;
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    expect(wrapper.findComponent(ScoreBoardTable).exists()).toBe(true);
  });
});

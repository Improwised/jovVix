import { describe, it, expect, beforeAll } from "vitest";
import { mount } from "@vue/test-utils";
import ScoreSpace from "~/components/Quiz/ScoreSpace.vue";

const props = {
  data: {
    status: "success",
    data: {
      duration: 20,
      options: {
        1: { value: "Blue", isAnswer: false },
        2: { value: "Red", isAnswer: true },
        3: { value: "Green", isAnswer: false },
        4: { value: "Pink", isAnswer: false },
      },
      options_media: "text",
      question: "What is the color of the Strawberry?",
      question_media: "text",
      quiz_id: "1",
      rankList: [
        {
          rank: 1,
          points: 0,
          score: 0,
          response_time: 3390,
          username: "john",
          firstname: "doe",
          img_key: "Eden",
          streak_count: 0,
        },
      ],
      resource: "",
      totalQuestions: 5,
      userResponses: [
        {
          id: "user1",
          answers: {
            String: "[3]",
            Valid: true,
          },
        },
      ],
    },
    event: "show_score",
    action: "show score page during quiz",
    component: "Score",
  },
  isAdmin: true,
  userName: "",
  selectedAnswer: 0,
  analysisTab: "ranking",
  quizState: "running",
};

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time.
let wrapper;
beforeAll(() => {
  wrapper = mount(ScoreSpace, {
    props,
    global: {
      stubs: {
        NuxtImg: true,
        NuxtLink: true,
        CodeBlockComponent: true,
      },
    },
  });
});

describe("ScoreSpace test", () => {
  it("renders correctly with default props", () => {
    expect(wrapper.find("h2").exists()).toBe(true);
    expect(wrapper.find("h3").exists()).toBe(true);
  });

  it("renders the correct answer and selected incorrect answer", async () => {
    expect(wrapper.text()).toContain("Red");
    expect(wrapper.text()).toContain("Blue");

    props.selectedAnswer = 3;
    props.isAdmin = false;

    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    expect(wrapper.text()).toContain("Green");
  });

  it("emits 'askSkipTimer' when the skip button is clicked", async () => {
    props.isAdmin = true;
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    const buttons = wrapper.findAll("button");
    const skipButton = buttons.find(
      (b) => b.text().includes("Skip") || b.text().includes("Finish")
    );
    expect(skipButton).toBeTruthy();
    await skipButton.trigger("click");
    expect(wrapper.emitted("askSkipTimer")).toBeTruthy();
  });

  it("changes the analysis tab when a tab is clicked", async () => {
    const chartButtons = wrapper.findAll("button");
    const chartTab = chartButtons.find((b) => b.text().includes("Chart"));
    if (chartTab) {
      await chartTab.trigger("click");
      expect(wrapper.emitted("changeAnalysisTab")[0]).toEqual(["chart"]);
    }
  });

  it("renders the rank list correctly", () => {
    expect(wrapper.text()).toContain("1");
    expect(wrapper.text()).toContain("doe");
    expect(wrapper.text()).toContain("0");
  });
});

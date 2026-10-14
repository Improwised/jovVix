import { mountSuspended } from "@nuxt/test-utils/runtime";
import { it, expect, describe } from "vitest";
import ScoreSpace from "~~/components/Quiz/ScoreSpace.vue";

const data = {
  status: "success",
  event: "show_score",
  component: "Score",
  data: {
    question_no: 1,
    totalQuestions: 2,
    duration: 20,
    question: "Which of these is a colour",
    question_media: "text",
    options_media: "text",
    options: {
      1: { value: "Table", isAnswer: false },
      2: { value: "Red", isAnswer: true },
    },
    rankList: [],
  },
};

const press = (key) =>
  window.dispatchEvent(new KeyboardEvent("keydown", { key, bubbles: true }));

describe("ScoreSpace keyboard shortcuts", () => {
  it("skips to the next question on Enter for the host, once", async () => {
    const wrapper = await mountSuspended(ScoreSpace, {
      props: { data, isAdmin: true, quizState: "running" },
    });

    press("Enter");
    press("Enter");

    expect(wrapper.emitted("askSkipTimer")).toHaveLength(1);
    wrapper.unmount();
  });

  it("ignores Enter for players", async () => {
    const wrapper = await mountSuspended(ScoreSpace, {
      props: { data, isAdmin: false, quizState: "running" },
    });

    press("Enter");

    expect(wrapper.emitted("askSkipTimer")).toBeUndefined();
    wrapper.unmount();
  });
});

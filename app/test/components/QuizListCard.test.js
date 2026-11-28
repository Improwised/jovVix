import { describe, it, expect, beforeAll } from "vitest";
import { mountSuspended } from "@nuxt/test-utils/runtime";
import QuizListCard from "~/components/QuizListCard.vue";
import UtilsStartQuiz from "~/components/utils/StartQuiz.vue";

const props = {
  details: {
    title: "Sample Quiz",
    description: { String: "This is a sample description." },
    total_questions: 10,
    created_at: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
    id: "quiz1",
  },
  isPlayedQuiz: false,
};

const mountComponent = () => mountSuspended(QuizListCard, { props });

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time. Nuxt 4.7 renders NuxtLink inline, so
// links are asserted on their href rather than through RouterLinkStub.
let wrapper;
beforeAll(async () => {
  wrapper = await mountComponent();
});

describe("QuizListCard test", () => {
  it("renders correctly with props", () => {
    // Check if title, description, and questions are rendered
    expect(wrapper.html()).toContain("Sample Quiz");
    expect(wrapper.html()).toContain("This is a sample description.");
    expect(wrapper.html()).toContain("10 Questions");
  });

  it("renders the `UtilsStartQuiz` component when `isPlayedQuiz` is false", () => {
    const startQuizComponent = wrapper.findComponent(UtilsStartQuiz);
    expect(startQuizComponent.exists()).toBe(true);
    expect(startQuizComponent.props("quizId")).toBe("quiz1");
  });

  it("renders the correct `NuxtLink` for played quizzes when `isPlayedQuiz` is true", async () => {
    props.isPlayedQuiz = true;
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();
    const viewQuizLink = wrapper.find('a[href="/admin/played_quiz/quiz1"]');
    expect(viewQuizLink.exists()).toBe(true);
  });

  it("renders the correct `NuxtLink` for unplayed quizzes when `isPlayedQuiz` is false", async () => {
    props.isPlayedQuiz = false;
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();
    const viewQuizLink = wrapper.find('a[href="/admin/quiz/list-quiz/quiz1"]');
    expect(viewQuizLink.exists()).toBe(true);
  });

  it("displays the correct creation time", () => {
    expect(wrapper.text()).toContain("2 days ago");
  });
});

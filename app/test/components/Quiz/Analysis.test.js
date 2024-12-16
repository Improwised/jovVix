import { mountSuspended, renderSuspended } from "@nuxt/test-utils/runtime";
import { it, expect, describe } from "vitest";
import Analysis from "~~/components/Quiz/Analysis.vue";
import QuizQuestionAnalysis from "~~/components/Quiz/QuestionAnalysis.vue";
import QuizOptionsAnalysis from "~~/components/Quiz/OptionsAnalysis.vue";
import { mount } from "@vue/test-utils";

const mountComponent = async (props) => {
  return await mountSuspended(Analysis, {
    props,
    global: {
      stubs: {
        FontAwesomeIcon: true,
        VProgressCircular: true,
      },
    },
  });
};

const data = [
  {
    username: "sd",
    firstname: "sd",
    selected_answer: { String: "[2]", Valid: true },
    correct_answer: "[2]",
    calculated_score: 973,
    is_attend: true,
    response_time: 1872,
    calculated_points: 1,
    question: "what is the output of this code?",
    raw_options: "eyIxIjoiMTAiLCIyIjoiMTIwIiwiMyI6IjI0MCIsIjQiOiIyMCJ9",
    options: { 1: "10", 2: "120", 3: "240", 4: "20" },
    question_media: "code",
    options_media: "text",
    resource:
      "function factorial(n) { \n let ans = 1; \n \n if(n === 0)\n return 1;\n for (let i = 2; i <= n; i++) \n ans = ans * i; \n return ans; \n}\nconsole.log(factorial(5));",
    points: 1,
    question_type_id: 1,
    question_type: "single answer",
    order_no: 1,
  },
  {
    username: "sd",
    firstname: "sd",
    selected_answer: { String: "[4]", Valid: true },
    correct_answer: "[3]",
    calculated_score: 0,
    is_attend: true,
    response_time: 1040,
    calculated_points: 0,
    question: "Which company is known for its Think Different slogan?",
    raw_options:
      "img",
    options: {
      1: "img1",
      2: "img2",
      3: "img3",
      4: "img4",
    },
    question_media: "text",
    options_media: "image",
    resource: "",
    points: 1,
    question_type_id: 1,
    question_type: "single answer",
    order_no: 2,
  },
  {
    username: "sd",
    firstname: "sd",
    selected_answer: { String: "[3]", Valid: true },
    correct_answer: "[1,2,3,4,5]",
    calculated_score: 987,
    is_attend: true,
    response_time: 1461,
    calculated_points: 1,
    question:
      "What features of a travel destination are most important to you when planning a vacation?",
    raw_options:
      "img",
    options: {
      1: "Scenic landscapes",
      2: "Cultural experiences",
      3: "Local cuisine",
      4: "Historical sites",
      5: "Adventure activities",
    },
    question_media: "text",
    options_media: "text",
    resource: "",
    points: 1,
    question_type_id: 2,
    question_type: "survey",
    order_no: 3,
  },
];

describe("Analysis Test", () => {
  it("can mount inner components", async () => {
    const wrapper = await mountComponent({ data });

    const questionAnalysisComponents =
      wrapper.findAllComponents(QuizQuestionAnalysis);
    expect(questionAnalysisComponents).toHaveLength(3);

    const optionAnalysisComponents =
      wrapper.findAllComponents(QuizOptionsAnalysis);
    expect(optionAnalysisComponents).toHaveLength(3);
  });

  it("check computed property", async () => {
    const wrapper = await mountComponent({ data });
    expect(wrapper.vm.questionsAnalysis.length).toBe(3);
  });
});

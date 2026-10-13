import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import OptionsAnalysis from "~/components/Quiz/OptionsAnalysis.vue";

const wrapper = mount(OptionsAnalysis, {
  props: {
    options: { 1: "Paris", 2: "Rome", 3: "Athens", 4: "Cairo" },
    selectedAnswers: { 3: ["lol_yJc3"] },
    correctAnswer: [2],
    optionsMedia: "text",
    isAdminAnalysis: true,
  },
  global: {
    stubs: {
      Option: {
        template: '<div class="option-item"><slot /></div>',
        props: [
          "order",
          "option",
          "selected",
          "isCorrect",
          "isPicked",
          "optionsMedia",
          "isAdminAnalysis",
        ],
      },
    },
  },
});

describe("OptionsAnalysis test", () => {
  it("renders all options", () => {
    const optionsEls = wrapper.findAll('[class*="option-item"]');
    expect(optionsEls.length).toBe(4);
  });

  it("applies correct styles for correct answers", () => {
    const optionWrappers = wrapper.findAll("div");
    const wrapperEls = optionWrappers.filter((w) =>
      w.classes().includes("bg-jv-mint")
    );
    expect(wrapperEls.length).toBe(1);
  });

  it("applies correct styles for wrong selected answers", async () => {
    let wrongOptions = wrapper.findAll('[class*="bg-jv-salmon"]');
    expect(wrongOptions.length).toBe(0);

    await wrapper.setProps({
      options: { 1: "10", 2: "120", 3: "240", 4: "20" },
      correctAnswer: "[2]",
      selectedAnswer: "[3]",
      selectedAnswers: {},
      optionsMedia: "text",
      isAdminAnalysis: false,
    });
    wrongOptions = wrapper.findAll('[class*="bg-jv-salmon"]');
    expect(wrongOptions.length).toBe(1);
  });

  it("handles empty props correctly", async () => {
    await wrapper.setProps({
      options: {},
      correctAnswer: [],
      selectedAnswer: "",
      selectedAnswers: {},
      optionsMedia: "",
      isAdminAnalysis: false,
    });

    const optionWrappers = wrapper.findAll("div");
    const wrapperEls = optionWrappers.filter((w) =>
      w.classes().includes("shadow-brutal-sm")
    );
    expect(wrapperEls.length).toBe(0);
  });
});

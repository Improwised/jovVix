import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import OptionsAnalysis from "~/components/Quiz/OptionsAnalysis.vue";
import constants from "~/test/constants";

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
        template: constants.slotTemplate,
      },
    },
  },
});

describe("OptionsAnalysis test", () => {
  it("renders all options", () => {
    const optionBoxes = wrapper.findAll(".option-box");
    expect(optionBoxes.length).toBe(4);
  });

  it("applies correct styles for correct answers", () => {
    const correctOptions = wrapper.findAll(".bg-light-success");
    expect(correctOptions.length).toBe(1); // Answers 1 and 2 are correct
  });

  it("applies correct styles for wrong selected answers", async () => {
    let wrongOptions = wrapper.findAll(".bg-light-danger");
    expect(wrongOptions.length).toBe(0);

    await wrapper.setProps({
      options: { 1: "10", 2: "120", 3: "240", 4: "20" },
      correctAnswer: "[2]",
      selectedAnswer: "[3]",
      selectedAnswers: {},
      optionsMedia: "text",
      isAdminAnalysis: false,
    });
    wrongOptions = wrapper.findAll(".bg-light-danger");
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

    const optionBoxes = wrapper.findAll(".option-box");
    expect(optionBoxes.length).toBe(0); // No options provided
  });
});

import { mount } from "@vue/test-utils";
import { describe, it, expect } from "vitest";
import EditQuestion from "~/components/Quiz/EditQuestion.vue";
import constants from "~/test/constants";

const props = {
  question: {
    question: "What is the capital of France?",
    question_media: "text",
    options: { 0: "Paris", 1: "Berlin", 2: "Madrid", 3: "Rome" },
    correct_answer: "[0]",
    question_id: "123",
    question_type_id: 1,
    options_media: "text",
  },
  quizId: "quiz-123",
  questionId: "q-123",
};

const mountComponent = () => {
  return mount(EditQuestion, {
    props,
    global: {
      stubs: {
        VFileInput: constants.slotTemplate,
      },
    },
  });
};
let wrapper = mountComponent();

describe("EditQuestion.vue", () => {
  it("renders the component with the provided props", () => {
    expect(wrapper.find("input").exists()).toBe(true);
    expect(wrapper.find("input").element.value).toBe(props.question.question);
    expect(wrapper.findAll(".option-box").length).toBe(
      Object.keys(props.question.options).length
    );
  });

  it("updates question text on input", async () => {
    const input = wrapper.find("input");
    await input.setValue("What is the largest planet?");

    expect(wrapper.vm.editableQuestion.question).toBe(
      "What is the largest planet?"
    );
  });

  it('renders image input when question_media is "image"', async () => {
    props.question.question_media = "image";
    wrapper = mountComponent();
    expect(wrapper.find("#image-attachment-question").exists()).toBe(true);
  });

  it('renders code editor when question_media is "code"', async () => {
    props.question.question_media = "code";
    wrapper = mountComponent();
    expect(wrapper.findComponent({ name: "CodeBlockComponent" }).exists()).toBe(
      true
    );
  });

  it("updates correct answer on radio button change", async () => {
    const radioButton = wrapper.find('input[type="radio"][value="1"]');
    await radioButton.setChecked();

    expect(wrapper.vm.picked).toBe("1");
  });
});

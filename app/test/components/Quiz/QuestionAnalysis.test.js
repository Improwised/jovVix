import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import QuestionAnalysis from "~/components/Quiz/QuestionAnalysis.vue";
import { VProgressCircular } from "vuetify/components";
import CodeBlockComponent from "~/components/CodeBlockComponent.vue";
import DeleteDialog from "~/components/DeleteDialog.vue";

const defaultProps = {
  question: {
    question_id: "1",
    question: "What is the color of the Strawberry?",
    type: 1,
    options: {
      1: "red",
      2: "pink",
      3: "black",
      4: "yellow",
    },
    question_media: "text",
    options_media: "text",
    resource: "",
    correct_answer: [1],
    selected_answers: { 1: ["husen_eSE1"] },
    duration: 30,
    avg_response_time: 6333,
    userParticipants: 1,
    correctPercentage: 100,
  },
  order: 1,
  isAdminAnalysis: true,
  isForQuiz: false,
  isEditable: true,
};

const mountComponent = () => {
  return mount(QuestionAnalysis, {
    props: defaultProps,
    global: {
      stubs: {
        VProgressCircular: true,
      },
    },
  });
};

let wrapper = mountComponent();

describe("QuestionDetail.vue", () => {
  it("renders the question order and text", () => {
    expect(wrapper.find(".text-primary").text()).toContain("Question: 1");
    expect(wrapper.find("h3").text()).toBe(defaultProps.question.question);
  });

  it("renders an image if question_media is 'image'", async () => {
    defaultProps.question.resource = "test-image";
    defaultProps.question.question_media = "image";
    defaultProps.question.question = "lol";

    wrapper.unmount();
    wrapper = mountComponent();

    await wrapper.setProps(defaultProps);
    const image = wrapper.find("img");
    expect(image.exists()).toBe(true);
    expect(image.attributes("src")).toBe(defaultProps.question.resource);
  });

  it("renders a code block if question_media is 'code'", () => {
    defaultProps.question.resource = "console.log('Hello!');";
    defaultProps.question.question_media = "code";

    wrapper.unmount();
    wrapper = mountComponent();

    const codeBlock = wrapper.findComponent(CodeBlockComponent);
    expect(codeBlock.exists()).toBe(true);
    expect(codeBlock.props("code")).toBe("console.log('Hello!');");
  });

  it("renders admin analysis section", () => {
    const avgTime = wrapper.find("span.bg-light-primary");
    expect(avgTime.text()).toContain("AVG. Response Time: 6.33/ 30 seconds");

    const mcqBadge = wrapper.find(".badge.bg-light-info");
    expect(mcqBadge.text()).toBe("M.C.Q.");

    const progressCircular = wrapper.findComponent(VProgressCircular);
    expect(progressCircular.exists()).toBe(true);
    expect(progressCircular.attributes("color")).toBe("teal");
    expect(progressCircular.attributes("modelvalue")).toContain(
      defaultProps.question.correctPercentage
    );
  });

  it("renders user response details for non-admin analysis", () => {
    defaultProps.isAdminAnalysis = false;
    defaultProps.question.response_time = 1924;
    defaultProps.question.is_attend = true;

    wrapper.unmount();
    wrapper = mountComponent();

    const responseTime = wrapper.find("span.bg-light-primary");
    expect(responseTime.text()).toContain("Response Time: 1.92 seconds");

    const attemptedBadge = wrapper.find(".badge.bg-success");
    expect(attemptedBadge.text()).toBe("Attempted");
  });

  it("renders not attempted badge when user has not attempted the question", () => {
    defaultProps.question.is_attend = false;
    wrapper.unmount();
    wrapper = mountComponent();
    const notAttemptedBadge = wrapper.find("span.bg-danger");
    expect(notAttemptedBadge.text()).toBe("Not Attempted");
  });

  it("renders edit and delete buttons for editable questions", async () => {
    defaultProps.isEditable = true;
    const editButton = wrapper.find("button[title='Edit question']");
    expect(editButton.exists()).toBe(true);

    const deleteButton = wrapper.find("button[title='Delete question']");
    expect(deleteButton.exists()).toBe(true);

    const deleteDialog = wrapper.findComponent(DeleteDialog);
    expect(deleteDialog.exists()).toBe(true);

    // Emitting editQuestion
    await editButton.trigger("click");
    expect(wrapper.emitted("editQuestion")[0]).toEqual(["1"]);

    // Emitting deleteQuestion
    deleteDialog.vm.$emit("confirm-delete");
    expect(wrapper.emitted("deleteQuestion")[0]).toEqual(["1"]);
  });

  it("does not render edit and delete buttons if not editable", () => {
    defaultProps.isEditable = false;

    wrapper.unmount();
    wrapper = mountComponent();

    const editButton = wrapper.find("button[title='Edit question']");
    const deleteButton = wrapper.find("button[title='Delete question']");

    expect(editButton.exists()).toBe(false);
    expect(deleteButton.exists()).toBe(false);
  });

  it("handles missing question props gracefully", () => {
    defaultProps.question = {};
    wrapper.unmount();
    wrapper = mountComponent();

    expect(wrapper.find("h3").text()).toBe("");
    expect(wrapper.find("img").exists()).toBe(false);
    expect(wrapper.find("v-progress-circular").exists()).toBe(false);
  });
});

import { describe, it, expect, vi, afterEach } from "vitest";
import { mount } from "@vue/test-utils";
import QuestionSpace from "~/components/Quiz/QuestionSpace.vue";
import constants from "~/config/constants";

vi.mock("notivue", () => ({
  usePush: vi.fn(() => ({
    error: vi.fn(),
    warning: vi.fn(),
  })),
}));

const props = {
  data: {
    status: "success",
    data: {
      duration: 30,
      id: "1",
      no: 3,
      options: {
        1: "Scenic landscapes",
        2: "Cultural experiences",
        3: "Local cuisine",
        4: "Historical sites",
        5: "Adventure activities",
      },
      options_media: "text",
      question:
        "What features of a travel destination are most important to you when planning a vacation?",
      question_media: "text",
      question_time: "",
      quiz_id: "1",
      resource: "",
      totalJoinUser: 1,
      totalQuestions: 3,
    },
    event: "send_question",
    action: "send single question to user",
    component: "Question",
  },
};
const mountComponent = () => {
  return mount(QuestionSpace, {
    props,
    global: {
      stubs: {
        NuxtImg: true,
        NuxtLink: true,
        Popover: true,
        PopoverContent: true,
        PopoverTrigger: true,
        CodeBlockComponent: true,
      },
    },
  });
};
let wrapper = mountComponent();

describe("QuestionSpace test", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
  });
  it("renders correctly when a question is present", () => {
    expect(wrapper.find("h3").exists()).toBe(true);
  });

  it("renders the question and options", () => {
    const options = wrapper.findAll("button");
    expect(options.length).toBe(6);
  });

  it("emits 'sendAnswer' when an answer is selected", async () => {
    const buttons = wrapper.findAll("button");
    const optionButton = buttons.find((b) =>
      b.text().includes("Scenic landscapes")
    );
    if (optionButton) {
      await optionButton.trigger("click");
      expect(wrapper.emitted("sendAnswer")).toBeTruthy();
    }
  });

  it("handles skipping the question", async () => {
    wrapper.vm.handleSkip({ preventDefault: vi.fn() });
    expect(wrapper.emitted("askSkip")).toBeTruthy();
  });

  it("renders the countdown if no question is present", async () => {
    props.data.event = constants.Counter;
    props.data.data.count = 3;

    wrapper.unmount();
    wrapper = mountComponent();

    expect(wrapper.text()).toContain("3");
  });

  it("renders admin skip button when not last question", async () => {
    props.isAdmin = true;
    props.data.event = "send_question";
    delete props.data.data.count;
    props.data.data.no = props.data.data.totalQuestions - 1;

    wrapper.unmount();
    wrapper = mountComponent();
    expect(wrapper.text()).toContain("Skip");
  });

  it("renders admin finish button when last question", async () => {
    props.isAdmin = true;
    props.data.event = "send_question";
    delete props.data.data.count;
    props.data.data.no = props.data.data.totalQuestions;

    wrapper.unmount();
    wrapper = mountComponent();
    expect(wrapper.text()).toContain("Finish Quiz");
  });
});

import { describe, it, expect, vi, afterAll, afterEach } from "vitest";
import { mount } from "@vue/test-utils";
import { useToast } from "vue-toastification";
import { createTestingPinia } from "@pinia/testing";
import { useMusicStore } from "~/store/music";
import QuestionSpace from "~/components/Quiz/QuestionSpace.vue";
import constants from "~/config/constants";
import QuizQuestionAnalysis from "~/components/Quiz/QuestionAnalysis.vue";
import { VProgressCircular } from "vuetify/components";
import CodeBlockComponent from "~/components/CodeBlockComponent.vue";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";
vi.mock("vue-toastification", () => ({
  useToast: vi.fn(() => ({
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
        VProgressCircular: true,
        FontAwesomeIcon: true,
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
    expect(wrapper.findComponent(QuizQuestionAnalysis).exists()).toBe(true);

    const progressCircular = wrapper.findComponent(VProgressCircular);
    expect(progressCircular.exists()).toBe(true);
  });

  it("renders the question and options", () => {
    const options = wrapper.findAll(".option-box");
    expect(options.length).toBe(5);

    options.forEach((option, index) => {
      expect(option.text()).toContain(props.data.data.options[index + 1]);
    });
  });

  it("emits 'sendAnswer' when an answer is selected", async () => {
    const inputs = wrapper.findAll("input[type='radio']");
    await inputs[0].setValue();
    await inputs[0].trigger("change");

    expect(wrapper.emitted("sendAnswer")[0]).toStrictEqual([[1]]);
  });

  it("disables options after submission", async () => {
    wrapper.vm.isSubmitted = true;
    await wrapper.vm.$nextTick();

    const inputs = wrapper.findAll("input[type='radio']");
    inputs.forEach((input) => {
      expect(input.attributes("disabled")).toBeDefined();
    });
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

    const countdown = wrapper.find(".d-flex");
    expect(countdown.exists()).toBe(true);
    expect(countdown.text()).toContain("3");
  });

  it("renders admin skip button", async () => {
    props.isAdmin = true;
    props.data.event = constants.GetQuestion;
    delete props.data.data.count;

    wrapper.unmount();
    wrapper = mountComponent();
    const skipButton = wrapper.find("button.btn-primary");
    expect(wrapper.find("button.btn-primary").exists()).toBe(true);
    expect(skipButton.text()).contain("skip");
  });
});

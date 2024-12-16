import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import ScoreSpace from "~/components/Quiz/ScoreSpace.vue";
import { VProgressLinear } from "vuetify/components";
import { FontAwesomeIcon } from "@fortawesome/vue-fontawesome";

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

let wrapper = mount(ScoreSpace, {
  props,
  global: {
    stubs: {
      VProgressLinear: true,
      FontAwesomeIcon: true,
    },
  },
});

describe("ScoreSpace test", () => {
  it("renders correctly with default props", () => {
    expect(wrapper.find("h1").exists()).toBe(true);
    expect(wrapper.findComponent(VProgressLinear).exists()).toBe(true);
  });

  it("renders the correct answer and selected incorrect answer", async () => {
    expect(wrapper.findAll(".option-box").length).toBe(4);

    const correctOption = wrapper.find(".bg-light-success");

    expect(correctOption.exists()).toBe(true);
    expect(correctOption.text()).toContain(props.data.data.options[2].value);

    props.selectedAnswer = 3;
    props.isAdmin = false;

    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    const incorrectOption = wrapper.find(".bg-light-danger");

    expect(incorrectOption.exists()).toBe(true);
    expect(incorrectOption.text()).toContain(
      props.data.data.options[props.selectedAnswer].value
    );
  });

  it("emits 'askSkipTimer' when the skip button is clicked", async () => {
    props.isAdmin = true;
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    const skipButton = wrapper.find("button.btn-primary");
    await skipButton.trigger("click");
    expect(wrapper.emitted("askSkipTimer")).toBeTruthy();
    expect(skipButton.attributes("disabled")).toBeDefined();
  });

  it("changes the analysis tab when a tab is clicked", async () => {
    const chartTab = wrapper.find("#pills-chart-tab");
    await chartTab.trigger("click");

    expect(wrapper.emitted("changeAnalysisTab")[0]).toEqual(["chart"]);
  });

  it("renders the rank list correctly", () => {
    const rows = wrapper.findAll("table tbody tr");
    expect(rows.length).toBe(1);

    const firstRow = rows[0];
    expect(firstRow.text()).toContain("1");
    expect(firstRow.text()).toContain("doe");
    expect(firstRow.text()).toContain("0");
  });
});

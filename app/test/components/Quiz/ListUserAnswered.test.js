import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import ListUserAnswerd from "~/components/Quiz/ListUserAnswered.vue";
import { createTestingPinia } from "@pinia/testing";
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";
import constants from "~/test/constants";

vi.mock("vue-toastification", () => ({
  useToast: vi.fn(() => ({
    error: vi.fn(),
  })),
}));

const pinia = createTestingPinia({
  createSpy: vi.fn,
});

const usersThatSubmittedAnswer = useUserThatSubmittedAnswer(pinia);

const mountComponent = () => {
  return mount(ListUserAnswerd, {
    props: {
      data: { status: "success", data: "" },
      runningQuizJoinUser: 10,
    },
    global: {
      plugins: [pinia],
      stubs: {
        FontAwesomeIcon: true,
        VCard: constants.slotTemplate,
      },
    },
  });
};

describe("LlistUserAnswered Test", () => {
  it("renders correctly when no answers are submitted", async () => {
    expect(usersThatSubmittedAnswer.usersSubmittedAnswers.length).toBe(0);

    const wrapper = mountComponent();

    expect(wrapper.find(".col-7").exists()).toBe(true);
    expect(wrapper.find("h5").text()).toBe("No One Answered Till Now..");
  });

  it("renders correctly when answers are submitted", async () => {
    const mockUsersSubmittedAnswers = [
      {
        UserId: 1,
        img_key: "avatar1.png",
        first_name: "Alice",
        username: "alice123",
      },
    ];

    usersThatSubmittedAnswer.usersSubmittedAnswers = mockUsersSubmittedAnswers;
    const wrapper = mountComponent();
    expect(wrapper.find(".col-6").exists()).toBe(true);
    expect(wrapper.find("h5").text()).toContain("1/10 People Answered");
  });

  it("handles new user join event correctly", async () => {
    const wrapper = mountComponent();
    await wrapper.setProps({
      data: { event: "send_question", data: { totalJoinUser: 12 } },
    });
    expect(wrapper.vm.totalUser).toBe(12);
  });

  it("emits 'autoSkip' when all users have answered", async () => {
    const mockUsersSubmittedAnswers = [
      {
        UserId: 1,
        img_key: "avatar1.png",
        first_name: "Alice",
        username: "alice123",
      },
      {
        UserId: 2,
        img_key: "avatar2.png",
        first_name: "Bob",
        username: "bob321",
      },
    ];

    usersThatSubmittedAnswer.usersSubmittedAnswers = mockUsersSubmittedAnswers;

    const wrapper = mountComponent();
    await wrapper.setProps({
      event: useNuxtApp.$GetQuestion,
      data: { totalJoinUser: 2 },
    });

    expect(wrapper.emitted("autoSkip")).toBeTruthy();
  });

  it("displays the correct number of user chips", async () => {
    const wrapper = mountComponent();

    const chips = wrapper.findAll(".chip");
    expect(chips.length).toBe(2); // Ensure all user chips are rendered
    expect(chips[0].text()).toContain("Alice (alice123)");
    expect(chips[1].text()).toContain("Bob (bob321)");
  });
});

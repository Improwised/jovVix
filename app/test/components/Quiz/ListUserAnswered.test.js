import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import ListUserAnswerd from "~/components/Quiz/ListUserAnswered.vue";
import { createTestingPinia } from "@pinia/testing";
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";

vi.mock("notivue", () => ({
  usePush: vi.fn(() => ({
    error: vi.fn(),
  })),
}));

vi.mock("~~/composables/avatar", () => ({
  getAvatarUrlByName: vi.fn(
    () => "https://api.dicebear.com/10.x/clay/svg?seed=Eden"
  ),
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
        NuxtImg: {
          template:
            '<img :src="$attrs.src" :alt="$attrs.alt" class="size-full object-cover" />',
        },
      },
    },
  });
};

describe("LlistUserAnswered Test", () => {
  it("renders correctly when no answers are submitted", async () => {
    expect(usersThatSubmittedAnswer.usersSubmittedAnswers.length).toBe(0);

    const wrapper = mountComponent();

    expect(wrapper.text()).toContain("No one answered till now");
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
    expect(wrapper.text()).toContain("1");
    expect(wrapper.text()).toContain("players answered");
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
      data: { totalJoinUser: 2 },
    });

    expect(wrapper.emitted("autoSkip")).toBeTruthy();
  });

  it("displays the correct number of user chips", async () => {
    const wrapper = mountComponent();

    const chips = wrapper.findAll("li");
    expect(chips.length).toBe(2);
    expect(chips[0].text()).toContain("Alice");
    expect(chips[1].text()).toContain("Bob");
  });
});

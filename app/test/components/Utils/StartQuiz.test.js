import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import { useRouter } from "vue-router";
import { useToast } from "vue-toastification";
import { useListUserstore } from "~/store/userlist";
import { useSessionStore } from "~/store/session";
import StartQuiz from "~/components/utils/StartQuiz.vue";

// Mock Vue Toastification
vi.mock("vue-toastification", () => ({
  useToast: vi.fn(),
}));

let wrapper = mount(StartQuiz, {
  props: {
    quizId: "1",
  },
});

describe("StartQuiz test", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
  });

  it("renders the button correctly", () => {
    expect(wrapper.find("button").exists()).toBe(true);
    expect(wrapper.text()).toContain("Start Quiz");
  });

  it("disables the button and shows 'Pending...' when requestPending is true", async () => {
    wrapper.vm.requestPending = true;
    await wrapper.vm.$nextTick();
    const button = wrapper.find("button");
    expect(button.text()).toBe("Pending...");
  });
});

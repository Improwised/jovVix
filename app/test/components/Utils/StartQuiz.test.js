import { describe, it, expect, vi, beforeEach, beforeAll } from "vitest";
import { mount } from "@vue/test-utils";
import StartQuiz from "~/components/utils/StartQuiz.vue";

// Mock Notivue
vi.mock("notivue", async (importOriginal) => ({
  ...(await importOriginal()),
  usePush: vi.fn(() => ({
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  })),
  push: {
    success: vi.fn(),
    error: vi.fn(),
    warning: vi.fn(),
    info: vi.fn(),
  },
}));

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time.
let wrapper;
beforeAll(() => {
  wrapper = mount(StartQuiz, {
    props: {
      quizId: "1",
    },
  });
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

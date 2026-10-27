import { describe, it, expect, afterEach, beforeAll, vi } from "vitest";
import { createTestingPinia } from "@pinia/testing";
import WaitingSpace from "~/components/Quiz/WaitingSpace.vue";
import { mount } from "@vue/test-utils";
import { useInvitationCodeStore } from "~/store/invitationcode";
import { useListUserstore } from "~/store/userlist";
import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import usecopyToClipboard from "~/composables/copy_to_clipboard";

const mock = vi.hoisted(() => {
  return {
    usecopyToClipboard: vi.fn(),
  };
});

// Mock composable for clipboard functionality
vi.mock("../../../composables/copy_to_clipboard.js", () => {
  return {
    default: mock.usecopyToClipboard,
  };
});

// Keep the real config so the Nuxt app (router, plugins) can still start;
// only the public baseUrl matters to these tests.
mockNuxtImport("useRuntimeConfig", (original) => {
  return () => {
    const config = original();
    return {
      ...config,
      public: { ...config.public, baseUrl: "https//quiz.example.com" },
    };
  };
});

// Stub after the Nuxt app has started: @nuxt/test-utils sets it up in a
// beforeAll hook and needs the real useNuxtApp to do so.
beforeAll(() => {
  vi.stubGlobal("useNuxtApp", () => ({}));
});
const pinia = createTestingPinia({
  createSpy: vi.fn,
});
const invitationCodeStore = useInvitationCodeStore(pinia);
invitationCodeStore.invitationCode = "123456";

const listUserStore = useListUserstore(pinia);

const mountComponent = () =>
  mount(WaitingSpace, {
    props: {
      data: { status: "success", data: "Welcome" },
      isAdmin: true,
    },
    global: {
      plugins: [pinia],
      stubs: {
        FontAwesomeIcon: true,
        NuxtImg: true,
      },
    },
  });

describe("WaitingSpace test", () => {
  afterEach(() => {
    vi.unstubAllGlobals();
    vi.restoreAllMocks();
    // Vitest 4's restoreAllMocks no longer resets vi.fn() call history.
    vi.clearAllMocks();
  });

  it("renders properly", async () => {
    const wrapper = mountComponent();
    expect(wrapper.find(".join-page-title").text()).toBe("Ready, Steady, Go!");
    expect(wrapper.find('button[type="submit"]').text()).toBe("Start Quiz");
  });

  it("shows correct invitation code", () => {
    const wrapper = mountComponent({
      data: { status: "success", data: "Welcome" },
      isAdmin: true,
    });
    expect(wrapper.find(".code").text()).toBe("123456");
  });

  it("calls `start_quiz` when the Start Quiz button is clicked", async () => {
    const wrapper = mountComponent();
    await wrapper.find("form").trigger("submit.prevent");
    expect(wrapper.emitted().startQuiz).toBeTruthy();
  });

  it("handles invitation code copying correctly", async () => {
    const wrapper = mountComponent();
    const copyBtn = wrapper.find("#OTP-input-container");
    expect(copyBtn.exists()).toBe(true);
    await copyBtn.trigger("click");
    await copyBtn.trigger("click");

    expect(usecopyToClipboard).toHaveBeenCalled();
    expect(usecopyToClipboard).toHaveBeenCalledWith("123456");
  });

  it("handles URL copying correctly", async () => {
    const wrapper = mountComponent();
    const urlCopyBtn = wrapper.find("#URL-input-container");
    await urlCopyBtn.trigger("click");
    expect(usecopyToClipboard).toHaveBeenCalledWith(
      "https//quiz.example.com/join?code=123456"
    );
  });

  it("pauses music on unmount if music is playing", () => {
    const pauseMock = vi.fn();
    const wrapper = mountComponent();
    wrapper.vm.waitingSound = { pause: pauseMock, play: vi.fn() };

    wrapper.unmount();
    expect(pauseMock).toHaveBeenCalled();
  });

  it("clears all users on unmount when `isAdmin` is true", () => {
    const wrapper = mountComponent();
    wrapper.unmount();
    expect(listUserStore.removeAllUsers).toHaveBeenCalled();
  });

  it("does not clear users on unmount when `isAdmin` is false", async () => {
    const wrapper = mountComponent();
    await wrapper.setProps({ isAdmin: false });
    wrapper.unmount();
    expect(listUserStore.removeAllUsers).not.toHaveBeenCalled();
  });

  it("emits `terminateQuiz` on unmount when quiz is not started", () => {
    const wrapper = mountComponent();
    wrapper.unmount();
    expect(wrapper.emitted().terminateQuiz).toBeTruthy();
  });

  it("does not emit `terminateQuiz` on unmount if quiz is started", async () => {
    const wrapper = mountComponent();
    wrapper.vm.startQuiz = true;
    wrapper.unmount();
    expect(wrapper.emitted().terminateQuiz).toBeFalsy();
  });
});

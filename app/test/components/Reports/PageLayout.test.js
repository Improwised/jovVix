import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import PageLayout from "~/components/reports/PageLayout.vue";

let wrapper = mount(PageLayout, {
  props: {
    currentTab: "report",
  },
});

describe("QuizAnalysisTabs.vue", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it("renders the navigation tabs", () => {
    const buttons = wrapper.findAll("button");
    expect(buttons).toHaveLength(2);
    expect(buttons[0].text()).toBe("Questions");
    expect(buttons[1].text()).toBe("Participants");
  });

  it("marks the correct tab as active based on `currentTab` prop", () => {
    expect(wrapper.find("button").text()).toBe("Questions");
  });

  it("emits `changeTab` with the correct value when a tab is clicked", async () => {
    const buttons = wrapper.findAll("button");
    const participantsTab = buttons[1];
    await participantsTab.trigger("click");

    expect(wrapper.emitted().changeTab).toBeTruthy();
    expect(wrapper.emitted().changeTab[0]).toEqual(["participants"]);
  });
});

import { describe, it, expect, beforeAll } from "vitest";
import { mountSuspended } from "@nuxt/test-utils/runtime";
import NavigationButton from "~/components/utils/NavigationButton.vue";

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time. Nuxt 4.7 renders NuxtLink inline, so
// it cannot be stubbed with RouterLinkStub.
let wrapper;
beforeAll(async () => {
  wrapper = await mountSuspended(NavigationButton, {
    props: {
      title: "Test Title",
      navigateTo: "/test-path",
    },
  });
});

describe("NavigationButton test", () => {
  it("renders the button with the correct title", () => {
    expect(wrapper.text()).toBe("Test Title");
  });

  it("renders a NuxtLink with the correct 'to' attribute", () => {
    const link = wrapper.find("a");
    expect(link.exists()).toBe(true);
    expect(link.attributes("href")).toBe("/test-path");
  });
});

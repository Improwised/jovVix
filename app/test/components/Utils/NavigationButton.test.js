import { describe, it, expect } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import NavigationButton from "~/components/utils/NavigationButton.vue";

let wrapper = mount(NavigationButton, {
  props: {
    title: "Test Title",
    navigateTo: "/test-path",
  },
  global: {
    stubs: {
      NuxtLink: RouterLinkStub,
    },
  },
});

describe("NavigationButton test", () => {
  it("renders the button with the correct title", () => {
    expect(wrapper.text()).toBe("Test Title");
  });

  it("renders a NuxtLink with the correct 'to' attribute", () => {
    const link = wrapper.findComponent(RouterLinkStub);
    expect(link.props().to).toBe("/test-path");
  });
});

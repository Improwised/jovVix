import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import Option from "~/components/Option.vue";

// Mock child component (e.g., CodeBlockComponent)
vi.mock("~/components/CodeBlockComponent.vue", () => ({
  default: {
    template: "<pre>{{ code }}</pre>",
    props: ["code"],
  },
}));

const mountComponent = (props) => {
  return mount(Option, {
    props,
  });
};

let wrapper = mountComponent({
  order: 1,
  optionsMedia: "text",
  option: "Option A",
});

describe("Option test", () => {
  it("renders the order correctly", () => {
    const button = wrapper.find("button");
    expect(button.text()).toBe("A");
  });

  it("renders image media correctly", () => {
    wrapper = mountComponent({
      optionsMedia: "image",
      option: "/path/to/image.jpg",
    });
    const img = wrapper.find("img");
    expect(img.exists()).toBe(true);
    expect(img.attributes("src")).toBe("/path/to/image.jpg");
    expect(img.attributes("alt")).toBe("/path/to/image.jpg");
  });

  it("renders text media correctly", () => {
    wrapper = mountComponent({
      optionsMedia: "text",
      option: "Option Text",
      isCorrect: true,
    });
    const textDiv = wrapper.find(".mx-3.font-weight-bold");
    expect(textDiv.exists()).toBe(true);
    expect(textDiv.text()).toBe("Option Text");
    expect(textDiv.classes()).toContain("text-success"); // isCorrect is true
  });

  it("renders code media correctly", () => {
    wrapper = mountComponent({
      optionsMedia: "code",
      option: "console.log('Hello, world!');",
    });
    const codeBlock = wrapper.find("pre");
    expect(codeBlock.exists()).toBe(true);
    expect(codeBlock.text()).toBe("console.log('Hello, world!');");
  });

  it("renders admin analysis badge when isAdminAnalysis is true", () => {
    wrapper = mountComponent({
      isAdminAnalysis: true,
      isCorrect: true,
      selected: 42,
      optionsMedia: "text",
    });
    const badge = wrapper.find(".badge");
    expect(badge.exists()).toBe(true);
    expect(badge.text()).toContain("42");
    expect(badge.classes()).toContain("bg-success"); // isCorrect is true
  });

  it("does not render admin analysis badge when isAdminAnalysis is false", () => {
    wrapper = mountComponent({
      isAdminAnalysis: false,
      selected: 42,
      optionsMedia: "text",
    });
    const badge = wrapper.find(".badge");
    expect(badge.exists()).toBe(false);
  });
});

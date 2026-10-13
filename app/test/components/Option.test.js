import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import Option from "~/components/Option.vue";

vi.mock("~/components/CodeBlockComponent.vue", () => ({
  default: {
    template: "<pre>{{ code }}</pre>",
    props: ["code"],
  },
}));

const mountComponent = (props) => {
  return mount(Option, {
    props,
    global: {
      stubs: {
        NuxtImg: {
          template:
            '<img :src="$attrs.src" :alt="$attrs.alt" class="max-h-[120px]" />',
        },
      },
    },
  });
};

let wrapper = mountComponent({
  order: 1,
  optionsMedia: "text",
  option: "Option A",
});

describe("Option test", () => {
  it("renders the order correctly", () => {
    expect(wrapper.text()).toContain("A");
  });

  it("renders image media correctly", () => {
    wrapper = mountComponent({
      optionsMedia: "image",
      option: "https://example.com/image.jpg",
    });
    const img = wrapper.find("img");
    expect(img.exists()).toBe(true);
    expect(img.attributes("src")).toBe("https://example.com/image.jpg");
  });

  it("renders text media correctly", () => {
    wrapper = mountComponent({
      optionsMedia: "text",
      option: "Option Text",
      isCorrect: true,
    });
    expect(wrapper.text()).toContain("Option Text");
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
    expect(wrapper.text()).toContain("42");
  });

  it("does not render admin analysis badge when isAdminAnalysis is false", () => {
    wrapper = mountComponent({
      isAdminAnalysis: false,
      selected: 42,
      optionsMedia: "text",
    });
    expect(wrapper.text()).not.toContain("42");
  });
});

import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import Playground from "~/components/Playground.vue";

Object.defineProperty(document, "fullscreenElement", {
  writable: true,
  value: null,
});

document.exitFullscreen = vi.fn().mockImplementation(() => {
  document.fullscreenElement = null;
});

const mountComponent = () => {
  return mount(Playground, {
    slots: {
      default: `<div>Slot Content</div>`,
    },
  });
};

const wrapper = mountComponent();

Element.prototype.requestFullscreen = vi.fn().mockImplementation(function () {
  document.fullscreenElement = this;
});

describe("Playground test", () => {
  it("renders the slot content", () => {
    expect(wrapper.html()).toContain("Slot Content");
  });

  it("requests full screen when fullScreenEnabled is true", async () => {
    await wrapper.setProps({ fullScreenEnabled: true });
    await wrapper.vm.$nextTick();

    expect(wrapper.emitted("isFullScreen")[0]).toEqual([true]);
  });

  it("exits full screen when fullScreenEnabled is false", async () => {
    const wrapper = mount(Playground, {
      props: {
        fullScreenEnabled: true,
      },
    });

    // Simulate the component entering full screen
    document.fullscreenElement = wrapper.element;

    await wrapper.setProps({ fullScreenEnabled: false });

    expect(document.exitFullscreen).toHaveBeenCalled();
    expect(wrapper.emitted("isFullScreen")[0]).toEqual([false]);
  });

  it("emits isFullScreen event on fullscreenchange event", async () => {
    // Simulate entering full screen
    document.fullscreenElement = wrapper.element;
    await wrapper.vm.$nextTick();

    document.dispatchEvent(new Event("fullscreenchange"));

    expect(wrapper.emitted("isFullScreen")[0]).toEqual([true]);
  });
});

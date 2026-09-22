import { mountSuspended } from "@nuxt/test-utils/runtime";
import { it, expect, describe } from "vitest";
import { nextTick } from "vue";
import QrFullscreenOverlay from "~~/components/QrFullscreenOverlay.vue";

describe("QrFullscreenOverlay Test", () => {
  it("closes on Escape", async () => {
    const wrapper = await mountSuspended(QrFullscreenOverlay, {
      props: { joinUrl: "https://example.com/join", code: "123456" },
    });

    await wrapper.vm.open();
    await nextTick();
    expect(document.body.querySelector('[role="dialog"]')).not.toBeNull();

    window.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));
    await nextTick();

    expect(document.body.querySelector('[role="dialog"]')).toBeNull();
    wrapper.unmount();
  });
});

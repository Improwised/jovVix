import { mount } from "@vue/test-utils";
import { describe, it, expect } from "vitest";
import ConfirmModal from "~/components/utils/ConfirmModal.vue";

const props = {
  modalTitle: "Test Modal",
  modalMessage: "test message",
  modelPositiveMessage: "this is positive message",
  modelNegativeMessage: "this is negative messsage",
};

const wrapper = mount(ConfirmModal, {
  props,
  attachTo: document.body,
  global: {
    stubs: {
      Teleport: true,
    },
  },
});

describe("ConfirmModal test", () => {
  it("renders with props", async () => {
    expect(wrapper.props("modalTitle")).toBe("Test Modal");
    expect(wrapper.props("modalMessage")).toBe("test message");
    expect(wrapper.props("modelPositiveMessage")).toBe(
      "this is positive message"
    );
    expect(wrapper.props("modelNegativeMessage")).toBe(
      "this is negative messsage"
    );
  });

  it("emits confirmMessage event on positive button click", async () => {
    const buttons = wrapper.findAll("button");
    const positiveButton = buttons.find((b) =>
      b.text().includes("this is positive message")
    );
    if (positiveButton) {
      await positiveButton.trigger("click");
      expect(wrapper.emitted("confirmMessage")).toBeTruthy();
      expect(wrapper.emitted("confirmMessage")[0]).toEqual([true]);
    }
  });
});

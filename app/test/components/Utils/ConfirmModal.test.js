import { mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import ConfirmModal from "~/components/utils/ConfirmModal.vue";

const props = {
  modalTitle: "Test Modal",
  modalMessage: "test message",
  modelPositiveMessage: "this is positive message",
  modelNegativeMessage: "this is negative messsage",
};

const wrapper = mount(ConfirmModal, {
  props,
  global: {
    mocks: {
      $bootstrap: {
        Modal: vi.fn(() => ({
          show: vi.fn(),
          hide: vi.fn(),
          _element: {
            addEventListener: vi.fn(),
          },
        })),
      },
    },
  },
});

describe("ConfirmModal test", () => {
  it("renders with props", async () => {
    // Assert default props
    expect(wrapper.props("modalTitle")).toBe("Test Modal");
    expect(wrapper.props("modalMessage")).toBe("test message");
    expect(wrapper.props("modelPositiveMessage")).toBe(
      "this is positive message"
    );
    expect(wrapper.props("modelNegativeMessage")).toBe(
      "this is negative messsage"
    );

    // Assert modal DOM structure
    expect(wrapper.find("#confirmModal").exists()).toBe(true);
    expect(wrapper.find(".modal-title").text()).toBe("Test Modal");
    expect(wrapper.find(".modal-body").text()).toBe("test message");
  });

  it("emits confirmMessage event on positive button click", async () => {
    const positiveButton = wrapper.find(".btn-primary");
    await positiveButton.trigger("click");

    // Assert the event is emitted with the correct value
    expect(wrapper.emitted("confirmMessage")).toBeTruthy();
    expect(wrapper.emitted("confirmMessage")[0]).toEqual([true]);
  });
});

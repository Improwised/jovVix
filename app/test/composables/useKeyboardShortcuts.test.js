import { describe, it, expect, afterEach, vi } from "vitest";
import { mount } from "@vue/test-utils";
import { h } from "vue";
import { useKeyboardShortcuts } from "~~/composables/useKeyboardShortcuts";

const press = (key, init = {}) =>
  window.dispatchEvent(
    new KeyboardEvent("keydown", { key, bubbles: true, ...init })
  );

const mountWithShortcuts = (handlers) =>
  mount({
    setup() {
      useKeyboardShortcuts(handlers);
      return () => h("div");
    },
  });

describe("useKeyboardShortcuts test", () => {
  afterEach(() => {
    document.body.innerHTML = "";
    vi.restoreAllMocks();
  });

  it("calls onDigit with the pressed number", () => {
    const onDigit = vi.fn();
    mountWithShortcuts({ onDigit });

    press("3");

    expect(onDigit).toHaveBeenCalledWith(3);
  });

  it("calls onEnter on Enter", () => {
    const onEnter = vi.fn();
    mountWithShortcuts({ onEnter });

    press("Enter");

    expect(onEnter).toHaveBeenCalledTimes(1);
  });

  it("stays silent while a modal is open", () => {
    const onDigit = vi.fn();
    const onEnter = vi.fn();
    mountWithShortcuts({ onDigit, onEnter });

    const modal = document.createElement("div");
    modal.setAttribute("role", "dialog");
    modal.setAttribute("aria-modal", "true");
    document.body.appendChild(modal);

    press("1");
    press("Enter");

    expect(onDigit).not.toHaveBeenCalled();
    expect(onEnter).not.toHaveBeenCalled();
  });

  it("stays silent while typing in an input", () => {
    const onDigit = vi.fn();
    mountWithShortcuts({ onDigit });

    const input = document.createElement("input");
    document.body.appendChild(input);
    input.focus();

    press("1");

    expect(onDigit).not.toHaveBeenCalled();
  });

  it("stays silent on Enter while a button is focused, so native activation survives", () => {
    const onEnter = vi.fn();
    mountWithShortcuts({ onEnter });

    const button = document.createElement("button");
    document.body.appendChild(button);
    button.focus();

    press("Enter");

    expect(onEnter).not.toHaveBeenCalled();
  });

  it("stays silent while a non-modal popover is open (data-state=open)", () => {
    const onDigit = vi.fn();
    const onEnter = vi.fn();
    mountWithShortcuts({ onDigit, onEnter });

    const popover = document.createElement("div");
    popover.setAttribute("data-state", "open");
    document.body.appendChild(popover);

    press("1");
    press("Enter");

    expect(onDigit).not.toHaveBeenCalled();
    expect(onEnter).not.toHaveBeenCalled();
  });

  it("ignores auto-repeat from a held key", () => {
    const onEnter = vi.fn();
    const onDigit = vi.fn();
    mountWithShortcuts({ onEnter, onDigit });

    press("Enter");
    press("Enter", { repeat: true });
    press("Enter", { repeat: true });
    press("1", { repeat: true });

    expect(onEnter).toHaveBeenCalledTimes(1);
    expect(onDigit).not.toHaveBeenCalled();
  });

  it("stops listening once the component unmounts", () => {
    const onEnter = vi.fn();
    const wrapper = mountWithShortcuts({ onEnter });

    wrapper.unmount();
    press("Enter");

    expect(onEnter).not.toHaveBeenCalled();
  });
});

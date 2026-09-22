import { onKeyStroke } from "@vueuse/core";

const DIGIT_KEYS = ["1", "2", "3", "4", "5", "6", "7", "8", "9"];

// Stay out of the way of typing, of a focused control's own native key
// handling (e.g. Enter activating a button), and of anything an open
// dialog/popover owns. Every modal in the app already closes itself on
// Escape, so that binding is intentionally not duplicated here.
// A held key repeats, which would fire the action dozens of times.
function isBlocked(e) {
  if (e.repeat) return true;
  const el = document.activeElement;
  if (el?.isContentEditable) return true;
  if (el?.matches?.("input, textarea, select")) return true;
  if (el?.matches?.("button, a[href], [role='button'], summary")) return true;
  return Boolean(
    document.querySelector('[role="dialog"], [data-state="open"]')
  );
}

export function useKeyboardShortcuts({ onDigit, onEnter } = {}) {
  if (onDigit) {
    onKeyStroke(DIGIT_KEYS, (e) => {
      if (isBlocked(e)) return;
      e.preventDefault();
      onDigit(Number(e.key));
    });
  }

  if (onEnter) {
    onKeyStroke("Enter", (e) => {
      if (isBlocked(e)) return;
      e.preventDefault();
      onEnter();
    });
  }
}

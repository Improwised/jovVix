import { onKeyStroke } from "@vueuse/core";

const DIGIT_KEYS = ["1", "2", "3", "4", "5", "6", "7", "8", "9"];

// Stay out of the way of typing, and of anything an open dialog/popover owns.
// Escape-to-close lives in each modal itself (they unmount on close), so it is
// not bound here. A held key repeats, which would fire the action dozens of times.
function isBlocked(e) {
  if (e.repeat) return true;
  const el = document.activeElement;
  if (el?.isContentEditable) return true;
  if (el?.matches?.("input, textarea, select")) return true;
  return Boolean(
    document.querySelector('[role="dialog"], [data-state="open"]')
  );
}

// Enter on a focused control must keep its native activation (e.g. Mute/Leave).
const isControlFocused = () =>
  Boolean(
    document.activeElement?.matches?.(
      "button, a[href], [role='button'], summary"
    )
  );

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
      if (isBlocked(e) || isControlFocused()) return;
      e.preventDefault();
      onEnter();
    });
  }
}

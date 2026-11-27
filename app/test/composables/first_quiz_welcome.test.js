import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import {
  markNewSignup,
  takeNewSignup,
} from "../../composables/first_quiz_welcome";

describe("new sign-up flag", () => {
  beforeEach(() => sessionStorage.clear());
  afterEach(() => vi.unstubAllGlobals());

  it("is false when nobody signed up in this tab", () => {
    expect(takeNewSignup()).toBe(false);
  });

  it("is true once after a sign-up, then false", () => {
    markNewSignup();
    expect(takeNewSignup()).toBe(true);
    expect(takeNewSignup()).toBe(false);
  });

  it("is false and does not throw when storage is blocked", () => {
    const blocked = () => {
      throw new Error("blocked");
    };
    vi.stubGlobal("sessionStorage", {
      getItem: blocked,
      setItem: blocked,
      removeItem: blocked,
    });
    expect(() => markNewSignup()).not.toThrow();
    expect(takeNewSignup()).toBe(false);
  });
});

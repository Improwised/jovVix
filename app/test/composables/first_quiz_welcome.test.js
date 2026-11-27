import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import {
  markNewSignup,
  takeNewSignup,
} from "../../composables/first_quiz_welcome";

describe("new sign-up flag", () => {
  beforeEach(() => sessionStorage.clear());
  afterEach(() => vi.restoreAllMocks());

  it("is false when nobody signed up in this tab", () => {
    expect(takeNewSignup()).toBe(false);
  });

  it("is true once after a sign-up, then false", () => {
    markNewSignup();
    expect(takeNewSignup()).toBe(true);
    expect(takeNewSignup()).toBe(false);
  });

  it("is false and does not throw when storage is blocked", () => {
    vi.spyOn(Storage.prototype, "getItem").mockImplementation(() => {
      throw new Error("blocked");
    });
    vi.spyOn(Storage.prototype, "setItem").mockImplementation(() => {
      throw new Error("blocked");
    });
    expect(() => markNewSignup()).not.toThrow();
    expect(takeNewSignup()).toBe(false);
  });
});

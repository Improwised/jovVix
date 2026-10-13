import { describe, expect, it } from "vitest";
import { readApiError } from "../../composables/apiError";

describe("readApiError", () => {
  it("returns error.data.data when available", () => {
    expect(
      readApiError({ data: { data: "data message", message: "fallback" } }, "")
    ).toBe("data message");
  });

  it("falls back to error.data.message", () => {
    expect(readApiError({ data: { message: "error message" } }, "")).toBe(
      "error message"
    );
  });

  it("falls back to error.message", () => {
    expect(readApiError({ message: "top level error" }, "")).toBe(
      "top level error"
    );
  });

  it("falls back to the fallback string when nothing is available", () => {
    expect(readApiError({}, "default message")).toBe("default message");
  });
});

import { describe, expect, it } from "vitest";
import { AI_PROVIDERS } from "../../config/aiProviders";
import { readAIGenerationError } from "../../composables/aiGenerationError";

describe("AI generation errors", () => {
  it.each([
    [429, "The AI provider is busy. Wait a moment, then try again.", false],
    [
      504,
      "The AI took too long to respond. Try again, or generate fewer questions.",
      false,
    ],
    [
      502,
      "The AI provider is temporarily unavailable. Try again in a moment.",
      false,
    ],
    [
      503,
      "The AI provider is temporarily unavailable. Try again in a moment.",
      false,
    ],
    [
      400,
      "Your AI settings need attention. Check the API key and selected model, then try again.",
      true,
    ],
    [
      401,
      "Your AI settings need attention. Check the API key and selected model, then try again.",
      true,
    ],
    [
      403,
      "Your AI settings need attention. Check the API key and selected model, then try again.",
      true,
    ],
  ])("maps HTTP %s to a useful next step", (status, message, needsSettings) => {
    expect(readAIGenerationError({ status })).toEqual({
      message,
      needsSettings,
    });
  });

  it("maps AI output and timeout failures without an HTTP status", () => {
    expect(
      readAIGenerationError({
        message: "could not read the questions returned by the ai service",
      }).message
    ).toBe("The AI did not return usable questions. Try again.");
    expect(
      readAIGenerationError({ message: "request timed out" }).message
    ).toBe(
      "The AI took too long to respond. Try again, or generate fewer questions."
    );
  });

  it("uses a safe retry message for unknown failures", () => {
    expect(readAIGenerationError(new Error("network failed"))).toEqual({
      message: "Could not generate questions. Try again.",
      needsSettings: false,
    });
  });
});

describe("AI provider labels", () => {
  it("marks only Groq as recommended", () => {
    const labels = AI_PROVIDERS.map((provider) =>
      provider.id === "groq"
        ? "Groq (Recommended · free tier)"
        : `${provider.label}${provider.freeTier ? " (free tier)" : ""}`
    );

    expect(labels).toContain("Groq (Recommended · free tier)");
    expect(
      labels.filter((label) => label.includes("Recommended"))
    ).toHaveLength(1);
  });
});

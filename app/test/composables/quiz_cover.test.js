import { describe, expect, it, beforeAll } from "vitest";
import { useQuizCover } from "../../composables/quiz_cover";

const fallbacks = [
  "/images/landing/homepage-public-quiz-1.png",
  "/images/landing/homepage-public-quiz-2.png",
  "/images/landing/homepage-public-quiz-3.png",
  "/images/landing/homepage-public-quiz-4.png",
];

const tilts = [
  "rotate-[-1deg]",
  "rotate-[1deg]",
  "rotate-[-0.4deg]",
  "rotate-[0.6deg]",
];

describe("useQuizCover", () => {
  let fallbackUrl, coverUrl, tiltClass;
  // useQuizCover reads runtime config, which is only available once
  // @nuxt/test-utils has started the Nuxt app in its beforeAll hook.
  beforeAll(() => {
    ({ fallbackUrl, coverUrl, tiltClass } = useQuizCover());
  });

  it("returns one of the known fallback images", () => {
    const url = fallbackUrl({ id: "42" });
    expect(fallbacks).toContain(url);
  });

  it("returns a cover URL when quiz has cover image", () => {
    const url = coverUrl({
      id: "42",
      has_cover_image: true,
      updated_at: "2026-10-13T00:00:00.000Z",
    });
    expect(url).toContain("/quizzes/42/cover");
    expect(url).toContain("v=");
  });

  it("falls back when quiz has no cover image", () => {
    const url = coverUrl({ id: "99", has_cover_image: false });
    expect(fallbacks).toContain(url);
  });

  it("returns the same fallback for the same quiz id", () => {
    const a = fallbackUrl({ id: "7" });
    const b = fallbackUrl({ id: "7" });
    expect(a).toBe(b);
  });

  it("returns one of the known tilt classes", () => {
    const t = tiltClass({ id: "42" });
    expect(tilts).toContain(t);
  });
});

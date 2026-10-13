import { describe, expect, it, vi, afterEach } from "vitest";
import { useGetTime } from "../../composables/quiz_details";

describe("useGetTime", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('returns "just now" for less than 60 seconds', () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 30000).toISOString());
    expect(result).toBe("This quiz was created just now.");
  });

  it("returns minutes ago for less than 60 minutes", () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 120000).toISOString());
    expect(result).toBe("This quiz was created 2 minutes ago");
  });

  it("returns hours ago for less than 24 hours", () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 3600000 * 3).toISOString());
    expect(result).toBe("This quiz was created 3 hours ago");
  });

  it("returns days ago for less than 30 days", () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 86400000 * 5).toISOString());
    expect(result).toBe("This quiz was created 5 days ago");
  });

  it("returns months ago for less than 12 months", () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 86400000 * 60).toISOString());
    expect(result).toBe("This quiz was created 2 months ago");
  });

  it("returns years ago for 12+ months", () => {
    const now = Date.now();
    vi.setSystemTime(now);
    const result = useGetTime(new Date(now - 86400000 * 400).toISOString());
    expect(result).toBe("This quiz was created 1 years ago");
  });
});

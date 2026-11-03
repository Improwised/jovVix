import { describe, expect, it } from "vitest";
import { joinRecoveryQuery } from "~/composables/join_recovery";

describe("joinRecoveryQuery", () => {
  it("returns a direct-link name conflict to the join form with its code", () => {
    const query = joinRecoveryQuery("123456", { status: 409 });
    expect(query.code).toBe("123456");
    expect(query.error).toContain("already taken");
    expect(query.error).not.toContain("ended");
  });

  it("distinguishes an invalid name from an expired invitation", () => {
    expect(
      joinRecoveryQuery("123456", {
        status: 400,
        data: { data: "player name must be between 1 and 12 characters" },
      }).error
    ).toContain("valid name");
    expect(
      joinRecoveryQuery("123456", {
        status: 400,
        data: { data: "invitation code not found" },
      }).error
    ).toBe("invitation code not found");
  });

  it("does not label network errors as ended quizzes", () => {
    expect(joinRecoveryQuery("123456", new Error("offline")).error).toContain(
      "try again"
    );
  });
});

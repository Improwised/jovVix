import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import UserOperation from "~/composables/user_operation";
import AdminOperations from "~/composables/admin_operation";

mockNuxtImport("setSocketObject", () => vi.fn());

class FakeSocket {
  static OPEN = 1;
  static CLOSED = 3;
  static instances = [];
  readyState = 1;
  constructor() {
    FakeSocket.instances.push(this);
  }
  close() {
    this.readyState = 3;
  }
  send() {}
}

describe("quiz connection recovery", () => {
  let handler;
  beforeEach(() => {
    vi.useFakeTimers();
    FakeSocket.instances = [];
    vi.stubGlobal("WebSocket", FakeSocket);
  });
  afterEach(() => {
    handler?.dispose();
    vi.useRealTimers();
    vi.unstubAllGlobals();
  });

  it("uses one retry scheduler for a refused player connection and reconciles the room first", async () => {
    const recover = vi.fn().mockResolvedValue(false);
    handler = new UserOperation(
      "123456",
      "player",
      vi.fn(),
      vi.fn(),
      vi.fn(),
      vi.fn(),
      recover
    );
    handler.onClose();
    handler.onError(new Error("Connection refused"));
    expect(vi.getTimerCount()).toBe(1);
    await vi.advanceTimersByTimeAsync(2000);
    expect(recover).toHaveBeenCalledTimes(1);
    expect(FakeSocket.instances).toHaveLength(1);
    expect(vi.getTimerCount()).toBe(0);
  });

  it.each(["player", "host"])(
    "does not restart a disposed %s connection when room recovery finishes late",
    async (role) => {
      let resolve;
      const recover = vi.fn(
        () =>
          new Promise((done) => {
            resolve = done;
          })
      );
      if (role === "player") {
        handler = new UserOperation(
          "123456",
          "player",
          vi.fn(),
          vi.fn(),
          vi.fn(),
          vi.fn(),
          recover
        );
      } else {
        handler = new AdminOperations(
          "session",
          vi.fn(),
          vi.fn(),
          vi.fn(),
          recover
        );
        handler.connectAdmin();
      }
      handler.onClose();
      await vi.advanceTimersByTimeAsync(2000);
      expect(recover).toHaveBeenCalledTimes(1);
      handler.dispose();
      resolve(true);
      await Promise.resolve();
      await Promise.resolve();
      expect(FakeSocket.instances).toHaveLength(1);
      expect(vi.getTimerCount()).toBe(0);
    }
  );
  it.each(["player", "host"])(
    "keeps recovering a %s after a long interruption",
    async (role) => {
      const recover = vi.fn().mockRejectedValue(new Error("offline"));
      handler =
        role === "player"
          ? new UserOperation(
              "123456",
              "player",
              vi.fn(),
              vi.fn(),
              vi.fn(),
              vi.fn(),
              recover
            )
          : new AdminOperations("session", vi.fn(), vi.fn(), vi.fn(), recover);
      handler.onClose();
      await vi.advanceTimersByTimeAsync(120000);
      expect(recover.mock.calls.length).toBeGreaterThan(5);
      expect(vi.getTimerCount()).toBe(1);
      handler.dispose();
      expect(vi.getTimerCount()).toBe(0);
    }
  );
});

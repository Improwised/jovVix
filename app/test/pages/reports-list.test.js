import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, beforeEach, vi } from "vitest";
import { ref } from "vue";
import Reports from "~/pages/admin/reports/index.vue";

const route = vi.hoisted(() => ({ path: "/admin/reports", query: {} }));
mockNuxtImport("useRoute", () => () => route);

const navigate = vi.hoisted(() => vi.fn());
mockNuxtImport("navigateTo", () => navigate);

const session = (id, quizId, sessions) => ({
  id,
  quiz_id: quizId,
  sessions,
  title: "Space%20Quiz",
  description: { String: "" },
  activated_from: { Time: "01-Nov-2026 10:00:00" },
  activated_to: { Time: "01-Nov-2026 10:10:00" },
  participants: 2,
  questions: 2,
  correct_answers: 2,
});

mockNuxtImport("useFetch", () => () => ({
  data: ref({
    data: {
      count: 2,
      data: [session("s1", "quiz-a", 2), session("s9", "quiz-b", 1)],
    },
  }),
  error: ref(null),
  pending: ref(false),
}));

let wrapper;
beforeEach(async () => {
  navigate.mockClear();
  wrapper = mount(Reports, {
    global: { stubs: { RoomReportList: true, Pagination: true } },
  });
  await flushPromises();
});

describe("Reports list combined report button", () => {
  it("shows only on quizzes with more than one ended session", () => {
    const buttons = wrapper
      .findAll("button")
      .filter((b) => b.text().includes("Combine"));
    expect(buttons.map((b) => b.text())).toEqual(["Combine 2 sessions"]);
  });

  it("explains the combine button in a tooltip on hover and focus", () => {
    const button = wrapper
      .findAll("button")
      .find((b) => b.text().includes("Combine"));
    const tip = wrapper.find(`#${button.attributes("aria-describedby")}`);
    expect(tip.attributes("role")).toBe("tooltip");
    expect(tip.text()).toBe(
      "Combine results from multiple sessions of this quiz into one leaderboard."
    );
    // Hidden until the button's wrapper is hovered or holds keyboard focus.
    expect(tip.classes()).toEqual(
      expect.arrayContaining([
        "invisible",
        "group-hover:visible",
        "group-focus-within:visible",
      ])
    );
  });

  it("opens the combined report, not the session report", async () => {
    const button = wrapper
      .findAll("button")
      .find((b) => b.text().includes("Combine"));
    await button.trigger("click");
    expect(navigate).toHaveBeenCalledTimes(1);
    expect(navigate).toHaveBeenCalledWith("/admin/reports/combined/quiz-a");
  });
});

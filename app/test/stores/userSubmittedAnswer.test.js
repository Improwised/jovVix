import { describe, it, expect } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import { useUserThatSubmittedAnswer } from "~/store/userSubmittedAnswer";

describe("submitted answer roster", () => {
  it("counts restored and live events for the same user only once", () => {
    setActivePinia(createPinia());
    const store = useUserThatSubmittedAnswer();
    store.addUserSubmittedAnswer({ id: "one", first_name: "One" });
    store.addUserSubmittedAnswer({ UserId: "one", first_name: "One" });
    store.addUserSubmittedAnswer({ user_id: "two", first_name: "Two" });
    expect(store.usersSubmittedAnswers).toHaveLength(2);
    store.resetUsersSubmittedAnswers();
    expect(store.usersSubmittedAnswers).toHaveLength(0);
  });
});

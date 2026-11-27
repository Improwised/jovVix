import { mount, flushPromises } from "@vue/test-utils";
import { describe, it, expect, beforeEach } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import { useUsersStore } from "~/store/users";
import { markNewSignup } from "~/composables/first_quiz_welcome";
import FirstQuizWelcome from "~/components/landing/FirstQuizWelcome.vue";

const mountPanel = () =>
  mount(FirstQuizWelcome, {
    global: {
      stubs: {
        NuxtLink: { props: ["to"], template: '<a :href="to"><slot /></a>' },
      },
    },
  });

describe("first quiz panel after sign-up", () => {
  let store;
  beforeEach(() => {
    sessionStorage.clear();
    setActivePinia(createPinia());
    store = useUsersStore();
  });

  it("shows two ways to start, and no AI, for a new sign-up", async () => {
    markNewSignup();
    store.setUserData({ role: "admin-user", firstname: "Asha" });
    const wrapper = mountPanel();
    await flushPromises();
    const hrefs = wrapper.findAll("a").map((a) => a.attributes("href"));
    expect(hrefs).toContain("#public-quiz");
    expect(hrefs).toContain("/admin/quiz/list-quiz?create=1");
    expect(wrapper.text()).not.toMatch(/AI/);
  });

  it("shows once: the flag is gone after the first mount", async () => {
    markNewSignup();
    store.setUserData({ role: "admin-user" });
    mountPanel();
    await flushPromises();
    expect(mountPanel().find("section").exists()).toBe(false);
  });

  it("waits for the user data that arrives after mount", async () => {
    markNewSignup();
    const wrapper = mountPanel();
    await flushPromises();
    expect(wrapper.find("section").exists()).toBe(false);
    store.setUserData({ role: "admin-user" });
    await flushPromises();
    expect(wrapper.find("section").exists()).toBe(true);
  });

  it("stays hidden for a visitor who did not just sign up", async () => {
    store.setUserData({ role: "admin-user" });
    const wrapper = mountPanel();
    await flushPromises();
    expect(wrapper.find("section").exists()).toBe(false);
  });
});

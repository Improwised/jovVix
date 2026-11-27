import { mount, flushPromises } from "@vue/test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { createPinia, setActivePinia } from "pinia";
import { useUsersStore } from "~/store/users";
import Sidebar from "~/components/common/Sidebar.vue";

vi.mock("@/composables/auth", () => ({
  handleLogout: vi.fn(),
  setUserDataStore: vi.fn(async () => {}),
}));
vi.mock("@/composables/email_verification", async () => {
  const { ref } = await import("vue");
  return {
    refreshEmailVerified: vi.fn(),
    useEmailVerified: () => ref(true),
  };
});

const mountSidebar = async () => {
  const wrapper = mount(Sidebar, {
    global: { stubs: { NuxtImg: true, NuxtLink: true } },
  });
  await flushPromises();
  return wrapper;
};

const aiLabel = "Generate with AI";

// A new host without an AI provider key stops at the key form, so signed-out
// visitors are not offered AI generation.
describe("sidebar AI entry", () => {
  let store;
  beforeEach(() => {
    setActivePinia(createPinia());
    store = useUsersStore();
  });

  it("is not offered to a signed-out visitor", async () => {
    const wrapper = await mountSidebar();
    expect(wrapper.find("aside").text()).not.toContain(aiLabel);
    expect(wrapper.vm.mobileNavItems.map((i) => i.label)).not.toContain(
      aiLabel
    );
  });

  it("stays for a signed-in host", async () => {
    store.setUserData({ role: "admin-user", firstname: "Asha" });
    const wrapper = await mountSidebar();
    expect(wrapper.find("aside").text()).toContain(aiLabel);
    expect(wrapper.vm.mobileNavItems.map((i) => i.label)).toContain(aiLabel);
  });
});

import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import ListJoinUser from "~/components/ListJoinUser.vue";
import { createTestingPinia } from "@pinia/testing";
import { useListUserstore } from "~/store/userlist";

vi.mock("~~/composables/avatar", () => ({
  getAvatarUrlByName: vi.fn(
    () => "https://api.dicebear.com/10.x/clay/svg?seed=Eden"
  ),
}));

const mockListUsers = [
  { UserId: 1, UserName: "Alice", Avatar: "Alice" },
  { UserId: 2, UserName: "Bob", Avatar: "Bob" },
];

const pinia = createTestingPinia({
  createSpy: vi.fn,
});

const listUserStore = useListUserstore(pinia);
listUserStore.listUsers = mockListUsers;

const mountComponent = () =>
  mount(ListJoinUser, {
    global: {
      plugins: [pinia],
      stubs: {
        NuxtImg: {
          template:
            '<img :src="$attrs.src" alt="" class="size-12 rounded-full border-2 border-jv-ink object-cover" />',
        },
      },
    },
  });

const wrapper = mountComponent();

describe("ListJoinUser test", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows the participant count when listUsers is not empty", () => {
    expect(wrapper.find("h5").text()).toBe("2 Participants");
  });

  it("renders user chips with correct data", () => {
    const cards = wrapper.findAll('[class*="jv-card"]');
    expect(cards.length).toBe(2);

    expect(cards[0].text()).toContain("Alice");
    expect(cards[1].text()).toContain("Bob");
  });

  it("shows 'Waiting for Participants' when listUsers is empty", async () => {
    listUserStore.listUsers = [];
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h5").text()).toBe("Waiting for Participants...");
    expect(wrapper.findAll('[class*="jv-card"]').length).toBe(0);
  });
});

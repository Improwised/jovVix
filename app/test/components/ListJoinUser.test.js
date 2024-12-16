import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import ListJoinUser from "~/components/ListJoinUser.vue";
import { createTestingPinia } from "@pinia/testing";
import { useListUserstore } from "~/store/userlist";
import constants from "~/test/constants";

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
        FontAwesomeIcon: true,
        VCard: constants.slotTemplate,
      },
    },
  });

const wrapper = mountComponent();

describe("ListJoinUser test", () => {
  beforeEach(() => {
    // Reset mocks
    vi.clearAllMocks();
  });

  it("shows the participant count when listUsers is not empty", () => {
    expect(wrapper.find("h5").text()).toBe("2 Participants");
  });

  it("renders user chips with correct data", () => {
    const chips = wrapper.findAll(".chip");
    expect(chips.length).toBe(2);

    expect(chips[0].text()).toContain("Alice");
    expect(chips[0].find("img").attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Eden"
    );

    expect(chips[1].text()).toContain("Bob");
    expect(chips[1].find("img").attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Eden"
    );
  });

  it("shows 'Waiting for Participants' when listUsers is empty", async () => {
    listUserStore.listUsers = [];
    await wrapper.vm.$nextTick();
    expect(wrapper.find("h5").text()).toBe("Waiting for Participants..");
    expect(wrapper.findAll(".chip").length).toBe(0);
  });
});

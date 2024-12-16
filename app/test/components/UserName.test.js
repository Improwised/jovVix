import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import UserName from "~/components/UserName.vue";
import { useUsersStore } from "~/store/users";
import { createTestingPinia } from "@pinia/testing";

const pinia = createTestingPinia({
  createSpy: vi.fn,
});
const userData = useUsersStore(pinia);
const mountComponent = () => {
  return mount(UserName, {
    props: { userName: "John Doe" },
    global: {
      plugins: [pinia],
    },
  });
};

let wrapper = mountComponent();

describe("UserName test", () => {
  beforeEach(() => {
    vi.resetAllMocks();
    vi.restoreAllMocks();
  });

  it("renders the user's avatar if available", async () => {
    userData.getUserData = vi.fn().mockReturnValue({
      avatar: "Sophia",
      name: "John Doe",
    });

    wrapper = mountComponent();
    const avatarImg = wrapper.find("img");

    expect(userData.getUserData).toBeCalled();
    expect(avatarImg.attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Sophia"
    );
    expect(wrapper.text()).toContain("John Doe");
  });

  it("renders the default avatar if user avatar is not available", () => {
    wrapper = mountComponent();
    const avatarImg = wrapper.find("img");
    expect(avatarImg.attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Eden"
    );
    expect(wrapper.text()).toContain("John Doe");
  });

  it("renders the user's name properly", () => {
    const wrapper = mount(UserName, {
      props: { userName: "John Doe" },
    });
    expect(wrapper.text()).toContain("John Doe");
  });
});

import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import WinnerCard from "@/components/WinnerCard.vue";
import { getAvatarUrlByName } from "../../composables/avatar";

const props = {
  winner: {
    rank: 1,
    img_key: "Sophia",
    firstname: "John",
    username: "john_doe",
    score: 8256,
  },
};

const mountComponent = () => {
  return mount(WinnerCard, {
    props,
  });
};

let wrapper = mountComponent();

vi.mock("../../composables/avatar.js", () => {
  return {
    getAvatarUrlByName: vi
      .fn()
      .mockReturnValue("https://api.dicebear.com/9.x/bottts/svg?seed=Sophia"),
  };
});

describe("WinnerCard test", () => {
  it("renders the correct medal for rank 1", () => {
    const medal = wrapper.find("img.bg-image");
    expect(medal.attributes("src")).toBe("/assets/images/medal/1.webp");
  });

  it("renders the correct medal for rank 2", async () => {
    props.winner.rank = 2;
    wrapper.unmount();
    wrapper = mountComponent();

    const medal = wrapper.find("img.bg-image");
    expect(medal.attributes("src")).toBe("/assets/images/medal/2.webp");
  });

  it("renders the correct medal for rank 3", async () => {
    props.winner.rank = 3;
    wrapper.unmount();
    wrapper = mountComponent();

    const medal = wrapper.find("img.bg-image");
    expect(medal.attributes("src")).toBe("/assets/images/medal/1.webp");
  });

  it("renders the correct avatar image", async () => {
    wrapper.unmount();
    wrapper = mountComponent();
    expect(getAvatarUrlByName).toBeCalledWith(props.winner.img_key);
    // await wrapper.vm.$nextTick();

    const avatarImg = wrapper.find("img.avatar-image");
    expect(avatarImg.attributes("src")).toContain(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Sophia"
    );
  });

  it("renders winner details correctly", () => {
    expect(wrapper.text()).toContain("JOHN");
    expect(wrapper.text()).toContain("john_doe");
    expect(wrapper.text()).toContain("8256");
  });
});

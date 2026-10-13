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
    global: {
      stubs: {
        NuxtImg: {
          template:
            '<img :src="$attrs.src" :alt="$attrs.alt" class="size-full object-cover" />',
        },
      },
    },
  });
};

let wrapper = mountComponent();

vi.mock("../../composables/avatar.js", () => {
  return {
    getAvatarUrlByName: vi
      .fn()
      .mockReturnValue("https://api.dicebear.com/10.x/clay/svg?seed=Sophia"),
  };
});

describe("WinnerCard test", () => {
  it("renders the correct rank label", () => {
    expect(wrapper.text()).toContain("1st Place");
  });

  it("renders the correct rank label for rank 2", async () => {
    props.winner.rank = 2;
    wrapper.unmount();
    wrapper = mountComponent();
    expect(wrapper.text()).toContain("2nd Place");
  });

  it("renders the correct rank label for rank 3", async () => {
    props.winner.rank = 3;
    wrapper.unmount();
    wrapper = mountComponent();
    expect(wrapper.text()).toContain("3rd Place");
  });

  it("renders the correct avatar image", async () => {
    wrapper.unmount();
    wrapper = mountComponent();
    expect(getAvatarUrlByName).toBeCalledWith(props.winner.img_key);

    const avatarImg = wrapper.find("img");
    expect(avatarImg.attributes("src")).toContain(
      "https://api.dicebear.com/10.x/clay/svg?seed=Sophia"
    );
  });

  it("renders winner details correctly", () => {
    expect(wrapper.text()).toContain("John");
    expect(wrapper.text()).toContain("john_doe");
    expect(wrapper.text()).toContain("8256");
  });
});

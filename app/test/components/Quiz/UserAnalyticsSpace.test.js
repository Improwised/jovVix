import { mount } from "@vue/test-utils";
import { describe, it, expect, vi } from "vitest";
import UserAnalyticsSpace from "~/components/Quiz/UserAnalyticsSpace.vue";

vi.mock("~~/composables/avatar", () => ({
  getAvatarUrlByName: vi.fn(
    () => "https://api.dicebear.com/10.x/clay/svg?seed=Eden"
  ),
}));

const props = {
  data: [
    {
      firstname: "John",
      username: "john_doe",
      avatar: "Eden",
      rank: 1,
    },
  ],
  userName: "john_doe",
};

const mountComponent = () =>
  mount(UserAnalyticsSpace, {
    props,
    global: {
      stubs: {
        NuxtImg: {
          template:
            '<img :src="$attrs.src" :alt="$attrs.alt" width="56" height="56" class="size-14 shrink-0 rounded-full border-[3px] border-jv-ink bg-jv-slate object-cover shadow-brutal-sm" />',
        },
      },
    },
  });

let wrapper = mountComponent();

describe("UserAnalyticsSpace test", () => {
  it("renders the component with provided props", () => {
    expect(wrapper.text()).toContain("John");
    expect(wrapper.text()).toContain("john_doe");
  });

  it("renders the avatar image", async () => {
    const w = mount(UserAnalyticsSpace, {
      props,
      global: {
        stubs: {
          NuxtImg: {
            template: '<img :src="$attrs.src" :alt="$attrs.alt" />',
          },
        },
      },
    });

    const avatarImg = w.find("img");
    expect(avatarImg.exists()).toBe(true);
  });
});

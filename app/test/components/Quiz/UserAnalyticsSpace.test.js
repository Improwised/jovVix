import { mount } from "@vue/test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import UserAnalyticsSpace from "~/components/Quiz/UserAnalyticsSpace.vue";

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
  });

let wrapper = mountComponent();

describe("UserAnalyticsSpace test", () => {
  it("renders the component with provided props", () => {
    expect(wrapper.find(".name").text()).toContain("John (john_doe)");
    expect(wrapper.find(".avatar").attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Eden"
    );
  });

  it("computes the correct avatar URL", async () => {
    const wrapper = mount(UserAnalyticsSpace, { props });

    const avatarImg = wrapper.find(".avatar");
    expect(avatarImg.attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Eden"
    );
  });
});

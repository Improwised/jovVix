import { describe, it, expect, vi } from "vitest";
import { mount } from "@vue/test-utils";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";

vi.mock("~~/composables/avatar", () => ({
  getAvatarUrlByName: vi.fn(
    (seed) => `https://api.dicebear.com/10.x/clay/svg?seed=${seed}`
  ),
}));

describe("ScoreBoardTable.vue", () => {
  const props = {
    scoreboardData: [
      {
        rank: 1,
        username: "john_doe",
        firstname: "John",
        score: 150,
        img_key: "Sophia",
      },
      {
        rank: 2,
        username: "jane_doe",
        firstname: "Jane",
        score: 140,
        img_key: "Jude",
      },
    ],
    isAdmin: true,
    userName: "",
  };

  const mountComponent = () => {
    return mount(ScoreBoardTable, {
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

  it("renders the table with the correct data", () => {
    const items = wrapper.findAll("li");
    expect(items.length).toBe(2);

    expect(wrapper.text()).toContain("John");
    expect(wrapper.text()).toContain("150");
  });

  it("does not highlight any row when the user is an admin", () => {
    const highlightedRow = wrapper.find('[class*="bg-jv-yellow-soft"]');
    expect(highlightedRow.exists()).toBe(false);
  });

  it("displays additional user information for admins", () => {
    expect(wrapper.text()).toContain("john_doe");
  });

  it("highlights the current user row when not an admin", async () => {
    const newProps = {
      scoreboardData: props.scoreboardData,
      isAdmin: false,
      userName: "jane_doe",
    };

    const w = mount(ScoreBoardTable, {
      props: newProps,
      global: {
        stubs: {
          NuxtImg: true,
        },
      },
    });

    expect(w.text()).toContain("Jane");
  });

  it("renders avatar URLs correctly", () => {
    const avatars = wrapper.findAll("li img");
    expect(avatars.length).toBe(2);
  });
});

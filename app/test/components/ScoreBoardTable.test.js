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
    expect(wrapper.text()).not.toContain("john_doe");
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
    expect(w.text()).not.toContain("jane_doe");
    expect(w.text()).not.toContain("john_doe");
    expect(w.text()).toContain("You");
  });

  it("renders avatar URLs correctly", () => {
    const avatars = wrapper.findAll("li img");
    expect(avatars.length).toBe(2);
  });

  it("renders an ellipsis row for a divider entry", () => {
    const w = mount(ScoreBoardTable, {
      props: {
        scoreboardData: [
          { rank: 1, username: "john_doe", firstname: "John", score: 150 },
          { rank: 2, username: "jane_doe", firstname: "Jane", score: 140 },
          { rank: 3, username: "carol", firstname: "Carol", score: 130 },
          { divider: true },
          { rank: 4, username: "you", firstname: "You", score: 120 },
        ],
        isAdmin: false,
        userName: "you",
      },
      global: {
        stubs: {
          NuxtImg: true,
        },
      },
    });

    const dots = w.find('[aria-hidden="true"] > .size-2.rounded-full');
    expect(dots.exists()).toBe(true);
    // One divider row plus four players.
    expect(w.findAll("li").length).toBe(5);
    expect(w.text()).toContain("4 players");
  });

  it("renders plain rank numbers when rankIcons is false", () => {
    const w = mount(ScoreBoardTable, {
      props: {
        scoreboardData: [
          { rank: 1, username: "a", firstname: "A", score: 100 },
          { rank: 2, username: "b", firstname: "B", score: 90 },
          { rank: 3, username: "c", firstname: "C", score: 80 },
        ],
        rankIcons: false,
        isAdmin: true,
        userName: "",
      },
      global: { stubs: { NuxtImg: true } },
    });

    expect(w.find('[aria-label="First place"]').exists()).toBe(false);
    expect(w.find('[aria-label="Rank 1"]').exists()).toBe(true);
    expect(w.find('[aria-label="Rank 2"]').exists()).toBe(true);
    expect(w.find('[aria-label="Rank 3"]').exists()).toBe(true);
    // Top 3 keep their colored boxes, only the icon is swapped for a number.
    expect(w.find('[aria-label="Rank 1"]').classes()).toContain("bg-jv-yellow");
    expect(w.find('[aria-label="Rank 2"]').classes()).toContain("bg-jv-salmon");
    expect(w.find('[aria-label="Rank 3"]').classes()).toContain("bg-jv-mint");
  });

  it("uses the totalPlayers prop for the count when provided", () => {
    const w = mount(ScoreBoardTable, {
      props: {
        scoreboardData: [
          { rank: 1, username: "a", firstname: "A", score: 100 },
          { rank: 6, username: "you", firstname: "You", score: 20 },
        ],
        totalPlayers: 35,
        isAdmin: false,
        userName: "you",
      },
      global: { stubs: { NuxtImg: true } },
    });

    expect(w.text()).toContain("35 players");
    expect(w.text()).not.toContain("2 players");
  });
});

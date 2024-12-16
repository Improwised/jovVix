import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import ScoreBoardTable from "~/components/ScoreBoardTable.vue";
import constants from "../constants";

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
          VCard: constants.slotTemplate,
          VCardItem: constants.slotTemplate,
        },
      },
    });
  };

  let wrapper = mountComponent();

  it("renders the table with the correct data", () => {
    const rows = wrapper.findAll("tbody tr");
    expect(rows.length).toBe(2);

    const firstRow = rows[0];
    expect(firstRow.html()).toContain("1");
    expect(firstRow.html()).toContain("John");
    expect(firstRow.html()).toContain("150");
  });

  it("does not highlight any row when the user is an admin", () => {
    const highlightedRow = wrapper.find("tr.table-primary");
    expect(highlightedRow.exists()).toBe(false);
  });

  it("displays additional user information for admins", () => {
    const adminRow = wrapper.find("tbody tr");
    expect(adminRow.html()).toContain("John (john_doe)");
  });

  it("highlights the current user row when not an admin", async () => {
    props.userName = "jane_doe";
    props.isAdmin = false;

    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    const highlightedRow = wrapper.find("tr.table-primary");
    expect(highlightedRow.exists()).toBe(true);
    expect(highlightedRow.html()).toContain("Jane");
  });

  it("renders avatar URLs correctly", () => {
    const avatars = wrapper.findAll("img");
    expect(avatars[0].attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Sophia&scale=75"
    );
    expect(avatars[1].attributes("src")).toBe(
      "https://api.dicebear.com/9.x/bottts/svg?seed=Jude&scale=75"
    );
  });
});

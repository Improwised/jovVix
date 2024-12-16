import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount } from "@vue/test-utils";
import ShareQuizAuthorizeUser from "~/components/Quiz/ShareQuizAuthorizeUser.vue";
import constants from "~/test/constants";

const props = {
  user: {
    id: 1,
    shared_to: "quiz123",
    permission: "Read",
    first_name: { Valid: true, String: "John" },
    last_name: { Valid: true, String: "Doe" },
  },
};

let wrapper = mount(ShareQuizAuthorizeUser, {
  props,
  global: {
    stubs: {
      FontAwesomeIcon: true,
      VListItemTitle: constants.slotTemplate,
      VBadge: constants.slotTemplate,
      VAvatar: constants.slotTemplate,
    },
  },
});

describe("ShareQuizAuthorizeUser test", () => {
  it("renders user details correctly", () => {
    expect(wrapper.text()).toContain("John Doe");
    expect(wrapper.text()).toContain("quiz123");
    expect(wrapper.text()).toContain("Read");
  });

  it("renders fallback when user name is invalid", async () => {
    await wrapper.setProps({
      user: {
        ...props.user,
        first_name: { Valid: false, String: "" },
        last_name: { Valid: false, String: "" },
      },
    });
    wrapper.vm.$nextTick();
    expect(wrapper.text()).toContain("Unknown");
  });

  it("renders avatar with correct image source", () => {
    const avatarImg = wrapper.find("img");
    expect(avatarImg.attributes("src")).toContain("https://api.dicebear.com");
    expect(avatarImg.attributes("alt")).toBe("props.user.title");
  });

  it("emits 'showEditForm' event with correct arguments when edit button is clicked", async () => {
    const editButton = wrapper.find("button[title='Edit Permission']");
    await editButton.trigger("click");

    expect(wrapper.emitted("showEditForm")).toBeTruthy();
    expect(wrapper.emitted("showEditForm")[0]).toEqual([1, "quiz123", "Read"]);
  });

  it("emits 'deleteUserPermission' event with user ID when delete button is clicked", async () => {
    const deleteButton = wrapper.find(
      "button[title='Edit Permission']:last-child"
    );
    await deleteButton.trigger("click");

    expect(wrapper.emitted("deleteUserPermission")).toBeTruthy();
    expect(wrapper.emitted("deleteUserPermission")[0]).toEqual([1]);
  });

  it("displays user permissions correctly", () => {
    const permissionText = wrapper.text();
    expect(permissionText).toContain("Read");
  });
});

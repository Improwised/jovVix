import { describe, it, expect, vi, afterEach } from "vitest";
import { mount } from "@vue/test-utils";
import ShareQuizModal from "~/components/Quiz/ShareQuizModal.vue";
import ShareQuizAuthorizeUser from "~/components/Quiz/ShareQuizAuthorizeUser.vue";
import ShareQuizForm from "~/components/Quiz/ShareQuizForm.vue";
import constants from "~/test/constants";

let wrapper = mount(ShareQuizModal, {
  props: {
    quizId: "1",
  },
  global: {
    stubs: {
      FontAwesomeIcon: true,
      VCard: constants.slotTemplate,
      VCardText: constants.slotTemplate,
      VList: constants.slotTemplate,
      VListItem: constants.slotTemplate,
      VListItemTitle: true,
    },
  },
});
const mockData = {
  data: [
    { id: "1", shared_to: "user1@example.com", permission: "read" },
    { id: "2", shared_to: "user2@example.com", permission: "write" },
  ],
};
describe("ShareQuizModal test", () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders the modal title", () => {
    expect(wrapper.text()).toContain("Share Quiz");
  });

  it("renders authorized users", async () => {
    wrapper.vm.quizAuthorizedUsersPending = false;
    wrapper.vm.quizAuthorizedUsersError = null;
    wrapper.vm.quizAuthorizedUsersData = mockData;
    await wrapper.vm.$nextTick();
    const authorizedUsers = wrapper.findAllComponents(ShareQuizAuthorizeUser);
    expect(authorizedUsers.length).toBe(2);
    expect(authorizedUsers[0].props("user")).toEqual(mockData.data[0]);
    expect(authorizedUsers[1].props("user")).toEqual(mockData.data[1]);
  });

  it("shows the 'Add People' form when button is clicked", async () => {
    const addButton = wrapper.find("button[title='Add People']");
    await addButton.trigger("click");

    expect(wrapper.findComponent(ShareQuizForm).exists()).toBe(true);
    expect(wrapper.findComponent(ShareQuizForm).props("formTitle")).toBe(
      "Add People"
    );
  });
});

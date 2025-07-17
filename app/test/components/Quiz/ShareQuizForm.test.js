import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import ShareQuizForm from "~/components/Quiz/ShareQuizForm.vue";

describe("ShareOrEditPermissionForm.vue", () => {
  const props = {
    formTitle: "Share Quiz",
    id: "",
    email: "",
    permission: "",
  };

  let wrapper = mount(ShareQuizForm, {
    props,
  });

  it("renders form title correctly", () => {
    expect(wrapper.text()).toContain("Share Quiz");
  });

  it("renders inputs and button for sharing a quiz", () => {
    expect(wrapper.find("input#email").exists()).toBe(true);
    expect(wrapper.find("select#permission").exists()).toBe(true);
    expect(wrapper.find("button").text()).toBe("Share Quiz");
  });

  it("renders inputs and button for updating permissions when id is provided", async () => {
    await wrapper.setProps({
      id: "123",
      email: "user@example.com",
      permission: "read",
    });
    await wrapper.vm.$nextTick();

    expect(wrapper.find("input#email").attributes("disabled")).toBe("");
    expect(wrapper.find("button").text()).toBe("Update Access");
  });

  it("updates local state when props change", async () => {
    await wrapper.setProps({ email: "new@example.com", permission: "write" });
    await wrapper.vm.$nextTick();

    expect(wrapper.find("input#email").element.value).toBe("new@example.com");
    expect(wrapper.find("select#permission").element.value).toBe("write");
  });

  it("emits 'shareQuiz' with correct payload when form is submitted without an id", async () => {
    await wrapper.setProps(props);
    await wrapper.vm.$nextTick();

    const emailInput = wrapper.find("input#email");
    const permissionSelect = wrapper.find("select#permission");

    await emailInput.setValue("test@example.com");
    await permissionSelect.setValue("write");
    await wrapper.find("form").trigger("submit.prevent");

    expect(wrapper.emitted("shareQuiz")).toBeTruthy();
    expect(wrapper.emitted("shareQuiz")[0]).toEqual([
      "test@example.com",
      "write",
    ]);
  });

  it("emits 'updateUserPermission' with correct payload when form is submitted with an id", async () => {
    await wrapper.setProps({ id: "123", email: "user@example.com" });
    const permissionSelect = wrapper.find("select#permission");

    await permissionSelect.setValue("share");
    await wrapper.find("form").trigger("submit.prevent");

    expect(wrapper.emitted("updateUserPermission")).toBeTruthy();
    expect(wrapper.emitted("updateUserPermission")[0]).toEqual([
      "123",
      "user@example.com",
      "share",
    ]);
  });

  it("clears input fields after form submission", async () => {
    const emailInput = wrapper.find("input#email");
    const permissionSelect = wrapper.find("select#permission");

    await emailInput.setValue("test@example.com");
    await permissionSelect.setValue("write");
    await wrapper.find("form").trigger("submit.prevent");

    expect(emailInput.element.value).toBe("");
    expect(permissionSelect.element.value).toBe("");
  });
});

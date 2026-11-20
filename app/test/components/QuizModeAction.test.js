import { it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import QuizModeAction from "~/components/QuizModeAction.vue";
it("shows a single action without a tooltip", async () => {
  const wrapper = mount(QuizModeAction, { props: { label: "Play solo" } });
  expect(wrapper.findAll("button")).toHaveLength(1);
  expect(wrapper.find('[role="tooltip"]').exists()).toBe(false);
  await wrapper.find("button").trigger("click");
  expect(wrapper.emitted("click")).toHaveLength(1);
});
it("disables the action while pending", async () => {
  const wrapper = mount(QuizModeAction, {
    props: { label: "Host Live", disabled: true },
  });
  await wrapper.find("button").trigger("click");
  expect(wrapper.emitted("click")).toBeUndefined();
});

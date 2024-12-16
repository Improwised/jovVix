import { mountSuspended } from "@nuxt/test-utils/runtime";
import { it, expect, describe } from "vitest";
import ActionCard from "~~/components/Quiz/ActionCard.vue";

const mountComponent = async (props) => {
  return await mountSuspended(ActionCard, {
    props,
  });
};

describe("ActionCard Test", () => {
  it("can mount component", async () => {
    const wrapper = await mountComponent({
      actionTitle: "Test-Action",
    });

    expect(wrapper.text()).toContain("Test-Action");
  });

  it("emmit the card click event", async () => {
    const wrapper = await mountComponent({
      actionTitle: "Test-Action",
    });
    await wrapper.find("div").trigger("click");

    expect(wrapper.emitted()).toHaveProperty("cardClick");
  });
});

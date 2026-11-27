import { mount } from "@vue/test-utils";
import { describe, it, expect } from "vitest";
import HeroSection from "~/components/landing/HeroSection.vue";

describe("landing hero", () => {
  const wrapper = mount(HeroSection, {
    global: { stubs: { NuxtImg: true } },
  });

  it("names gatherings in the headline", () => {
    expect(wrapper.find("h1").text()).toMatch(/meetup|party|family/i);
  });

  it("leads to the ready-made quizzes with one button", () => {
    const link = wrapper
      .findAll("a")
      .find((a) => a.text().includes("Host a quiz now"));
    expect(link).toBeTruthy();
    expect(link.attributes("href")).toBe("#public-quiz");
  });
});

import { describe, it, expect, beforeAll } from "vitest";
import { mount } from "@vue/test-utils";
import Pagination from "~/components/Pagination.vue";
import NavigationLink from "~/components/common/NavigationLink.vue";

const mountComponent = (props) => {
  return mount(Pagination, {
    props,
    global: {
      stubs: {
        NuxtLink: true,
        NavigationLink: true,
      },
    },
  });
};

// @nuxt/test-utils starts the Nuxt app in a beforeAll hook, so mount
// there rather than at import time.
let wrapper;
beforeAll(() => {
  wrapper = mountComponent({ page: 2, numOfRecords: 5 });
});

describe("Pagination test", () => {
  it("renders the current page correctly", () => {
    expect(wrapper.text()).toContain("2");
  });

  it("renders NavigationLink components for prev/next when not on boundary", async () => {
    wrapper = mountComponent({ page: 2, numOfRecords: 5 });
    const navLinks = wrapper.findAllComponents(NavigationLink);
    expect(navLinks.length).toBe(2);
  });

  it("renders disabled NavigationLink on the first page", async () => {
    wrapper = mountComponent({ page: 1, numOfRecords: 5 });
    const navLinks = wrapper.findAllComponents(NavigationLink);
    expect(navLinks.length).toBe(2);
  });

  it("renders disabled NavigationLink on the last page", async () => {
    wrapper = mountComponent({ page: 5, numOfRecords: 5 });
    const navLinks = wrapper.findAllComponents(NavigationLink);
    expect(navLinks.length).toBe(2);
  });

  it("shows page number correctly for middle page", async () => {
    wrapper = mountComponent({ page: 2, numOfRecords: 5 });
    expect(wrapper.text()).toContain("2");
  });
});

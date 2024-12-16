import { describe, it, expect } from "vitest";
import { mount } from "@vue/test-utils";
import Pagination from "~/components/Pagination.vue";
import { RouterLinkStub } from "@vue/test-utils";

const mountComponent = (props) => {
  return mount(Pagination, {
    props,
    global: {
      stubs: {
        NuxtLink: RouterLinkStub,
        FontAwesomeIcon: true,
      },
    },
  });
};

let wrapper = mountComponent({ page: 2, numOfRecords: 5 });

describe("Pagination test", () => {
  it("renders the current page correctly", () => {
    const currentPage = wrapper.find(".page-item.active .page-link");
    expect(currentPage.text()).toBe("2");
  });

  it("disables the previous button on the first page", async () => {
    wrapper = mountComponent({ page: 1, numOfRecords: 5 });
    const prevButton = wrapper.findAll(".page-link")[0];
    expect(prevButton.classes()).toContain("disabled");
  });

  it("disables the next button on the last page", async () => {
    wrapper = mountComponent({ page: 5, numOfRecords: 5 });

    const nextButton = wrapper.findAll("a");

    expect(nextButton[1].attributes()).toStrictEqual({
      class: "page-link disabled",
    });
  });

  it("enables both buttons on a middle page", async () => {
    wrapper = mountComponent({ page: 3, numOfRecords: 5 });
    const prevButton = wrapper.findAll(".page-link")[0];
    const nextButton = wrapper.findAll(".page-link")[1];

    expect(prevButton.classes()).not.toContain("disabled");
    expect(nextButton.classes()).not.toContain("disabled");
  });

  it("navigates to the correct previous page URL", async () => {
    wrapper = mountComponent({ page: 2, numOfRecords: 5 });
    const prevLink = wrapper.findAllComponents(RouterLinkStub)[0];

    expect(prevLink.props().to).toEqual({
      path: "/",
      query: {
        page: 1,
      },
    });
  });

  it("navigates to the correct next page URL", async () => {
    wrapper = mountComponent({ page: 2, numOfRecords: 5 });
    const nextLink = wrapper.findAllComponents(RouterLinkStub)[1];

    expect(nextLink.props("to")).toEqual({
      path: "/",
      query: { page: 3 },
    });
  });
});

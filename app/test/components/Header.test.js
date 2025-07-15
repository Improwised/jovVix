import { describe, it, expect, vi, beforeEach } from "vitest";
import { RouterLinkStub, mount } from "@vue/test-utils";
import Header from "~/components/Header.vue";
import { useRouter } from "vue-router";
import { useSessionStore } from "~~/store/session";
import { useUsersStore } from "~~/store/users";
import { createTestingPinia } from "@pinia/testing";

const pinia = createTestingPinia({
  createSpy: vi.fn,
});

const mockRouter = {
  push: vi.fn(),
};

const userData = useUsersStore(pinia);

const mountComponent = () => {
  return mount(Header, {
    global: {
      plugins: [pinia],
      stubs: {
        NuxtLink: RouterLinkStub,
        FontAwesomeIcon: true,
      },
      mocks: {
        router: mockRouter,
      },
    },
  });
};

let wrapper = mountComponent();

describe("Header test", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it("renders the logo and links to the home page", () => {
    const logoLink = wrapper.findComponent(RouterLinkStub);
    const logo = logoLink.find("img");
    expect(logo.attributes("src")).toBe("/jovvix-logo.png");
    expect(logoLink.props().to).toBe("/");
  });

  it("renders login and signup buttons for non-admin users", () => {
    const navs = wrapper.find("ul").findAll("li");
    expect(navs[1].text()).toContain("Log in");
    expect(navs[2].text()).toContain("Sign up");
  });

  it("conditionally renders admin options when the user is an admin", async () => {
    userData.getUserData = vi.fn().mockReturnValue({
      avatar: "Sophia",
      name: "John Doe",
      role: "admin-user",
    });
    wrapper = mountComponent();
    const navs = wrapper.find("ul").findAll("li");
    expect(navs.length).toBe(6);
    expect(navs[1].text()).toContain("Quizzes");
  });
});

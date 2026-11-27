import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import { takeNewSignup } from "~/composables/first_quiz_welcome";
import RegisterPage from "~/pages/account/register.vue";

vi.mock("notivue", () => ({
  usePush: () => ({ success: vi.fn(), error: vi.fn(), warning: vi.fn() }),
}));

const mountPage = async () => {
  // The page asks Kratos for a registration flow on load; answer with nothing.
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response("{}", { status: 503 }))
  );
  const Parent = defineComponent({
    render: () => h(Suspense, null, { default: () => h(RegisterPage) }),
  });
  const wrapper = mount(Parent, {
    global: { stubs: { NuxtLink: true, NuxtImg: true } },
  });
  await flushPromises();
  return wrapper.findComponent(RegisterPage);
};

describe("register page: new sign-up flag", () => {
  beforeEach(() => sessionStorage.clear());

  it("marks the tab when the form is submitted", async () => {
    const page = await mountPage();
    const submit = vi.fn();
    page.vm.onSubmit({ target: { submit } });
    expect(submit).toHaveBeenCalled();
    expect(takeNewSignup()).toBe(true);
  });

  it("does not mark the tab when the form has errors", async () => {
    const page = await mountPage();
    page.vm.password = "short";
    const submit = vi.fn();
    page.vm.onSubmit({ target: { submit } });
    expect(submit).not.toHaveBeenCalled();
    expect(takeNewSignup()).toBe(false);
  });
});

import { mockNuxtImport, registerEndpoint } from "@nuxt/test-utils/runtime";
import { flushPromises, mount } from "@vue/test-utils";
import { describe, it, expect, beforeAll, beforeEach, vi } from "vitest";
import { defineComponent, h, Suspense } from "vue";
import { useRuntimeConfig } from "#imports";
import AdminProfile from "~/pages/admin/index.vue";

// The page calls useFetch without await; hand it empty refs so it renders without network access.
const vue = await vi.hoisted(() => import("vue"));
mockNuxtImport("useFetch", () => () => ({
  data: vue.ref(null),
  error: vue.ref(null),
  status: vue.ref("success"),
  pending: vue.ref(false),
  refresh: () => {},
}));

const toast = vi.hoisted(() => ({
  success: vi.fn(),
  error: vi.fn(),
  warning: vi.fn(),
}));
vi.mock("notivue", () => ({ usePush: () => toast }));

// What DELETE /kratos/user answers next.
let deleteResponse;

let page;
const Parent = defineComponent({
  render: () => h(Suspense, null, { default: () => h(AdminProfile) }),
});

beforeAll(async () => {
  vi.stubGlobal("confirm", () => true);
  registerEndpoint(useRuntimeConfig().public.apiUrl + "/kratos/user", {
    method: "DELETE",
    handler: () => deleteResponse(),
  });
  const wrapper = mount(Parent, {
    global: {
      stubs: {
        ClientOnly: { template: "<div><slot /></div>" },
        NuxtImg: true,
        NuxtLink: true,
        QuizListCard: true,
        Pagination: true,
        UtilsQuizListWaiting: true,
      },
    },
  });
  await flushPromises();
  page = wrapper.findComponent(AdminProfile);
});

beforeEach(() => {
  toast.error.mockReset();
});

const jsonError = (status, body) =>
  new Response(JSON.stringify(body), {
    status,
    headers: { "content-type": "application/json" },
  });

describe("Admin profile: delete account", () => {
  it("shows the API's error message when deletion fails", async () => {
    // The API's JSONError body (jsend).
    deleteResponse = () =>
      jsonError(500, {
        status: "error",
        code: 500,
        message: "error while deleting user, please try after sometime",
      });
    await page.vm.deleteAccount();
    expect(toast.error).toHaveBeenCalledWith(
      "error while deleting user, please try after sometime"
    );
  });

  it("falls back to a generic message when the API sends none", async () => {
    deleteResponse = () => jsonError(502, {});
    await page.vm.deleteAccount();
    expect(toast.error).toHaveBeenCalledWith("Failed to delete the user.");
  });
});

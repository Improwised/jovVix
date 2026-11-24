import { mockNuxtImport } from "@nuxt/test-utils/runtime";
import { mount } from "@vue/test-utils";
import { describe, it, expect, vi, beforeEach } from "vitest";
import { defineComponent, nextTick, ref } from "vue";
import {
  SESSION_EXPIRED_MESSAGE,
  isSessionExpired,
  useSessionExpiredRedirect,
} from "~/composables/session_expired";

const nav = vi.hoisted(() => vi.fn());
mockNuxtImport("navigateTo", () => nav);
mockNuxtImport("useRoute", () => () => ({
  fullPath: "/admin/reports?page=2",
}));
const toast = vi.hoisted(() => ({ warning: vi.fn() }));
vi.mock("notivue", () => ({ usePush: () => toast }));

const mountWith = (error) =>
  mount(
    defineComponent({
      setup() {
        useSessionExpiredRedirect(error);
        return () => null;
      },
    })
  );

beforeEach(() => {
  nav.mockReset();
  toast.warning.mockReset();
});

// GitHub #77
describe("session expired", () => {
  it("recognises a 401 in the shapes useFetch and $fetch give", () => {
    expect(isSessionExpired({ statusCode: 401 })).toBe(true);
    expect(isSessionExpired({ data: { code: 401 } })).toBe(true);
    expect(isSessionExpired({ response: { status: 401 } })).toBe(true);
    expect(isSessionExpired({ statusCode: 403 })).toBe(false);
    expect(isSessionExpired(null)).toBe(false);
  });

  it("tells the person and sends them to sign in, then back", async () => {
    const error = ref(null);
    mountWith(error);
    expect(nav).not.toHaveBeenCalled();

    error.value = { statusCode: 401, data: { code: 401 } };
    await nextTick();

    expect(toast.warning).toHaveBeenCalledWith(SESSION_EXPIRED_MESSAGE);
    expect(nav).toHaveBeenCalledWith(
      "/account/login?returnTo=%2Fadmin%2Freports%3Fpage%3D2",
      { replace: true }
    );
  });

  it("redirects at once when the page loads with an expired session", () => {
    mountWith(ref({ statusCode: 401 }));
    expect(nav).toHaveBeenCalledTimes(1);
  });

  it("leaves other errors to the page", async () => {
    const error = ref(null);
    mountWith(error);
    error.value = { statusCode: 500 };
    await nextTick();
    expect(nav).not.toHaveBeenCalled();
    expect(toast.warning).not.toHaveBeenCalled();
  });
});

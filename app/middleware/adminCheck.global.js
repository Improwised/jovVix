import { callWithNuxt } from "nuxt/app";
import { constructPath } from "~/composables/url_operation";

export default defineNuxtRouteMiddleware(async (to) => {
  if (to.fullPath.startsWith("/admin")) {
    const app = useNuxtApp();
    const is_admin = await useIsAdmin();

    if (!is_admin.ok) {
      to.query["error"] = is_admin.err;
      delete to.query["url"];
      to.query["t"] = new Date().valueOf();
      const url = constructPath(to.path, to.params, to.query);

      const login_url = constructPath("/account/login", {}, url);
      return callWithNuxt(app, () => navigateTo(login_url));
    }
  }
});

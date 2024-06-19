import { callWithNuxt } from "nuxt/app";

export default defineNuxtRouteMiddleware(async (to) => {
  if (to.fullPath.startsWith("/admin")) {
    let is_admin;

    is_admin = await useIsAdmin();

    if (!is_admin.ok) {
      const app = useNuxtApp();
      return callWithNuxt(app, () =>
        navigateTo(
          "/account/login?t=" + new Date().valueOf() + "&url=" + to.fullPath
        )
      );
    }
  }
});

import { callWithNuxt } from "nuxt/app";

export default defineNuxtRouteMiddleware(async (to) => {
  if (to.fullPath.startsWith("/admin")) {
    let is_admin;
    let err;

    is_admin = await useIsAdmin();

    if (err != undefined || !is_admin.ok) {
      const app = useNuxtApp();
      return callWithNuxt(app, () =>
        navigateTo(
          "/account/login?error=" +
            (err || is_admin.err) +
            "&url=" +
            to.fullPath +
            "&t=" +
            new Date().valueOf()
        )
      );
    }
  }
});

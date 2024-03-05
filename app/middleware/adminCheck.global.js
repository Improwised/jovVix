import { callWithNuxt } from "nuxt/app";

export default defineNuxtRouteMiddleware(async (to) => {

  if (to.fullPath.startsWith("/admin")) {
    const is_admin = await useIsAdmin();

    if (!is_admin.ok) {
      const nuxtInstance = useNuxtApp();
      return callWithNuxt(nuxtInstance, () =>
        navigateTo(
          "/account/login?error=" + is_admin.err + "&url=" + to.fullPath
        )
      );
    }
  }
});

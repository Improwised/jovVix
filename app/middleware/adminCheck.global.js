import { callWithNuxt } from "nuxt/app";

export default defineNuxtRouteMiddleware(async (to) => {
  const { session, update, reset } = await useSession();
  const cookie = useCookie("user");

  if (!session.value.loginCheck || !cookie.value) {
    const user = await useGetUser();
    if (user.value.ok) {
      await update({ loginCheck: true, user: user.value.data });
    } else {
      await reset();
    }
  }

  if (to.fullPath.startsWith("/admin")) {
    const is_admin = await useIsAdmin();

    if (!is_admin.ok) {
      const nuxtInstance = useNuxtApp();
      console.log(is_admin);
      return callWithNuxt(nuxtInstance, () =>
        navigateTo(
          "/account/login?error=" + is_admin.err + "&url=" + to.fullPath
        )
      );
    }
  }
});

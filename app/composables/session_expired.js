import { watch } from "vue";
import { usePush } from "notivue";

// GitHub #77: when the login session expires on a registered-only page, say so in words
// and send the person to sign in, instead of printing a raw "401 Unauthorized" fetch error.
export const SESSION_EXPIRED_MESSAGE =
  "Your session has expired. Please sign in again.";

export const isSessionExpired = (error) =>
  !!error &&
  [
    error.statusCode,
    error.status,
    error.response?.status,
    error.data?.code,
  ].includes(401);

// Pass the error refs of the page's API calls. On the first 401 the page shows a toast and
// goes to the login page, which returns here after sign-in.
export const useSessionExpiredRedirect = (...errors) => {
  const route = useRoute();
  const toast = import.meta.client ? usePush() : null;
  watch(
    errors,
    (values) => {
      if (!values.some(isSessionExpired)) return;
      toast?.warning(SESSION_EXPIRED_MESSAGE);
      return navigateTo(
        `/account/login?returnTo=${encodeURIComponent(route.fullPath)}`,
        { replace: true }
      );
    },
    { immediate: true }
  );
};

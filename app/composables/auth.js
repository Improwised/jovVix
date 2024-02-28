import { useSystemEnv } from "./envs";

export async function useIsAdmin() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);

  const { error: err } = await useFetch(cfg.api_url + "/users/adminAccess", {
    method: "GET",
    credentials: "include",
    headers: headers,
    mode: "cors",
  });

  if (err?.value) {
    return { ok: false, err: err.value.data?.data || "unknown error" };
  }

  return { ok: true, err: null };
}

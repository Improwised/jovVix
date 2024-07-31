import { useSystemEnv } from "./envs";

export async function useIsAdmin() {
  const cfg = useRuntimeConfig().public;
  const headers = useRequestHeaders(["cookie"]);

  const { error: err, data: data } = await useFetch(
    cfg.api_url + "/user/is_admin",
    {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    }
  );

  if (err?.value) {
    return { ok: false, err: err.value.data?.data || "unknown error" };
  }

  return { ok: data?.value.data == true, err: null };
}

export async function useGetUser() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);
  const isLogin = useState("guestUser", () => {
    return { ok: false, data: "" };
  });

  const { error: err, data: data } = await useFetch(
    cfg.value.api_url + "/user/who",
    {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    }
  );

  if (err?.value) {
    isLogin.value.ok = false;
    return isLogin;
  }

  isLogin.value.ok = true;
  isLogin.value.data = data?.value.data;
  return isLogin;
}

export async function useGetKratosUser() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);
  const isKratosUser = useState("kratosUser", () => {
    return { ok: false, data: "" };
  });

  const { error: err, data: data } = await useFetch(
    cfg.value.kratos_url + "/sessions/whoami",
    {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    }
  );

  if (err?.value) {
    isKratosUser.value.ok = false;
    return isKratosUser;
  }

  isKratosUser.value.ok = true;
  isKratosUser.value.data = data?.value.data;
  return isKratosUser;
}

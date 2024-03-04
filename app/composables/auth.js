import { useSystemEnv } from "./envs";

export async function useIsAdmin() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);

  const { error: err, data: data } = await useFetch(
    cfg.value.api_url + "/user/is_admin",
    {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    }
  );

  if (err.value) {
    return { ok: false, err: err.value.data?.data };
  }

  return { ok: data?.value.data == true, err: null };
}

export async function useGetUser() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);
  const isLogin = useState("user", () => {
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

import { useSystemEnv } from "./envs";

export async function useIsAdmin() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);

  const { error: err, data: data } = await useFetch(
    cfg.value.api_url + "/users/adminAccess",
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

export async function updateSession(force = false) {
  const { session, update, reset } = await useSession();
  const cookie = useCookie("user");

  if (force || !session.value.user || !cookie.value) {
    const user = await useGetUser();
    if (user.value.ok) {
      await update({ user: user.value.data });
    } else {
      await reset();
    }
  }
}

export async function useGetUser() {
  const cfg = useSystemEnv();
  const headers = useRequestHeaders(["cookie"]);
  const isLogin = useState("user", () => {
    return { ok: false, data: "" };
  });

  const { error: err, data: data } = await useFetch(
    cfg.value.api_url + "/users/meta",
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

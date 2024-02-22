export async function useIsAdmin() {
  const cfg = useRuntimeConfig();
  const headers = useRequestHeaders(["cookie"]);
  const isLogin = useState("login", () => {
    return { ok: false, err: "" };
  });

  const { error: err, data: data } = await useFetch(
    cfg.public.api_url + "/users/adminAccess",
    {
      method: "GET",
      credentials: "include",
      headers: headers,
      mode: "cors",
    }
  );

  if (err?.value) {
    isLogin.value = false;
    return { ok: false, err: err.value.data.data };
  }

  isLogin.value = data?.value.data == true;
  return { ok: true, err: null };
}

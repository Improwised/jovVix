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
  const { kratos_url } = useRuntimeConfig().public;
  const isKratosUser = useState("kratosUser", () => {
    return { ok: false, data: "" };
  });
  try {
    await $fetch(kratos_url + "/sessions/whoami", {
      method: "GET",
      credentials: "include",
      headers: {
        Accept: "application/json",
      },
      onResponse({ response }) {
        if (response.status >= 200 && response.status < 300) {
          isKratosUser.value.ok = true;
          isKratosUser.value.data = response;
        }
      },
    });
  } catch (error) {}

  return isKratosUser;
}

export const handleLogout = async () => {
  const { kratos_url } = useRuntimeConfig().public;

  try {
    // Step 1: Fetch logout URL and token from the first API endpoint
    const response = await fetch(`${kratos_url}/self-service/logout/browser`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });

    if (!response.ok) {
      throw new Error("Failed to fetch logout URL and token");
    }

    const { logout_url } = await response.json();

    // Step 2: Use the fetched URL and token to make another API call
    const secondApiResponse = await fetch(`${logout_url}`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      credentials: "include",
    });

    if (!secondApiResponse.ok) {
      throw new Error(
        "Failed to perform logout with the provided URL and token"
      );
    }
    console.log("Logged out successfully");
    const user = useState("kratosUser");
    navigateTo("/");
    user.value.ok = false;
    user.value.data = null;
  } catch (error) {
    console.error("Error during logout:", error);
  }
};
